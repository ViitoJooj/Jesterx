package service

import (
	"errors"

	"github.com/ViitoJooj/Jesterx/internal/domain"
	"github.com/ViitoJooj/Jesterx/internal/repository"
	"github.com/ViitoJooj/Jesterx/internal/security"
)

type ReportService struct {
	reportRepo  repository.ReportRepository
	websiteRepo repository.WebsiteRepository
}

func NewReportService(reportRepo repository.ReportRepository, websiteRepo repository.WebsiteRepository) *ReportService {
	return &ReportService{reportRepo: reportRepo, websiteRepo: websiteRepo}
}

type CreateReportInput struct {
	WebsiteID      string
	ReporterUserID *string
	ReporterName   string
	ReporterEmail  string
	Reason         string
	Description    string
	EvidenceURLs   []string
}

type UpdateReportInput struct {
	Status        domain.ReportStatus
	AdminResponse *string
}

func (s *ReportService) CreateReport(input CreateReportInput) (*domain.Report, error) {
	if input.WebsiteID == "" || input.ReporterName == "" || input.ReporterEmail == "" ||
		input.Reason == "" || input.Description == "" {
		return nil, errors.New("todos os campos são obrigatórios")
	}

	validReasons := map[string]bool{
		string(domain.ReportReasonSpam):          true,
		string(domain.ReportReasonFraud):         true,
		string(domain.ReportReasonScam):          true,
		string(domain.ReportReasonInappropriate): true,
		string(domain.ReportReasonCounterfeit):   true,
		string(domain.ReportReasonOther):         true,
	}
	if !validReasons[input.Reason] {
		return nil, errors.New("motivo inválido")
	}

	// Validate evidence (max 5, each base64 item max ~1.3MB)
	if len(input.EvidenceURLs) > 5 {
		return nil, errors.New("máximo de 5 imagens de evidência permitidas")
	}
	for _, ev := range input.EvidenceURLs {
		if len(ev) > 1_400_000 {
			return nil, errors.New("cada imagem deve ter no máximo 1MB")
		}
	}

	_, err := s.websiteRepo.FindWebSiteByID(input.WebsiteID)
	if err != nil {
		return nil, errors.New("loja não encontrada")
	}

	evidenceURLs := input.EvidenceURLs
	if evidenceURLs == nil {
		evidenceURLs = []string{}
	}

	report := domain.Report{
		WebsiteID:      input.WebsiteID,
		ReporterUserID: input.ReporterUserID,
		ReporterName:   input.ReporterName,
		ReporterEmail:  input.ReporterEmail,
		Reason:         domain.ReportReason(input.Reason),
		Description:    input.Description,
		EvidenceURLs:   evidenceURLs,
	}

	return s.reportRepo.SaveReport(report)
}

func (s *ReportService) ListReports(status string, page, perPage int) ([]domain.Report, int, error) {
	return s.reportRepo.ListReports(status, page, perPage)
}

func (s *ReportService) GetReport(id string) (*domain.Report, error) {
	return s.reportRepo.FindReportByID(id)
}

func (s *ReportService) UpdateReport(id string, input UpdateReportInput) (*domain.Report, error) {
	if id == "" {
		return nil, errors.New("id inválido")
	}

	validStatuses := map[domain.ReportStatus]bool{
		domain.ReportStatusOpen:       true,
		domain.ReportStatusInProgress: true,
		domain.ReportStatusResolved:   true,
		domain.ReportStatusDismissed:  true,
	}
	if !validStatuses[input.Status] {
		return nil, errors.New("status inválido")
	}

	existing, err := s.reportRepo.FindReportByID(id)
	if err != nil {
		return nil, errors.New("denúncia não encontrada")
	}

	updated, err := s.reportRepo.UpdateReport(id, input.Status, input.AdminResponse)
	if err != nil {
		return nil, err
	}

	// Send email notification if there is an admin response and reporter email exists
	if input.AdminResponse != nil && *input.AdminResponse != "" && existing.ReporterEmail != "" {
		_ = security.SendTicketResponseEmail(existing.ReporterEmail, existing.ReporterName, updated.TicketNumber, *input.AdminResponse, string(updated.Status))
	}

	return updated, nil
}
