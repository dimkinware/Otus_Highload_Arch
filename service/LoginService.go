package service

import (
	"HighArch/api"
	"HighArch/entity"
	"HighArch/storage"
	"github.com/google/uuid"
	"time"
)

type LoginService struct {
	userStore           storage.UserStore
	tokenStore          storage.TokenStore
	feedCacheController FeedCacheController
}

func NewLoginService(userStore storage.UserStore, tokenStore storage.TokenStore, feedCacheController FeedCacheController) *LoginService {
	return &LoginService{
		userStore:           userStore,
		tokenStore:          tokenStore,
		feedCacheController: feedCacheController,
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

		s.feedCacheController.InvalidateFeedCacheForUser(loginData.UserId)

		return &api.LoginSuccessApiModel{Token: newTokenInfo.Token}, nil
	}
}

func (s *LoginService) Authenticate(token string) (userId *string, err error) {
	tokenInfo, err := s.tokenStore.FindToken(token)
	if err != nil {
		return nil, err
	}
	if tokenInfo == nil {
		return nil, ErrorNotFound
	}
	if !checkTokenInfoIsValid(*tokenInfo) {
		return nil, ErrorTokenExpired
	}
	return &tokenInfo.UserId, nil
}

func checkTokenInfoIsValid(tokenInfo entity.TokenInfo) bool {
	if tokenInfo.ExpireTime <= time.Now().UnixMilli() || tokenInfo.UserId == "" {
		return false
	}
	return true
}

const ThirtyDaysMs = 30 * 24 * 60 * 60 * 1000
