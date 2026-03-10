package service

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/ViitoJooj/Jesterx/internal/domain"
	"github.com/ViitoJooj/Jesterx/internal/repository"
)

type RouteInput struct {
	Path         string
	Title        string
	RequiresAuth bool
}

type WebSiteService struct {
	webSiteRepo repository.WebsiteRepository
	userRepo    repository.UserRepository
	paymentRepo repository.PaymentRepository
}

func NewWebSiteService(webSiteRepo repository.WebsiteRepository, userRepo repository.UserRepository, paymentRepo repository.PaymentRepository) *WebSiteService {
	return &WebSiteService{
		webSiteRepo: webSiteRepo,
		userRepo:    userRepo,
		paymentRepo: paymentRepo,
	}
}

var acceptedTypes = [5]string{"ECOMMERCE", "LANDING_PAGE", "SOFTWARE_SELL", "COURSE", "VIDEO"}
var acceptedSourceTypes = map[string]bool{
	"JXML":           true,
	"REACT":          true,
	"SVELTE":         true,
	"ELEMENTOR_JSON": true,
}

func containsInvalidChars(name string) bool {
	for _, r := range name {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || r == ' ' {
			continue
		}
		return true
	}
	return false
}

func isValidType(rawType string) bool {
	normalizedType := strings.ToUpper(strings.TrimSpace(rawType))
	for _, allowedType := range acceptedTypes {
		if normalizedType == allowedType {
			return true
		}
	}
	return false
}

func normalizeRoutePath(path string) string {
	trimmed := strings.TrimSpace(path)
	if trimmed == "" {
		return ""
	}
	if strings.HasPrefix(trimmed, "/") {
		return trimmed
	}
	return "/" + trimmed
}

func routeMatches(pattern string, target string) bool {
	p := normalizeRoutePath(pattern)
	t := normalizeRoutePath(target)
	if p == t {
		return true
	}

	pSeg := strings.Split(strings.Trim(p, "/"), "/")
	tSeg := strings.Split(strings.Trim(t, "/"), "/")

	if len(pSeg) == 1 && pSeg[0] == "" && len(tSeg) == 1 && tSeg[0] == "" {
		return true
	}
	if len(pSeg) != len(tSeg) {
		return false
	}

	for i := range pSeg {
		if strings.HasPrefix(pSeg[i], ":") {
			continue
		}
		if pSeg[i] != tSeg[i] {
			return false
		}
	}
	return true
}

func (s *WebSiteService) getPlanLimits(planName string) (maxSites int, maxRoutes int) {
	if s.paymentRepo != nil {
		plan, err := s.paymentRepo.FindPlanByName(planName)
		if err == nil && plan != nil {
			return plan.MaxSites, plan.MaxRoutes
		}
	}
	// Fallback: infer from name
	normalized := strings.ToLower(strings.TrimSpace(planName))
	if strings.Contains(normalized, "enterprise") || strings.Contains(normalized, "ultra") {
		return 50, 100
	}
	if strings.Contains(normalized, "pro") || strings.Contains(normalized, "business") {
		return 10, 30
	}
	return 3, 8
}

func (s *WebSiteService) ensureActivePlan(userID string) (string, error) {
	user, err := s.userRepo.FindUserByID(userID)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("user not found")
	}
	if user.Plan == nil || strings.TrimSpace(*user.Plan) == "" {
		return "", errors.New("active plan required")
	}

	return strings.TrimSpace(*user.Plan), nil
}

func (s *WebSiteService) ensureOwnership(websiteID string, userID string) (*domain.WebSite, error) {
	website, err := s.webSiteRepo.FindWebSiteByID(websiteID)
	if err != nil {
		return nil, err
	}
	if website == nil {
		return nil, errors.New("website not found")
	}
	if website.Creator_id != userID {
		return nil, errors.New("forbidden")
	}
	return website, nil
}

