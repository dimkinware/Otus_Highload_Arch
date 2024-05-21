package main

import (
	"HighArch/api"
	"HighArch/service"
	"HighArch/storage"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js) // TODO: should handle error???
}

func parseJSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

// TODO move all Server's stuff to separate file

type Server struct {
	userService     service.UserService
	registerService service.RegisterService
	loginService    service.LoginService
}

func (s *Server) GetRegisterHandler(w http.ResponseWriter, req *http.Request) {
	var userDataModel api.RegisterApiModel
	var err = parseJSON(req, &userDataModel)
	if err != nil {
		// validation error
		println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
	} else {
		var res, err = s.registerService.Register(userDataModel)
		if err != nil {
			log.Println(err)
			if errors.Is(err, service.ErrorValidation) {
				w.WriteHeader(http.StatusBadRequest)
			} else if errors.Is(err, service.ErrorStoreError) {
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
		} else {
			w.WriteHeader(http.StatusOK)
			renderJSON(w, res)
		}
	}
}

func (s *Server) GetUserHandler(w http.ResponseWriter, req *http.Request) {
	var userId = mux.Vars(req)["id"]
	var res, err = s.userService.GetUser(userId)
	if err != nil {
		log.Println(err)
		if errors.Is(err, service.ErrorNotFound) {
			w.WriteHeader(http.StatusNotFound)
		} else if errors.Is(err, service.ErrorStoreError) {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(http.StatusOK)
		renderJSON(w, res)
	}
}

func (s *Server) GetLoginHandler(w http.ResponseWriter, req *http.Request) {
	var loginDataModel api.LoginApiModel
	err := parseJSON(req, &loginDataModel)
	if err != nil {
		// validation error
		println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
	} else {
		res, err := s.loginService.Login(loginDataModel)
		if err != nil {
			log.Println(err)
			if errors.Is(err, service.ErrorNotFound) {
				w.WriteHeader(http.StatusNotFound)
			} else if errors.Is(err, service.ErrorStoreError) {
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
		} else {
			w.WriteHeader(http.StatusOK)
			renderJSON(w, res)
		}
	}
}

func NewServer(db *sqlx.DB) *Server {
	userStore := storage.NewDbUserStore(db)
	tokenStore := storage.NewDbTokenStore(db)
	return &Server{
		userService:     *service.NewUserService(userStore),
		registerService: *service.NewRegisterService(userStore),
		loginService:    *service.NewLoginService(userStore, tokenStore),
	}
}

func main() {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
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
	router.HandleFunc("/user/get/{id}", server.GetUserHandler).Methods("GET")
	router.HandleFunc("/user/register", server.GetRegisterHandler).Methods("POST")
	router.HandleFunc("/login", server.GetLoginHandler).Methods("POST")

	appPort := os.Getenv("APP_PORT")
	println("Start listening server on port " + appPort)
	http.ListenAndServe("0.0.0.0:"+appPort, router)
}
