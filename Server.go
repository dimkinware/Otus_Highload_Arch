package main

import (
	"HighArch/api"
	"HighArch/service"
	"HighArch/storage"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type Server struct {
	userService        service.UserService
	registerService    service.RegisterService
	loginService       service.LoginService
	searchService      service.SearchService
	friendLinksService service.FriendLinksService
}

func NewServer(db *sqlx.DB) *Server {
	userStore := storage.NewDbUserStore(db)
	tokenStore := storage.NewDbTokenStore(db)
	friendLinksStore := storage.NewDbFriendLinksStore(db)
	return &Server{
		userService:        *service.NewUserService(userStore),
		registerService:    *service.NewRegisterService(userStore),
		loginService:       *service.NewLoginService(userStore, tokenStore),
		searchService:      *service.NewSearchService(userStore),
		friendLinksService: *service.NewFriendLinksService(friendLinksStore),
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
	_, err := getUserIdFromContext(req.Context())
	if err != nil {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	userId := mux.Vars(req)["id"]
	res, err := s.userService.GetUser(userId)
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
	//check auth
	_, err := getUserIdFromContext(req.Context())
	if err != nil {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

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

func (s *Server) GetFriendSetHandler(w http.ResponseWriter, req *http.Request) {
	//check auth
	currentUserId, err := getUserIdFromContext(req.Context())
	if err != nil {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}
	var newFriendUserId = mux.Vars(req)["id"]
	err = s.friendLinksService.SetFriendsLink(currentUserId, newFriendUserId)
	if err != nil {
		log.Println(err)
		if errors.Is(err, service.ErrorNotFound) {
			http.Error(w, "", http.StatusNotFound)
		} else if errors.Is(err, service.ErrorStoreError) {
			http.Error(w, "", http.StatusInternalServerError)
		} else {
			http.Error(w, "", http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (s *Server) GetFriendDeleteHandler(w http.ResponseWriter, req *http.Request) {
	//check auth
	currentUserId, err := getUserIdFromContext(req.Context())
	if err != nil {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}
	var friendUserId = mux.Vars(req)["id"]
	err = s.friendLinksService.DeleteFriendsLink(currentUserId, friendUserId)
	if err != nil {
		log.Println(err)
		if errors.Is(err, service.ErrorNotFound) {
			http.Error(w, "", http.StatusNotFound)
		} else if errors.Is(err, service.ErrorStoreError) {
			http.Error(w, "", http.StatusInternalServerError)
		} else {
			http.Error(w, "", http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

// Auth middleware methods

func (s *Server) GetAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		userId, err := s.loginService.Authenticate(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), userIdKey, *userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserIdFromContext(ctx context.Context) (string, error) {
	userId, ok := ctx.Value(userIdKey).(string)
	if !ok {
		return "", fmt.Errorf("user id not found in context")
	}
	return userId, nil
}

// Utils methods

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

const userIdKey string = "user_id"
