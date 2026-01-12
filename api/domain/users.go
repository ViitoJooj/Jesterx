package domain

type User struct {
	UUID 	string `json:"uuid"`
	SiteID 	string `json:"site_id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Email 	string `json:"email"`
	Password string `json:"password"`
	Role 	string `json:"role"`
	UpdatedAt string `json:"updated_at"`
	CreatedAt string `json:"created_at"`
}