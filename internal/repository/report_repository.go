package repository

import "github.com/ViitoJooj/Jesterx/internal/domain"

type ReportRepository interface {
	SaveReport(report domain.Report) (*domain.Report, error)
	FindReportByID(id string) (*domain.Report, error)
	ListReports(status string, page, perPage int) ([]domain.Report, int, error)
	ListReportsByWebsiteID(websiteID string) ([]domain.Report, error)
	UpdateReport(id string, status domain.ReportStatus, adminResponse *string) (*domain.Report, error)
}
