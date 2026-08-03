package dtos

type CreateTermsAcceptedRequest struct {
	OwnerUUID string `json:"owner_uuid"`
	OwnerType string `json:"owner_type"`
}

type TermsAcceptedResponse struct {
	UUID         string `json:"uuid"`
	WebSiteUUID  string `json:"website_uuid"`
	OwnerUUID    string `json:"owner_uuid"`
	OwnerType    string `json:"owner_type"`
	AcceptedWhen string `json:"accepted_when"`
}
