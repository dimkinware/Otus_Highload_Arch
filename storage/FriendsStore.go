package storage

import (
	"HighArch/entity"
	"log"

	"github.com/jmoiron/sqlx"
)

type FriendLinksStore interface {
	SetFriends(link entity.FriendsLink) error
	DeleteFriends(link entity.FriendsLink) error
	GetFriendsIds(userId string) ([]string, error)
}

type dbFriendLinksStore struct {
	db *sqlx.DB
}

func (d dbFriendLinksStore) GetFriendsIds(userId string) ([]string, error) {
	query := "select user_id_f2 as friend_id from friends where user_id_f1=$1 union select user_id_f1 as friend_id from friends where user_id_f2=$1;"
	rows, err := d.db.Queryx(query, userId)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var firendsIds []string
	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		firendsIds = append(firendsIds, id)
	}
	return firendsIds, nil
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
