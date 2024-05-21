package service

import (
	"HighArch/api"
	"HighArch/entity"
	"HighArch/storage"
	"time"
)

type RegisterService struct {
	store storage.UserStore
}

func NewRegisterService(store storage.UserStore) *RegisterService {
	return &RegisterService{store: store}
}

func (s *RegisterService) Register(userDataModel api.RegisterApiModel) (*api.RegisterSuccessApiModel, error) {
	var err = validateRegisterModel(userDataModel)
	if err != nil {
		return nil, err
	}
	var pwdHash = hashAndSalt(userDataModel.Password)
	var newUser = entity.User{
		FirstName:  userDataModel.FirstName,
		SecondName: userDataModel.SecondName,
		Birthdate:  parseToUnixTimestamp(userDataModel.Birthdate, time.DateOnly),
		Gender:     userDataModel.Gender,
		Biography:  userDataModel.Biography,
		City:       userDataModel.City,
		PwdHash:    pwdHash,
	}
	id, err := s.store.CreateUser(newUser)
	if err != nil {
		return nil, ErrorStoreError
	}
	return &api.RegisterSuccessApiModel{UserId: *id}, nil
}

func validateRegisterModel(userDataModel api.RegisterApiModel) error {
	if len(userDataModel.FirstName) <= 0 {
		return ErrorValidation
	}
	if len(userDataModel.Password) <= 0 {
		return ErrorValidation
	}
	if validateTime(userDataModel.Birthdate, time.DateOnly) != nil {
		return ErrorValidation
	}
	// TODO validate length of all string fields
	return nil
}
