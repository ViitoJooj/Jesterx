package dtos

type CreateCupomRequest struct {
	TagUUID     string `json:"tag_uuid"`
	Label       string `json:"label"`
	Description string `json:"description"`
	Value       string `json:"value"`
	ValueType   string `json:"value_type"`
}

type CupomResponse struct {
	UUID        string `json:"uuid"`
	TagUUID     string `json:"tag_uuid"`
	Label       string `json:"label"`
	Description string `json:"description"`
	Value       string `json:"value"`
	ValueType   string `json:"value_type"`
}
