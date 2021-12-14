package auth

type Authentication struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type Token struct {
	Designation string `json:"designation"`
	Email string `json:"email"`
	TokenString string `json:"token"`
}