func (s *WebSiteService) CreateWebSite(Type string, Image []byte, Name string, Short_description string, Description string, Creator_id string) (*domain.WebSite, error) {
	Type = strings.ToUpper(strings.TrimSpace(Type))
	Name = strings.TrimSpace(Name)

	plan, err := s.ensureActivePlan(Creator_id)
	if err != nil {
		return nil, err
	}

	maxSites, _ := s.getPlanLimits(plan)
	count, err := s.webSiteRepo.CountWebSitesByUserID(Creator_id)
	if err != nil {
		return nil, err
	}
	if count >= maxSites {
		return nil, fmt.Errorf("site limit reached for your plan (%d/%d)", count, maxSites)
	}

	if !isValidType(Type) {
		return nil, errors.New("invalid type")
	}

	if len(Name) < 3 || len(Name) > 50 || containsInvalidChars(Name) {
		return nil, errors.New("invalid name")
	}

	existing, err := s.webSiteRepo.FindWebSiteByName(Name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("this site already exists")
	}

	website := domain.NewWebSite(Type, Image, Name, strings.TrimSpace(Short_description), strings.TrimSpace(Description), Creator_id)
	if err := s.webSiteRepo.SaveWebSite(*website); err != nil {
		return nil, err
	}

	// Auto-register the creator as admin of the new store so they can log in immediately.
	if creator, err := s.userRepo.FindUserByID(Creator_id); err == nil && creator != nil {
		// Only create if no user with this email exists in this website yet
		if existing, _ := s.userRepo.FindUserByEmailAndWebsite(creator.Email, website.Id); existing == nil {
			adminUser := domain.NewUser(
				website.Id,
				creator.First_name, creator.Last_name,
				creator.Email, creator.Password,
				creator.AccountType,
			)
			adminUser.Role = "admin"
			adminUser.Verified_email = true
			adminUser.CpfCnpj = creator.CpfCnpj
			adminUser.AvatarUrl = creator.AvatarUrl
			adminUser.CompanyName = creator.CompanyName
			adminUser.TradeName = creator.TradeName
			adminUser.Phone = creator.Phone
			adminUser.ZipCode = creator.ZipCode
			adminUser.AddressStreet = creator.AddressStreet
			adminUser.AddressNumber = creator.AddressNumber
			adminUser.AddressComplement = creator.AddressComplement
			adminUser.AddressCity = creator.AddressCity
			adminUser.AddressState = creator.AddressState
			adminUser.AddressCountry = creator.AddressCountry
			_ = s.userRepo.UserRegister(*adminUser)
		}
	}

	return website, nil
}

func (s *WebSiteService) ReplaceRoutes(userID string, websiteID string, routes []RouteInput) ([]domain.WebSiteRoute, error) {
	plan, err := s.ensureActivePlan(userID)
	if err != nil {
		return nil, err
	}

	if _, err := s.ensureOwnership(websiteID, userID); err != nil {
		return nil, err
	}

	_, limit := s.getPlanLimits(plan)
	if len(routes) == 0 {
		return nil, errors.New("at least one route is required")
	}
	if len(routes) > limit {
		return nil, fmt.Errorf("route limit exceeded for plan: %d", limit)
	}

	normalized := make([]domain.WebSiteRoute, 0, len(routes))
	pathMap := make(map[string]bool, len(routes))
	for i, item := range routes {
		path := normalizeRoutePath(item.Path)
		title := strings.TrimSpace(item.Title)
		if len(path) < 1 || strings.Contains(path, " ") {
			return nil, errors.New("invalid route path")
		}
		if len(title) < 2 || len(title) > 100 {
			return nil, errors.New("invalid route title")
		}
		if pathMap[path] {
			return nil, errors.New("duplicate route path")
		}
		pathMap[path] = true

		route := domain.NewWebSiteRoute(websiteID, path, title, item.RequiresAuth, i)
		normalized = append(normalized, *route)
	}

	if err := s.webSiteRepo.ReplaceRoutesByWebsiteID(websiteID, normalized); err != nil {
		return nil, err
	}

	return s.webSiteRepo.ListRoutesByWebsiteID(websiteID)
}

func (s *WebSiteService) ListRoutes(userID string, websiteID string) ([]domain.WebSiteRoute, error) {
	if _, err := s.ensureOwnership(websiteID, userID); err != nil {
		return nil, err
	}
	return s.webSiteRepo.ListRoutesByWebsiteID(websiteID)
}

