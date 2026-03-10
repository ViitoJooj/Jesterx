package domain

import "time"

type ReportStatus string
type ReportReason string

const (
	ReportStatusOpen       ReportStatus = "OPEN"
	ReportStatusInProgress ReportStatus = "IN_PROGRESS"
	ReportStatusResolved   ReportStatus = "RESOLVED"
	ReportStatusDismissed  ReportStatus = "DISMISSED"
)

const (
	ReportReasonSpam          ReportReason = "SPAM"
	ReportReasonFraud         ReportReason = "FRAUD"
	ReportReasonScam          ReportReason = "SCAM"
	ReportReasonInappropriate ReportReason = "INAPPROPRIATE"
	ReportReasonCounterfeit   ReportReason = "COUNTERFEIT"
	ReportReasonOther         ReportReason = "OTHER"
)

type Report struct {
	ID            string
	TicketNumber  int
	WebsiteID     string
	ReporterName  string
	ReporterEmail string
	Reason        ReportReason
	Description   string
	Status        ReportStatus
	AdminResponse *string
	ResolvedAt    *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
