package entity

type FriendsLink struct {
	Friend1UserId string `db:"user_id_f1"`
	Friend2UserId string `db:"user_id_f2"`
}
