package api

type TokenGetResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

type TokenPostRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenPostResponse struct {
	Token string `json:"token"`
}
