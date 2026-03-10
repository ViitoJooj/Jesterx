package repository

import "github.com/ViitoJooj/Jesterx/internal/domain"

type StoreSocialRepository interface {
	GetStoreFullInfo(websiteID string) (*domain.StoreFullInfo, error)

	// Comments
	SaveComment(comment domain.StoreComment) (*domain.StoreComment, error)
	ListComments(websiteID string) ([]domain.StoreComment, error)
	FindCommentByID(id string) (*domain.StoreComment, error)
	DeleteComment(id string) error

	// Ratings
	UpsertRating(websiteID, userID string, stars int) (*domain.StoreRating, error)
	GetUserRating(websiteID, userID string) (*domain.StoreRating, error)
	RecalcRating(websiteID string) error

	// Visits
	RecordVisit(websiteID string) error
	GetVisitStats(websiteID string, days int) ([]domain.VisitDay, error)

	// Admin
	SetMatureContent(websiteID string, mature bool) error

	// Owner update
	UpdateStoreProfile(websiteID string, name, shortDesc, description string, image []byte) error

	// Team members
	AddMember(member domain.StoreMember) (*domain.StoreMember, error)
	RemoveMember(websiteID, userID string) error
	ListMembers(websiteID string) ([]domain.StoreMember, error)
	FindMember(websiteID, userID string) (*domain.StoreMember, error)
	// GetUserRoleInStore returns the user's role: "owner", a member role, "admin" (global), or "" if none.
	GetUserRoleInStore(userID, websiteID string) (string, error)
	// GetPlatformRole returns the user's global platform role ("admin", "manager", "user", etc.).
	GetPlatformRole(userID string) (string, error)
}
