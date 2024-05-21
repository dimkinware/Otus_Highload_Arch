package api

type UserApiModel struct {
	UserId     string `json:"id"`
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	Birthdate  string `json:"birthdate"` // unix timestamp
	Gender     int    `json:"gender"`    // 0 - female, 1 - male TODO: make enum consts???
	Biography  string `json:"biography"`
	City       string `json:"city"`
}
