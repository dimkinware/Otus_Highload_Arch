package service

import (
	"HighArch/api"
	"HighArch/entity"
	"HighArch/storage"
	"github.com/google/uuid"
	"time"
)

type LoginService struct {
	userStore  storage.UserStore
	tokenStore storage.TokenStore
}

func NewLoginService(userStore storage.UserStore, tokenStore storage.TokenStore) *LoginService {
	return &LoginService{
		userStore:  userStore,
		tokenStore: tokenStore,
	}
}

func (s *LoginService) Login(loginData api.LoginApiModel) (*api.LoginSuccessApiModel, error) {
	user, err := s.userStore.GetUser(loginData.UserId)
	if err != nil {
		return nil, ErrorStoreError
	}
	if user == nil {
		return nil, ErrorNotFound
	}

	if !comparePasswords(user.PwdHash, []byte(loginData.Password)) {
		return nil, ErrorValidation
	} else {
		// success
		newTokenInfo := entity.TokenInfo{
			UserId:     loginData.UserId,
			Token:      uuid.NewString(),
			ExpireTime: time.Now().UnixMilli() + ThirtyDaysMs,
		}
		err := s.tokenStore.CreateNewToken(newTokenInfo)
		if err != nil {
			return nil, ErrorStoreError
		}
		return &api.LoginSuccessApiModel{Token: newTokenInfo.Token}, nil
	}
}

const ThirtyDaysMs = 30*24*60*60 + 1000
