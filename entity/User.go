package entity

type User struct {
	Id         string `db:"id"`
	FirstName  string `db:"first_name"`
	SecondName string `db:"second_name"`
	Birthdate  int64  `db:"birth_date"` // unix timestamp
	Gender     int    `db:"gender"`     // 0 - female, 1 - male
	Biography  string `db:"bio"`
	City       string `db:"city"`     // TODO: reference to City table?
	PwdHash    string `db:"pwd_hash"` // TODO: should be stored in dedicated table?
}
