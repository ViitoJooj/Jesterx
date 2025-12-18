package responses

type MeResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Data    UserData `json:"data"`
}
