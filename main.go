package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	appPort := os.Getenv("APP_PORT")

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

	server := NewServer(db)

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

	println("Start listening server on port " + appPort)
	http.ListenAndServe("0.0.0.0:"+appPort, router)
}
