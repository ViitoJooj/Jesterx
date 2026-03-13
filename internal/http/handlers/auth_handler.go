package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ViitoJooj/Jesterx/internal/config"
	"github.com/ViitoJooj/Jesterx/internal/domain"
	middleware "github.com/ViitoJooj/Jesterx/internal/http/middlewares"
	"github.com/ViitoJooj/Jesterx/internal/security"
	"github.com/ViitoJooj/Jesterx/internal/service"
	"github.com/ViitoJooj/Jesterx/pkg/validate"
)

type RegisterRequest struct {
	First_name        string  `json:"first_name"`
	Last_name         string  `json:"last_name"`
	Email             string  `json:"email"`
	Password          string  `json:"password"`
	AccountType       string  `json:"account_type"`
	CompanyName       *string `json:"company_name"`
	TradeName         *string `json:"trade_name"`
	CpfCnpj           *string `json:"cpf_cnpj"`
	Phone             *string `json:"phone"`
	ZipCode           *string `json:"zip_code"`
	AddressStreet     *string `json:"address_street"`
	AddressNumber     *string `json:"address_number"`
	AddressComplement *string `json:"address_complement"`
	AddressCity       *string `json:"address_city"`
	AddressState      *string `json:"address_state"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserData struct {
	Id         string    `json:"id"`
	WebsiteId  string    `json:"website_id"`
	Email      string    `json:"email"`
	Plan       string    `json:"user_plan"`
	Updated_at time.Time `json:"updated_at"`
	Created_at time.Time `json:"created_at"`
}

type AuthResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Data    UserData `json:"data"`
}

type ResponseRefreshToken struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type UserMeResponse struct {
	ID                string `json:"id"`
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	Email             string `json:"email"`
	Role              string `json:"role"`
	Plan              string `json:"user_plan"`
	CpfCnpj           string `json:"cpf_cnpj"`
	AvatarUrl         string `json:"avatar_url"`
	PlanMaxSites      int    `json:"plan_max_sites"`
	PlanMaxRoutes     int    `json:"plan_max_routes"`
	AccountType       string `json:"account_type"`
	CompanyName       string `json:"company_name"`
	TradeName         string `json:"trade_name"`
	DisplayName       string `json:"display_name"`
	BirthDate         string `json:"birth_date"`
	Gender            string `json:"gender"`
	Bio               string `json:"bio"`
	Instagram         string `json:"instagram"`
	WebsiteUrl        string `json:"website_url"`
	Whatsapp          string `json:"whatsapp"`
	Phone             string `json:"phone"`
	ZipCode           string `json:"zip_code"`
	AddressStreet     string `json:"address_street"`
	AddressNumber     string `json:"address_number"`
	AddressComplement string `json:"address_complement"`
	AddressDistrict   string `json:"address_district"`
	AddressCity       string `json:"address_city"`
	AddressState      string `json:"address_state"`
	AddressCountry    string `json:"address_country"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
}

