package service

import (
	"errors"

	"github.com/ViitoJooj/Jesterx/internal/domain"
	"github.com/ViitoJooj/Jesterx/internal/repository"
)

type StoreSocialService struct {
	repo        repository.StoreSocialRepository
	websiteRepo repository.WebsiteRepository
}

func NewStoreSocialService(repo repository.StoreSocialRepository, websiteRepo repository.WebsiteRepository) *StoreSocialService {
	return &StoreSocialService{repo: repo, websiteRepo: websiteRepo}
}

// ─── Store Info ──────────────────────────────────────────────────────────────

func (s *StoreSocialService) GetStoreFullInfo(websiteID string) (*domain.StoreFullInfo, error) {
	return s.repo.GetStoreFullInfo(websiteID)
}

// ─── Comments ────────────────────────────────────────────────────────────────

func (s *StoreSocialService) PostComment(websiteID, userID, content string, stars int) (*domain.StoreComment, error) {
	if len(content) < 3 || len(content) > 1000 {
		return nil, errors.New("comentário deve ter entre 3 e 1000 caracteres")
	}
	if stars < 1 || stars > 5 {
		return nil, errors.New("nota é obrigatória (1 a 5 estrelas)")
	}
	site, err := s.websiteRepo.FindWebSiteByID(websiteID)
	if err != nil || site == nil {
		return nil, errors.New("loja não encontrada")
	}
	// Auto-upsert the rating so the aggregate stays in sync
	_, _ = s.repo.UpsertRating(websiteID, userID, stars)
	_ = s.repo.RecalcRating(websiteID)

	return s.repo.SaveComment(domain.StoreComment{
		WebsiteID: websiteID,
		UserID:    userID,
		Content:   content,
		Stars:     &stars,
	})
}

func (s *StoreSocialService) ListComments(websiteID string) ([]domain.StoreComment, error) {
	return s.repo.ListComments(websiteID)
}

func (s *StoreSocialService) DeleteComment(commentID, requestingUserID, websiteID string) error {
	c, err := s.repo.FindCommentByID(commentID)
	if err != nil || c == nil {
		return errors.New("comentário não encontrado")
	}
	// Only the comment author OR a Jesterx platform admin/manager may delete.
	// Store owners and store team members have NO permission to delete user reviews.
	if c.UserID == requestingUserID {
		return s.repo.DeleteComment(commentID)
	}
	platformRole, err := s.repo.GetPlatformRole(requestingUserID)
	if err != nil {
		return errors.New("erro ao verificar permissão")
	}
	if platformRole != "admin" && platformRole != "manager" {
		return errors.New("sem permissão para excluir este comentário")
	}
	return s.repo.DeleteComment(commentID)
}

// ─── Ratings ─────────────────────────────────────────────────────────────────

func (s *StoreSocialService) RateStore(websiteID, userID string, stars int) (*domain.StoreRating, error) {
	if stars < 1 || stars > 5 {
		return nil, errors.New("estrelas devem ser entre 1 e 5")
	}
	site, err := s.websiteRepo.FindWebSiteByID(websiteID)
	if err != nil || site == nil {
		return nil, errors.New("loja não encontrada")
	}
	// Cannot rate own store
	if site.Creator_id == userID {
		return nil, errors.New("você não pode avaliar sua própria loja")
	}
	rating, err := s.repo.UpsertRating(websiteID, userID, stars)
	if err != nil {
		return nil, err
	}
	_ = s.repo.RecalcRating(websiteID)
	return rating, nil
}

func (s *StoreSocialService) GetMyRating(websiteID, userID string) (*domain.StoreRating, error) {
	return s.repo.GetUserRating(websiteID, userID)
}

// ─── Visits ──────────────────────────────────────────────────────────────────

func (s *StoreSocialService) RecordVisit(websiteID string) error {
	return s.repo.RecordVisit(websiteID)
}

func (s *StoreSocialService) GetVisitStats(websiteID string, days int) ([]domain.VisitDay, error) {
	return s.repo.GetVisitStats(websiteID, days)
}

// ─── Admin ───────────────────────────────────────────────────────────────────

func (s *StoreSocialService) SetMatureContent(websiteID string, mature bool) error {
	site, err := s.websiteRepo.FindWebSiteByID(websiteID)
	if err != nil || site == nil {
		return errors.New("loja não encontrada")
	}
	return s.repo.SetMatureContent(websiteID, mature)
}

// ─── Owner ───────────────────────────────────────────────────────────────────

func (s *StoreSocialService) UpdateStoreProfile(websiteID, ownerID, name, shortDesc, description string, image []byte) error {
	site, err := s.websiteRepo.FindWebSiteByID(websiteID)
	if err != nil || site == nil {
		return errors.New("loja não encontrada")
	}
	role, err := s.repo.GetUserRoleInStore(ownerID, websiteID)
	if err != nil {
		return err
	}
	if role != "owner" && role != "manager" && role != "admin" {
		return errors.New("sem permissão")
	}
	return s.repo.UpdateStoreProfile(websiteID, name, shortDesc, description, image)
}

