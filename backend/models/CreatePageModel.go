package models

type CreatePageModels struct {
	Name   string `json:"name" binding:"required"`
	PageID string `json:"page_id"`
	Svelte string `json:"svelte" binding:"required"`
}
