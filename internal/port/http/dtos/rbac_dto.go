package dtos

type CreateRbacRequest struct {
	Label      string `json:"label"`
	CanRead    bool   `json:"can_read"`
	CanWrite   bool   `json:"can_write"`
	CanUpdate  bool   `json:"can_update"`
	CanUpgrade bool   `json:"can_upgrade"`
	CanDelete  bool   `json:"can_delete"`
}

type RbacResponse struct {
	UUID        string `json:"uuid"`
	WebSiteUUID string `json:"website_uuid"`
	Label       string `json:"label"`
	CanRead     bool   `json:"can_read"`
	CanWrite    bool   `json:"can_write"`
	CanUpdate   bool   `json:"can_update"`
	CanUpgrade  bool   `json:"can_upgrade"`
	CanDelete   bool   `json:"can_delete"`
	UpdatedAt   string `json:"updated_at"`
	CreatedAt   string `json:"created_at"`
}
