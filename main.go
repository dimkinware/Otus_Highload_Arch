package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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

	// connect to postgres
	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sqlx.Open("postgres", psqlconn)
	if err != nil {
		log.Fatal(err)
	}
	println("Ping PostgresSql")
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// connect to redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: redisPassword,
		DB:       0, // use default DB
	})
	println("Ping Redis")
	err = redisClient.Ping().Err()
	if err != nil {
		log.Fatal(err)
	}

	// starting Server
	server := NewServer(db, redisClient)
	router := mux.NewRouter()
	router.HandleFunc("/user/register", server.GetRegisterHandler).Methods("POST")
	router.HandleFunc("/login", server.GetLoginHandler).Methods("POST")

	//defining authenticated route
	privateRouter := router.PathPrefix("/").Subrouter()
	privateRouter.Use(server.GetAuthMiddleware)

	privateRouter.HandleFunc("/user/get/{id}", server.GetUserHandler).Methods("GET")
	privateRouter.HandleFunc("/user/search", server.GetSearchHandler).Methods("GET")
	privateRouter.HandleFunc("/friend/set/{id}", server.GetFriendSetHandler).Methods("PUT")
	privateRouter.HandleFunc("/friend/delete/{id}", server.GetFriendDeleteHandler).Methods("PUT")
	privateRouter.HandleFunc("/post/get/{id}", server.GetPostGetHandler).Methods("GET")
	privateRouter.HandleFunc("/post/create", server.GetPostCreateHandler).Methods("POST")
	privateRouter.HandleFunc("/post/feed", server.GetPostFeedHandler).Methods("GET")

	// start listening cache queue
	go server.FeedCacheController.ListenHandleFeedUpdate()

	println("Start listening server on port " + appPort)
	http.ListenAndServe("0.0.0.0:"+appPort, router)
}
