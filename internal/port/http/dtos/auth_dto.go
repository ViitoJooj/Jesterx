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

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token       string `json:"token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	UserUUID    string `json:"user_uuid"`
	WebsiteUUID string `json:"website_uuid"`
	Name        string `json:"name"`
	Email       string `json:"email"`
}