func (s *WebSiteService) CreateVersion(userID string, websiteID string, sourceType string, source string) (*domain.WebSiteVersion, *ScanReport, error) {
	if _, err := s.ensureActivePlan(userID); err != nil {
		return nil, nil, err
	}
	if _, err := s.ensureOwnership(websiteID, userID); err != nil {
		return nil, nil, err
	}

	normalizedSourceType := strings.ToUpper(strings.TrimSpace(sourceType))
	if !acceptedSourceTypes[normalizedSourceType] {
		return nil, nil, errors.New("invalid source_type")
	}

	if err := s.webSiteRepo.DeleteVersionsByWebsiteID(websiteID); err != nil {
		return nil, nil, err
	}

	scan := ScanWebsiteSource(normalizedSourceType, source)

	version := domain.NewWebSiteVersion(websiteID, 1, normalizedSourceType, source, userID)
	version.CompiledHTML = scan.CompiledHTML
	version.ScanStatus = scan.Report.Status
	version.ScanScore = scan.Report.Score
	version.ScanFindings = strings.Join(append(scan.Report.Findings, scan.Report.Errors...), " | ")
	version.Published = false

	if err := s.webSiteRepo.SaveVersion(*version); err != nil {
		return nil, nil, err
	}

	return version, &scan.Report, nil
}

func (s *WebSiteService) PublishVersion(userID string, websiteID string, versionNumber int) (*domain.WebSiteVersion, error) {
	if _, err := s.ensureOwnership(websiteID, userID); err != nil {
		return nil, err
	}

	version, err := s.webSiteRepo.FindVersionByWebsiteID(websiteID, versionNumber)
	if err != nil {
		return nil, err
	}
	if version == nil {
		return nil, errors.New("version not found")
	}
	if version.ScanStatus == "blocked" {
		return nil, errors.New("version blocked by security scan")
	}

	now := time.Now()
	if err := s.webSiteRepo.UpdateVersionPublishState(websiteID, versionNumber, true, &now); err != nil {
		return nil, err
	}

	return s.webSiteRepo.FindVersionByWebsiteID(websiteID, versionNumber)
}

func (s *WebSiteService) GetScanReport(userID string, websiteID string, versionNumber int) (*domain.WebSiteVersion, error) {
	if _, err := s.ensureOwnership(websiteID, userID); err != nil {
		return nil, err
	}

	version, err := s.webSiteRepo.FindVersionByWebsiteID(websiteID, versionNumber)
	if err != nil {
		return nil, err
	}
	if version == nil {
		return nil, errors.New("version not found")
	}

	return version, nil
}

func (s *WebSiteService) ListWebSites(userID string) ([]domain.WebSite, error) {
	return s.webSiteRepo.ListWebSitesByUserID(userID)
}

func (s *WebSiteService) ListVersions(userID string, websiteID string) ([]domain.WebSiteVersion, error) {
	if _, err := s.ensureOwnership(websiteID, userID); err != nil {
		return nil, err
	}
	return s.webSiteRepo.ListVersionsByWebsiteID(websiteID)
}

func (s *WebSiteService) GetPublicCompiledPage(websiteID string, path string) (string, error) {
	website, err := s.webSiteRepo.FindWebSiteByID(websiteID)
	if err != nil {
		return "", err
	}
	if website == nil {
		return "", errors.New("website not found")
	}
	if website.Banned {
		return "", errors.New("website is banned")
	}

	normalizedPath := normalizeRoutePath(path)
	if normalizedPath == "" {
		normalizedPath = "/"
	}
	route, err := s.webSiteRepo.FindRouteByWebsiteIDAndPath(websiteID, normalizedPath)
	if err != nil {
		return "", err
	}
	if route == nil {
		routes, err := s.webSiteRepo.ListRoutesByWebsiteID(websiteID)
		if err != nil {
			return "", err
		}
		matched := false
		for _, route := range routes {
			if routeMatches(route.Path, normalizedPath) {
				matched = true
				break
			}
		}
		if !matched {
			return "", errors.New("route not found")
		}
	}

	version, err := s.webSiteRepo.FindPublishedVersionByWebsiteID(websiteID)
	if err != nil {
		return "", err
	}
	if version == nil {
		return "", errors.New("published version not found")
	}

	return version.CompiledHTML, nil
}

func (s *WebSiteService) DeleteWebSite(userID string, websiteID string) error {
	if _, err := s.ensureOwnership(websiteID, userID); err != nil {
		return err
	}
	return s.webSiteRepo.DeleteWebSiteByID(websiteID)
}

func (s *WebSiteService) GetWebSiteByID(websiteID string) (*domain.WebSite, error) {
	return s.webSiteRepo.FindWebSiteByID(websiteID)
}
