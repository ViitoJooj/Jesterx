package responses

type CreatePageDTO struct {
	Id        string `json:"id"`
	Page_id   string `json:"page_id"`
	Name      string `json:"name"`
	Tenant_id string `json:"tenant_id"`
}
