package api

type RegisterApiModel struct {
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	Birthdate  string `json:"birthdate"`
	Gender     int    `json:"gender"` // 0 - female, 1 - male TODO: make enum consts???
	Biography  string `json:"biography"`
	City       string `json:"city"`
	Password   string `json:"password"`
}

type RegisterSuccessApiModel struct {
	UserId string `json:"user_id"`
}
