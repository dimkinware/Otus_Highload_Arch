package api

type LoginApiModel struct {
	UserId   string `json:"id"`
	Password string `json:"password"`
}

type LoginSuccessApiModel struct {
	Token string `json:"token"`
}
