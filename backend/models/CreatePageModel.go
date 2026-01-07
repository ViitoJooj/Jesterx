package models

type CreatePageModels struct {
	Name       string   `json:"name" binding:"required"`
	PageID     string   `json:"page_id"`
	PageType   string   `json:"page_type" binding:"required"`
	Template   string   `json:"template"`
	Domain     string   `json:"domain"`
	Goal       string   `json:"goal"`
	Header     string   `json:"header"`
	Footer     string   `json:"footer"`
	ShowHeader *bool    `json:"show_header"`
	ShowFooter *bool    `json:"show_footer"`
	Components []string `json:"components"`
}
