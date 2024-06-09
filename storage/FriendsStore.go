package storage

import (
	"HighArch/entity"
	"log"

	"github.com/jmoiron/sqlx"
)

type FriendLinksStore interface {
	SetFriends(link entity.FriendsLink) error
	DeleteFriends(link entity.FriendsLink) error
}

type dbFriendLinksStore struct {
	db *sqlx.DB
}

func NewDbFriendLinksStore(db *sqlx.DB) FriendLinksStore {
	return &dbFriendLinksStore{
		db: db,
	}
}

func (d dbFriendLinksStore) SetFriends(link entity.FriendsLink) error {
	query := "INSERT INTO friends(user_id_f1, user_id_f2) VALUES ($1, $2)"
	_, err := d.db.Exec(query, link.Friend1UserId, link.Friend2UserId)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (d dbFriendLinksStore) DeleteFriends(link entity.FriendsLink) error {
	query := "DELETE FROM friends where user_id_f1 = $1 AND user_id_f2 = $2"
	_, err := d.db.Exec(query, link.Friend1UserId, link.Friend2UserId)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
