package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ViitoJooj/Jesterx/internal/domain"
	middleware "github.com/ViitoJooj/Jesterx/internal/http/middlewares"
	"github.com/ViitoJooj/Jesterx/internal/service"
)

type ReportHandler struct {
	reportService *service.ReportService
	authService   *service.AuthService
}

func NewReportHandler(reportService *service.ReportService, authService *service.AuthService) *ReportHandler {
	return &ReportHandler{reportService: reportService, authService: authService}
}

type CreateReportRequest struct {
	WebsiteID     string   `json:"website_id"`
	ReporterName  string   `json:"reporter_name"`
	ReporterEmail string   `json:"reporter_email"`
	Reason        string   `json:"reason"`
	Description   string   `json:"description"`
	EvidenceURLs  []string `json:"evidence_urls"`
}

type ReportData struct {
	ID             string     `json:"id"`
	TicketNumber   int        `json:"ticket_number"`
	WebsiteID      string     `json:"website_id"`
	ReporterUserID *string    `json:"reporter_user_id,omitempty"`
	ReporterName   string     `json:"reporter_name"`
	ReporterEmail  string     `json:"reporter_email"`
	Reason         string     `json:"reason"`
	Description    string     `json:"description"`
	EvidenceURLs   []string   `json:"evidence_urls"`
	Status         string     `json:"status"`
	AdminResponse  *string    `json:"admin_response,omitempty"`
	ResolvedAt     *time.Time `json:"resolved_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type ReportResponse struct {
	Success bool       `json:"success"`
	Message string     `json:"message"`
	Data    ReportData `json:"data"`
}

type ReportsListResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Total   int          `json:"total"`
	Page    int          `json:"page"`
	PerPage int          `json:"per_page"`
	Data    []ReportData `json:"data"`
}

type UpdateReportRequest struct {
	Status        domain.ReportStatus `json:"status"`
	AdminResponse *string             `json:"admin_response,omitempty"`
}

func reportToData(r *domain.Report) ReportData {
	evidenceURLs := r.EvidenceURLs
	if evidenceURLs == nil {
		evidenceURLs = []string{}
	}
	return ReportData{
		ID:             r.ID,
		TicketNumber:   r.TicketNumber,
		WebsiteID:      r.WebsiteID,
		ReporterUserID: r.ReporterUserID,
		ReporterName:   r.ReporterName,
		ReporterEmail:  r.ReporterEmail,
		Reason:         string(r.Reason),
		Description:    r.Description,
		EvidenceURLs:   evidenceURLs,
		Status:         string(r.Status),
		AdminResponse:  r.AdminResponse,
		ResolvedAt:     r.ResolvedAt,
		CreatedAt:      r.CreatedAt,
		UpdatedAt:      r.UpdatedAt,
	}
}

// PublicCreateReport handles POST /api/v1/reports – anyone can submit a report.
// If the user is authenticated, reporter_name and reporter_email are auto-filled from their profile.
func (h *ReportHandler) PublicCreateReport(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req CreateReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Auto-fill reporter info from authenticated user if missing
	var reporterUserID *string
	if userID, ok := middleware.UserID(r.Context()); ok {
		user, err := h.authService.GetUserByID(userID)
		if err == nil {
			reporterUserID = &userID
			if strings.TrimSpace(req.ReporterName) == "" {
				req.ReporterName = strings.TrimSpace(user.First_name + " " + user.Last_name)
			}
			if strings.TrimSpace(req.ReporterEmail) == "" {
				req.ReporterEmail = user.Email
			}
		}
	}

	report, err := h.reportService.CreateReport(service.CreateReportInput{
		WebsiteID:      strings.TrimSpace(req.WebsiteID),
		ReporterUserID: reporterUserID,
		ReporterName:   strings.TrimSpace(req.ReporterName),
		ReporterEmail:  strings.TrimSpace(req.ReporterEmail),
		Reason:         strings.TrimSpace(req.Reason),
		Description:    strings.TrimSpace(req.Description),
		EvidenceURLs:   req.EvidenceURLs,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ReportResponse{
		Success: true,
		Message: "Denúncia enviada com sucesso",
		Data:    reportToData(report),
	})
}

// AdminListReports handles GET /api/v1/admin/reports – admin only.
func (h *ReportHandler) AdminListReports(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}

	reports, total, err := h.reportService.ListReports(status, page, perPage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := make([]ReportData, 0, len(reports))
	for i := range reports {
		data = append(data, reportToData(&reports[i]))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ReportsListResponse{
		Success: true,
		Message: "success",
		Total:   total,
		Page:    page,
		PerPage: perPage,
		Data:    data,
	})
}

// AdminGetReport handles GET /api/v1/admin/reports/{reportID} – admin only.
func (h *ReportHandler) AdminGetReport(w http.ResponseWriter, r *http.Request) {
	reportID := strings.TrimSpace(r.PathValue("reportID"))
	if reportID == "" {
		http.Error(w, "reportID required", http.StatusBadRequest)
		return
	}

	report, err := h.reportService.GetReport(reportID)
	if err != nil {
		http.Error(w, "denúncia não encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ReportResponse{
		Success: true,
		Message: "success",
		Data:    reportToData(report),
	})
}

// AdminUpdateReport handles PATCH /api/v1/admin/reports/{reportID} – admin only.
func (h *ReportHandler) AdminUpdateReport(w http.ResponseWriter, r *http.Request) {
	reportID := strings.TrimSpace(r.PathValue("reportID"))
	if reportID == "" {
		http.Error(w, "reportID required", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()
	var req UpdateReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	report, err := h.reportService.UpdateReport(reportID, service.UpdateReportInput{
		Status:        req.Status,
		AdminResponse: req.AdminResponse,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ReportResponse{
		Success: true,
		Message: "Denúncia atualizada",
		Data:    reportToData(report),
	})
}
