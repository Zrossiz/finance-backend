package handler

type createUserDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginUserDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
