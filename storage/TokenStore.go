package storage

import (
	"HighArch/entity"
	"github.com/jmoiron/sqlx"
	"log"
)

type TokenStore interface {
	CreateNewToken(tokenInfo entity.TokenInfo) error
	FindToken(string) (*entity.TokenInfo, error)
}

type dbTokenStore struct {
	db *sqlx.DB
}

func NewDbTokenStore(db *sqlx.DB) TokenStore {
	return &dbTokenStore{
		db: db,
	}
}

func (d dbTokenStore) CreateNewToken(tokenInfo entity.TokenInfo) error {
	query := "INSERT INTO tokens(user_id, token, expired_time) VALUES ($1, $2, $3)"
	_, err := d.db.Exec(query, tokenInfo.UserId, tokenInfo.Token, tokenInfo.ExpireTime)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (d dbTokenStore) FindToken(token string) (*entity.TokenInfo, error) {
	panic("implement me")
}
