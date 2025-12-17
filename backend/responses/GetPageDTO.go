package responses

type GetPageResponse struct {
	Success   bool   `json:"success"`
	Id        string `json:"id"`
	Page_id   string `json:"page_id"`
	Tenant_id string `json:"tenant_id"`
	Svelte    string `json:"svelte"`
}
