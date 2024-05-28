package main

import (
	"HighArch/api"
	"HighArch/service"
	"HighArch/storage"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
)

type Server struct {
	userService     service.UserService
	registerService service.RegisterService
	loginService    service.LoginService
	searchService   service.SearchService
}

func NewServer(db *sqlx.DB) *Server {
	userStore := storage.NewDbUserStore(db)
	tokenStore := storage.NewDbTokenStore(db)
	return &Server{
		userService:     *service.NewUserService(userStore),
		registerService: *service.NewRegisterService(userStore),
		loginService:    *service.NewLoginService(userStore, tokenStore),
		searchService:   *service.NewSearchService(userStore),
	}
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

func (s *Server) GetSearchHandler(w http.ResponseWriter, req *http.Request) {
	firstNameStr := req.URL.Query().Get("first_name")
	lastNameStr := req.URL.Query().Get("last_name")
	if firstNameStr == "" || lastNameStr == "" {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		res, err := s.searchService.SearchByName(firstNameStr, lastNameStr)
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
			renderJSON(w, res)
		}
	}
}

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
