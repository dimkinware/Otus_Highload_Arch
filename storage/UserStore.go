package storage

import (
	"HighArch/entity"
	"errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"log"
)

type UserStore interface {
	GetUser(id string) (*entity.User, error)
	CreateUser(user entity.User) (*string, error)
}

type dbUserStore struct {
	db *sqlx.DB
}

func NewDbUserStore(db *sqlx.DB) UserStore {
	return &dbUserStore{
		db: db,
	}
}

func (p dbUserStore) GetUser(id string) (*entity.User, error) {
	rows, err := p.db.Queryx("SELECT * FROM users WHERE id = $1 limit 1", id)
	defer rows.Close()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if rows.Next() {
		var user entity.User
		err = rows.StructScan(&user)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return &user, nil
	}

	return nil, nil
}

func (p dbUserStore) CreateUser(user entity.User) (*string, error) {
	var userId = user.Id
	if len(userId) <= 0 {
		userId = uuid.NewString()
	}
	query := "INSERT INTO users(id, first_name, second_name, birth_date, gender, bio, city, pwd_hash) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
	_, err := p.db.Exec(query, userId, user.FirstName, user.SecondName, user.Birthdate, user.Gender, user.Biography, user.City, user.PwdHash)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &userId, nil
}

// region MockUserStore

type mockUserStore struct{}

func NewMockUserStore() UserStore {
	return &mockUserStore{}
}

func (m mockUserStore) GetUser(id string) (*entity.User, error) {
	// simulate success
	if id == "42" {
		var mockUser = &entity.User{
			Id:         "42",
			FirstName:  "42",
			SecondName: "42 42",
			Birthdate:  1485907200000,
			Gender:     1,
			Biography:  "42 42 42 42",
			City:       "424242",
		}
		return mockUser, nil
	}
	// simulate not found
	if id == "41" {
		return nil, nil
	}
	// simulate unknown error
	return nil, errors.New("internal store error")
}

func (m mockUserStore) CreateUser(user entity.User) (*string, error) {
	var userId = "100500"
	return &userId, nil
}

// endregion mockUserStore
