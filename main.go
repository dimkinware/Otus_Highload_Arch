package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
	"os"
)

func main() {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	appPort := os.Getenv("APP_PORT")
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	rabbitUser := os.Getenv("RABBITMQ_USER")
	rabbitPassword := os.Getenv("RABBITMQ_PASS")
	rabbitHost := os.Getenv("RABBITMQ_HOST")
	rabbitPort := os.Getenv("RABBITMQ_PORT")

	// connect to Postgres
	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	log.Println(psqlconn)
	db, err := sqlx.Open("postgres", psqlconn)
	failOnError(err, "Error connecting to PostgresSql database")
	log.Println("Ping PostgresSql")
	err = db.Ping()
	failOnError(err, "Error pinging PostgresSql")

	// connect to Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: redisPassword,
		DB:       0, // use default DB
	})
	log.Println("Ping Redis")
	err = redisClient.Ping().Err()
	failOnError(err, "Error pinging Redis")

	// connect to RabbitMQ
	rabbitConn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbitUser, rabbitPassword, rabbitHost, rabbitPort))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer rabbitConn.Close()
	log.Println("Create RabbitMQ channel")
	rabbitChannel, err := rabbitConn.Channel()
	failOnError(err, "Failed to open a channel to RabbitMQ")
	defer rabbitChannel.Close()

	// create Server
	server := NewServer(db, redisClient, rabbitChannel)
	router := mux.NewRouter()
	router.HandleFunc("/user/register", server.GetRegisterHandler).Methods("POST")
	router.HandleFunc("/login", server.GetLoginHandler).Methods("POST")

	// define authenticated route
	privateRouter := router.PathPrefix("/").Subrouter()
	privateRouter.Use(server.GetAuthMiddleware)
	privateRouter.HandleFunc("/user/get/{id}", server.GetUserHandler).Methods("GET")
	privateRouter.HandleFunc("/user/search", server.GetSearchHandler).Methods("GET")
	privateRouter.HandleFunc("/friend/set/{id}", server.GetFriendSetHandler).Methods("PUT")
	privateRouter.HandleFunc("/friend/delete/{id}", server.GetFriendDeleteHandler).Methods("PUT")
	privateRouter.HandleFunc("/post/get/{id}", server.GetPostGetHandler).Methods("GET")
	privateRouter.HandleFunc("/post/create", server.GetPostCreateHandler).Methods("POST")
	privateRouter.HandleFunc("/post/feed", server.GetPostFeedHandler).Methods("GET")
	privateRouter.HandleFunc("/post/feed/posted", server.GetPostFeedWsHandler)

	// start listening cache queue
	go server.FeedCacheController.ListenHandleFeedUpdate()

	// start server
	log.Println("Start listening server on port " + appPort)
	http.ListenAndServe("0.0.0.0:"+appPort, router)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
