package responses

type UserData struct {
	Id         string `json:"id"`
	ProfileImg string `json:"profile_img"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	Plan       string `json:"plan"`
}