// ─── Team Members ─────────────────────────────────────────────────────────────

// AddMember adds or updates a team member. Only owner, manager, or global admin can do this.
// Managers cannot assign the "manager" role — only the owner/admin can.
func (s *StoreSocialService) AddMember(websiteID, requesterID, targetUserID, role string) (*domain.StoreMember, error) {
	switch role {
	case domain.MemberRoleManager, domain.MemberRoleCatalogManager,
		domain.MemberRoleSupport, domain.MemberRoleLogistics:
	default:
		return nil, errors.New("role inválida")
	}
	requesterRole, err := s.repo.GetUserRoleInStore(requesterID, websiteID)
	if err != nil {
		return nil, err
	}
	if requesterRole != "owner" && requesterRole != "admin" {
		if requesterRole != "manager" {
			return nil, errors.New("sem permissão para gerenciar membros")
		}
		// managers cannot assign manager role
		if role == domain.MemberRoleManager {
			return nil, errors.New("apenas o dono pode definir gerentes")
		}
	}
	return s.repo.AddMember(domain.StoreMember{
		WebsiteID: websiteID,
		UserID:    targetUserID,
		Role:      role,
		InvitedBy: &requesterID,
	})
}

func (s *StoreSocialService) RemoveMember(websiteID, requesterID, targetUserID string) error {
	requesterRole, err := s.repo.GetUserRoleInStore(requesterID, websiteID)
	if err != nil {
		return err
	}
	if requesterRole != "owner" && requesterRole != "admin" && requesterRole != "manager" {
		return errors.New("sem permissão")
	}
	// Managers cannot remove other managers
	if requesterRole == "manager" {
		target, err := s.repo.FindMember(websiteID, targetUserID)
		if err != nil {
			return err
		}
		if target != nil && target.Role == domain.MemberRoleManager {
			return errors.New("gerentes não podem remover outros gerentes")
		}
	}
	return s.repo.RemoveMember(websiteID, targetUserID)
}

func (s *StoreSocialService) ListMembers(websiteID, requesterID string) ([]domain.StoreMember, error) {
	role, err := s.repo.GetUserRoleInStore(requesterID, websiteID)
	if err != nil {
		return nil, err
	}
	if role == "" {
		return nil, errors.New("sem permissão")
	}
	return s.repo.ListMembers(websiteID)
}

// UpdateMemberRole changes the role of an existing member.
// Permission rules mirror AddMember: managers cannot set/change to the "manager" role.
func (s *StoreSocialService) UpdateMemberRole(websiteID, requesterID, targetUserID, newRole string) (*domain.StoreMember, error) {
	switch newRole {
	case domain.MemberRoleManager, domain.MemberRoleCatalogManager,
		domain.MemberRoleSupport, domain.MemberRoleLogistics:
	default:
		return nil, errors.New("role inválida")
	}
	requesterRole, err := s.repo.GetUserRoleInStore(requesterID, websiteID)
	if err != nil {
		return nil, err
	}
	if requesterRole != "owner" && requesterRole != "admin" {
		if requesterRole != "manager" {
			return nil, errors.New("sem permissão para gerenciar membros")
		}
		if newRole == domain.MemberRoleManager {
			return nil, errors.New("apenas o dono pode definir gerentes")
		}
		// Managers cannot change other managers' roles
		target, err := s.repo.FindMember(websiteID, targetUserID)
		if err != nil {
			return nil, err
		}
		if target != nil && target.Role == domain.MemberRoleManager {
			return nil, errors.New("gerentes não podem alterar outros gerentes")
		}
	}
	return s.repo.UpdateMemberRole(websiteID, targetUserID, newRole)
}

// ─── Comment Replies ──────────────────────────────────────────────────────────

// ReplyComment allows store team (owner/manager/support/admin) to reply to a comment.
func (s *StoreSocialService) ReplyComment(websiteID, parentCommentID, userID, content string) (*domain.StoreComment, error) {
	if len(content) < 3 || len(content) > 1000 {
		return nil, errors.New("resposta deve ter entre 3 e 1000 caracteres")
	}
	role, err := s.repo.GetUserRoleInStore(userID, websiteID)
	if err != nil {
		return nil, err
	}
	switch role {
	case "owner", "admin", domain.MemberRoleManager, domain.MemberRoleSupport:
		// allowed
	default:
		return nil, errors.New("sem permissão para responder comentários")
	}
	parent, err := s.repo.FindCommentByID(parentCommentID)
	if err != nil || parent == nil {
		return nil, errors.New("comentário pai não encontrado")
	}
	if parent.ParentCommentID != nil {
		return nil, errors.New("não é possível responder a uma resposta")
	}
	return s.repo.SaveComment(domain.StoreComment{
		WebsiteID:       websiteID,
		UserID:          userID,
		Content:         content,
		ParentCommentID: &parentCommentID,
	})
}

// GetUserRoleInStore exposes the role resolution for handlers.
func (s *StoreSocialService) GetUserRoleInStore(userID, websiteID string) (string, error) {
	return s.repo.GetUserRoleInStore(userID, websiteID)
}