type UpdateProfileRequest struct {
	FirstName         string  `json:"first_name"`
	LastName          string  `json:"last_name"`
	CpfCnpj           *string `json:"cpf_cnpj"`
	AvatarUrl         *string `json:"avatar_url"`
	CompanyName       *string `json:"company_name"`
	TradeName         *string `json:"trade_name"`
	DisplayName       *string `json:"display_name"`
	BirthDate         *string `json:"birth_date"`
	Gender            *string `json:"gender"`
	Bio               *string `json:"bio"`
	Instagram         *string `json:"instagram"`
	WebsiteUrl        *string `json:"website_url"`
	Whatsapp          *string `json:"whatsapp"`
	Phone             *string `json:"phone"`
	ZipCode           *string `json:"zip_code"`
	AddressStreet     *string `json:"address_street"`
	AddressNumber     *string `json:"address_number"`
	AddressComplement *string `json:"address_complement"`
	AddressDistrict   *string `json:"address_district"`
	AddressCity       *string `json:"address_city"`
	AddressState      *string `json:"address_state"`
	AddressCountry    *string `json:"address_country"`
}

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println(err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := validate.New().
		Required("email", req.Email).
		Email("email", req.Email).
		Required("password", req.Password).
		MinLen("password", req.Password, 8).
		Required("first_name", req.First_name).
		MaxLen("first_name", req.First_name, 50).
		Required("last_name", req.Last_name).
		MaxLen("last_name", req.Last_name, 50).
		Err(); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	websiteId := r.Header.Get("X-Website-Id")
	if websiteId == "" {
		log.Println("websiteId is null")
		http.Error(w, "invalid website id", http.StatusBadRequest)
		return
	}

	user, err := h.authService.Register(service.RegisterInput{
		WebsiteId:         websiteId,
		FirstName:         req.First_name,
		LastName:          req.Last_name,
		Email:             req.Email,
		Password:          req.Password,
		AccountType:       req.AccountType,
		CompanyName:       req.CompanyName,
		TradeName:         req.TradeName,
		CpfCnpj:           req.CpfCnpj,
		Phone:             req.Phone,
		ZipCode:           req.ZipCode,
		AddressStreet:     req.AddressStreet,
		AddressNumber:     req.AddressNumber,
		AddressComplement: req.AddressComplement,
		AddressCity:       req.AddressCity,
		AddressState:      req.AddressState,
	})
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = security.SendVerifyEmail(user.Email, user.Id, user.WebsiteId)
	if err != nil {
		log.Println("Error on sending email")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := AuthResponse{
		Success: true,
		Message: "registered, please verify you email.",
		Data: UserData{
			Id:         user.Id,
			WebsiteId:  user.WebsiteId,
			Email:      user.Email,
			Plan:       derefString(user.Plan),
			Created_at: user.Created_at,
		},
	}

	log.Printf("User: %s; is registred", user.First_name+" "+user.Last_name)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	path := r.URL.Path
	prefix := "/api/v1/auth/verify/"

	if !strings.HasPrefix(path, prefix) {
		http.NotFound(w, r)
		return
	}

	id := strings.TrimPrefix(path, prefix)
	if id == "" {
		log.Println("ID is requered")
		http.Error(w, "token error", http.StatusBadRequest)
		return
	}

	websiteID := r.URL.Query().Get("website_id")
	user, err := h.authService.VerifyEmail(id, websiteID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	refreshClaims := security.RefreshTokenClaims{
		Iss:       "https://jesterx.com.br",
		Sub:       user.Id,
		WebsiteId: user.WebsiteId,
		Exp:       time.Now().Add(30 * 24 * time.Hour).Unix(),
	}
	refreshToken, err := security.RefreshToken(refreshClaims)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	accessClaims := security.AccessTokenClaims{
		Iss:       "https://jesterx.com.br",
		Sub:       user.Id,
		Aud:       "https://api.jesterx.com.br",
		WebsiteId: user.WebsiteId,
		Role:      user.Role,
		Exp:       time.Now().Add(15 * time.Minute).Unix(),
	}
	accessToken, err := security.AccessToken(accessClaims)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	secure := !config.IsDev
	http.SetCookie(w, &http.Cookie{
		Name:     security.RefreshCookieName(user.WebsiteId),
		Value:    refreshToken,
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 30,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     security.AccessCookieName(user.WebsiteId),
		Value:    accessToken,
		Path:     "/",
		MaxAge:   60 * 15,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, strings.TrimRight(config.FrontendURL, "/"), http.StatusFound)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	websiteId := r.Header.Get("X-Website-Id")
	if websiteId == "" {
		http.Error(w, "invalid website id", http.StatusBadRequest)
		return
	}

	user, err := h.authService.Login(
		websiteId,
		req.Email,
		req.Password,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	refreshClaims := security.RefreshTokenClaims{
		Iss:       "https://jesterx.com.br",
		Sub:       user.Id,
		WebsiteId: user.WebsiteId,
		Exp:       time.Now().Add(30 * 24 * time.Hour).Unix(),
	}

	refreshToken, err := security.RefreshToken(refreshClaims)
	if err != nil {
		http.Error(w, "Internal error", http.StatusBadGateway)
		return
	}

	accessClaims := security.AccessTokenClaims{
		Iss:       "https://jesterx.com.br",
		Sub:       user.Id,
		Aud:       "https://api.jesterx.com.br",
		WebsiteId: user.WebsiteId,
		Role:      user.Role,
		Exp:       time.Now().Add(15 * time.Minute).Unix(),
	}

	accessToken, err := security.AccessToken(accessClaims)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	secure := !config.IsDev

	var plan string
	if user.Plan != nil {
		plan = *user.Plan
	}

	resp := AuthResponse{
		Success: true,
		Message: "logged in.",
		Data: UserData{
			Id:         user.Id,
			WebsiteId:  user.WebsiteId,
			Email:      user.Email,
			Plan:       plan,
			Updated_at: user.Updated_at,
			Created_at: user.Created_at,
		},
	}

	http.SetCookie(w, &http.Cookie{
		Name:     security.RefreshCookieName(user.WebsiteId),
		Value:    refreshToken,
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 30,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     security.AccessCookieName(user.WebsiteId),
		Value:    accessToken,
		Path:     "/",
		MaxAge:   60 * 15,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	websiteId := r.Header.Get("X-Website-Id")
	if websiteId == "" {
		http.Error(w, "invalid website id", http.StatusBadRequest)
		return
	}

	refreshCookie, err := r.Cookie(security.RefreshCookieName(websiteId))
	if err != nil {
		http.Error(w, "refresh token missing", http.StatusUnauthorized)
		return
	}

	accessToken, err := h.authService.Refresh(refreshCookie.Value)
	if err != nil {
		log.Println(err)
		http.Error(w, "not allowed", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     security.AccessCookieName(websiteId),
		Value:    accessToken,
		Path:     "/",
		MaxAge:   60 * 15,
		HttpOnly: true,
		Secure:   !config.IsDev,
		SameSite: http.SameSiteLaxMode,
	})

	resp := ResponseRefreshToken{
		Success: true,
		Message: "success",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user, planLimits, err := h.authService.MeWithPlan(userID)
	if err != nil {
		http.Error(w, "not allowed", http.StatusUnauthorized)
		return
	}

	resp := UserMeResponse{
		ID:                user.Id,
		FirstName:         user.First_name,
		LastName:          user.Last_name,
		Email:             user.Email,
		Role:              user.Role,
		Plan:              derefString(user.Plan),
		CpfCnpj:           derefString(user.CpfCnpj),
		AvatarUrl:         derefString(user.AvatarUrl),
		PlanMaxSites:      planLimits[0],
		PlanMaxRoutes:     planLimits[1],
		AccountType:       user.AccountType,
		CompanyName:       derefString(user.CompanyName),
		TradeName:         derefString(user.TradeName),
		DisplayName:       derefString(user.DisplayName),
		BirthDate:         formatDate(user.BirthDate),
		Gender:            derefString(user.Gender),
		Bio:               derefString(user.Bio),
		Instagram:         derefString(user.Instagram),
		WebsiteUrl:        derefString(user.WebsiteUrl),
		Whatsapp:          derefString(user.Whatsapp),
		Phone:             derefString(user.Phone),
		ZipCode:           derefString(user.ZipCode),
		AddressStreet:     derefString(user.AddressStreet),
		AddressNumber:     derefString(user.AddressNumber),
		AddressComplement: derefString(user.AddressComplement),
		AddressDistrict:   derefString(user.AddressDistrict),
		AddressCity:       derefString(user.AddressCity),
		AddressState:      derefString(user.AddressState),
		AddressCountry:    derefString(user.AddressCountry),
		CreatedAt:         user.Created_at.Format(time.RFC3339),
		UpdatedAt:         user.Updated_at.Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		http.Error(w, `{"success":false,"message":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	websiteId := r.Header.Get("X-Website-Id")

	if err := h.authService.DeleteAccount(userID); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{"success": false, "message": err.Error()})
		return
	}

	if websiteId != "" {
		http.SetCookie(w, &http.Cookie{Name: security.RefreshCookieName(websiteId), Value: "", Path: "/", MaxAge: -1})
		http.SetCookie(w, &http.Cookie{Name: security.AccessCookieName(websiteId), Value: "", Path: "/", MaxAge: -1})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{"success": true, "message": "conta excluída"})
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func formatDate(d *time.Time) string {
	if d == nil {
		return ""
	}
	return d.Format("2006-01-02")
}

func parseBirthDate(raw *string) (*time.Time, error) {
	if raw == nil {
		return nil, nil
	}
	trimmed := strings.TrimSpace(*raw)
	if trimmed == "" {
		return nil, nil
	}
	parsed, err := time.Parse("2006-01-02", trimmed)
	if err != nil {
		return nil, errors.New("data de nascimento inválida")
	}
	return &parsed, nil
}

func (h *AuthHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	defer r.Body.Close()
	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	birthDate, err := parseBirthDate(req.BirthDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.authService.UpdateProfile(userID, domain.UpdateProfileData{
		FirstName:         req.FirstName,
		LastName:          req.LastName,
		CpfCnpj:           req.CpfCnpj,
		AvatarUrl:         req.AvatarUrl,
		CompanyName:       req.CompanyName,
		TradeName:         req.TradeName,
		DisplayName:       req.DisplayName,
		BirthDate:         birthDate,
		Gender:            req.Gender,
		Bio:               req.Bio,
		Instagram:         req.Instagram,
		WebsiteUrl:        req.WebsiteUrl,
		Whatsapp:          req.Whatsapp,
		Phone:             req.Phone,
		ZipCode:           req.ZipCode,
		AddressStreet:     req.AddressStreet,
		AddressNumber:     req.AddressNumber,
		AddressComplement: req.AddressComplement,
		AddressDistrict:   req.AddressDistrict,
		AddressCity:       req.AddressCity,
		AddressState:      req.AddressState,
		AddressCountry:    req.AddressCountry,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{"success": true, "message": "profile updated"})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	websiteId := r.Header.Get("X-Website-Id")
	if websiteId == "" {
		http.Error(w, "invalid website id", http.StatusBadRequest)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   security.RefreshCookieName(websiteId),
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	http.SetCookie(w, &http.Cookie{
		Name:   security.AccessCookieName(websiteId),
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	w.WriteHeader(http.StatusNoContent)
}
