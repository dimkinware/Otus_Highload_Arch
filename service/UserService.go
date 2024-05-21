package service

import (
	"HighArch/api"
	"HighArch/storage"
	"time"
)

type UserService struct {
	userStore storage.UserStore
}

func NewUserService(userStore storage.UserStore) *UserService {
	return &UserService{userStore: userStore}
}

func (s *UserService) GetUser(id string) (*api.UserApiModel, error) {
	var user, err = s.userStore.GetUser(id)
	if err != nil {
		return nil, ErrorStoreError
	}
	if user == nil {
		return nil, ErrorNotFound
	}

	var result = api.UserApiModel{
		UserId:     user.Id,
		FirstName:  user.FirstName,
		SecondName: user.SecondName,
		Birthdate:  formatUnixTimestampToString(user.Birthdate, time.DateOnly),
		Gender:     user.Gender,
		Biography:  user.Biography,
		City:       user.City,
	}

	return &result, nil
}
