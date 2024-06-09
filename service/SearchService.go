package service

import (
	"HighArch/api"
	"HighArch/storage"
	"time"
)

type SearchService struct {
	userStore storage.UserStore
}

func NewSearchService(userStore storage.UserStore) *SearchService {
	return &SearchService{userStore: userStore}
}

func (s *SearchService) SearchByName(firstNameStr, secondNameStr string) ([]api.UserApiModel, error) {
	users, err := s.userStore.SearchByName(firstNameStr, secondNameStr)
	if err != nil {
		return nil, ErrorStoreError
	}
	if len(users) == 0 {
		return nil, ErrorNotFound
	}

	usersCount := len(users)
	result := make([]api.UserApiModel, usersCount)
	for i := 0; i < usersCount; i++ {
		user := users[i]
		// TODO move to mapper
		result[i] = api.UserApiModel{
			UserId:     user.Id,
			FirstName:  user.FirstName,
			SecondName: user.SecondName,
			Birthdate:  formatUnixTimestampToString(user.Birthdate, time.DateOnly),
			Gender:     user.Gender,
			Biography:  user.Biography,
			City:       user.City,
		}
	}
	return result, nil
}
