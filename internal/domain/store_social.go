package domain

import "time"

// Store team member roles
const (
	MemberRoleManager        = "manager"
	MemberRoleCatalogManager = "catalog_manager"
	MemberRoleSupport        = "support"
	MemberRoleLogistics      = "logistics"
)

// StoreMember represents a team member of a store.
type StoreMember struct {
	ID        string    `json:"id"`
	WebsiteID string    `json:"website_id"`
	UserID    string    `json:"user_id"`
	UserName  string    `json:"user_name"`
	AvatarURL *string   `json:"avatar_url,omitempty"`
	Role      string    `json:"role"`
	InvitedBy *string   `json:"invited_by,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type StoreComment struct {
	ID              string         `json:"id"`
	WebsiteID       string         `json:"website_id"`
	UserID          string         `json:"user_id"`
	UserName        string         `json:"user_name"`
	AvatarURL       *string        `json:"avatar_url,omitempty"`
	Content         string         `json:"content"`
	Stars           *int           `json:"stars,omitempty"`
	ParentCommentID *string        `json:"parent_comment_id,omitempty"`
	Replies         []StoreComment `json:"replies,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

type StoreRating struct {
	ID        string
	WebsiteID string
	UserID    string
	Stars     int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type VisitDay struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type StoreCreator struct {
	ID          string  `json:"id"`
	FullName    string  `json:"full_name"`
	CompanyName *string `json:"company_name,omitempty"`
	TradeName   *string `json:"trade_name,omitempty"`
	AvatarURL   *string `json:"avatar_url,omitempty"`
	AccountType string  `json:"account_type"`
}

type StoreFullInfo struct {
	ID               string        `json:"id"`
	Name             string        `json:"name"`
	ShortDescription string        `json:"short_description"`
	Description      string        `json:"description"`
	Image            []byte        `json:"image,omitempty"`
	Type             string        `json:"type"`
	MatureContent    bool          `json:"mature_content"`
	RatingAvg        float64       `json:"rating_avg"`
	RatingCount      int           `json:"rating_count"`
	EditorType       string        `json:"editor_type"`
	Creator          StoreCreator  `json:"creator"`
	Managers         []StoreMember `json:"managers,omitempty"`
}
