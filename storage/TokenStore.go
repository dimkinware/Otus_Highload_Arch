package storage

import (
	"HighArch/entity"
	"log"

	"github.com/jmoiron/sqlx"
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
	rows, err := d.db.Queryx("SELECT * FROM tokens WHERE token = $1 ORDER BY expired_time DESC LIMIT 1", token)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		var result entity.TokenInfo
		err := rows.StructScan(&result)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return &result, nil
	}

	return nil, nil
}
