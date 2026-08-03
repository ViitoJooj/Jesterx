package dtos

type RegisterRequest struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	SaveLogin bool   `json:"save_login"`
}

type RegisterResponse struct {
	UUID      string `json:"uuid"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}
