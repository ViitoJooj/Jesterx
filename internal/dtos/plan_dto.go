package dtos

type CreatePlanRequest struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	MaxWebsites     int    `json:"max_websites"`
	MaxRouters      int    `json:"max_routers"`
	MaxProducts     int    `json:"max_products"`
	CostPerSaleRate int    `json:"cost_per_sale_rate"`
	Coin            string `json:"coin"`
	Price           int    `json:"price"`
}

type PlanResponse struct {
	UUID            string `json:"uuid"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	MaxWebsites     int    `json:"max_websites"`
	MaxRouters      int    `json:"max_routers"`
	MaxProducts     int    `json:"max_products"`
	CostPerSaleRate int    `json:"cost_per_sale_rate"`
	Coin            string `json:"coin"`
	Price           int    `json:"price"`
	CreatedAt       string `json:"created_at"`
}
