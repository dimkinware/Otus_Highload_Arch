package entity

type TokenInfo struct {
	UserId     string `db:"user_id"`
	Token      string `db:"token"`
	ExpireTime int64  `db:"expired_time"`
}
