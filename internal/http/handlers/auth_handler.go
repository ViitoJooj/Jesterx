package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/ViitoJooj/Jesterx/internal/config"
	"github.com/ViitoJooj/Jesterx/internal/domain"
	middleware "github.com/ViitoJooj/Jesterx/internal/http/middlewares"
	"github.com/ViitoJooj/Jesterx/internal/security"
	"github.com/ViitoJooj/Jesterx/internal/service"
	"github.com/ViitoJooj/Jesterx/pkg/validate"
)

type CompanyRegisterRequest struct {
	CompanyName       string  `json:"company_name"`
	TradeName         *string `json:"trade_name"`
	Cnpj              *string `json:"cnpj"`
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

type RegisterRequest struct {
	First_name string                  `json:"first_name"`
	Last_name  string                  `json:"last_name"`
	Email      string                  `json:"email"`
	Password   string                  `json:"password"`
	Cpf        *string                 `json:"cpf"`
	Phone      *string                 `json:"phone"`
	Company    *CompanyRegisterRequest `json:"company"`
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

type CompanyResponse struct {
	ID                string `json:"id"`
	CompanyName       string `json:"company_name"`
	TradeName         string `json:"trade_name"`
	Cnpj              string `json:"cnpj"`
	Phone             string `json:"phone"`
	ZipCode           string `json:"zip_code"`
	AddressStreet     string `json:"address_street"`
	AddressNumber     string `json:"address_number"`
	AddressComplement string `json:"address_complement"`
	AddressDistrict   string `json:"address_district"`
	AddressCity       string `json:"address_city"`
	AddressState      string `json:"address_state"`
	AddressCountry    string `json:"address_country"`
}

type UserMeResponse struct {
	ID                string           `json:"id"`
	FirstName         string           `json:"first_name"`
	LastName          string           `json:"last_name"`
	Email             string           `json:"email"`
	Role              string           `json:"role"`
	Plan              string           `json:"user_plan"`
	Cpf               string           `json:"cpf"`
	AvatarUrl         string           `json:"avatar_url"`
	PlanMaxSites      int              `json:"plan_max_sites"`
	PlanMaxRoutes     int              `json:"plan_max_routes"`
	DisplayName       string           `json:"display_name"`
	BirthDate         string           `json:"birth_date"`
	Gender            string           `json:"gender"`
	Bio               string           `json:"bio"`
	Instagram         string           `json:"instagram"`
	WebsiteUrl        string           `json:"website_url"`
	Whatsapp          string           `json:"whatsapp"`
	Phone             string           `json:"phone"`
	Company           *CompanyResponse `json:"company"`
	CreatedAt         string           `json:"created_at"`
	UpdatedAt         string           `json:"updated_at"`
}

type UpdateProfileRequest struct {
	FirstName         string  `json:"first_name"`
	LastName          string  `json:"last_name"`
	Cpf               *string `json:"cpf"`
	AvatarUrl         *string `json:"avatar_url"`
	DisplayName       *string `json:"display_name"`
	BirthDate         *string `json:"birth_date"`
	Gender            *string `json:"gender"`
	Bio               *string `json:"bio"`
	Instagram         *string `json:"instagram"`
	WebsiteUrl        *string `json:"website_url"`
	Whatsapp          *string `json:"whatsapp"`
	Phone             *string `json:"phone"`
}

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
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
		http.Error(w, "invalid website id", http.StatusBadRequest)
		return
	}

	var companyInput *service.CompanyInput
	if req.Company != nil {
		co := req.Company
		companyInput = &service.CompanyInput{
			CompanyName:       co.CompanyName,
			TradeName:         co.TradeName,
			Cnpj:              co.Cnpj,
			Phone:             co.Phone,
			ZipCode:           co.ZipCode,
			AddressStreet:     co.AddressStreet,
			AddressNumber:     co.AddressNumber,
			AddressComplement: co.AddressComplement,
			AddressDistrict:   co.AddressDistrict,
			AddressCity:       co.AddressCity,
			AddressState:      co.AddressState,
			AddressCountry:    co.AddressCountry,
		}
	}

	user, err := h.authService.Register(service.RegisterInput{
		WebsiteId: websiteId,
		FirstName: req.First_name,
		LastName:  req.Last_name,
		Email:     req.Email,
		Password:  req.Password,
		Cpf:       req.Cpf,
		Phone:     req.Phone,
		Company:   companyInput,
	})
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := security.SendVerifyEmail(user.Email, user.Id, user.WebsiteId); err != nil {
		log.Println("error sending verification email:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("user registered: %s %s", user.First_name, user.Last_name)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(AuthResponse{
		Success: true,
		Message: "registered, please verify your email.",
		Data: UserData{
			Id:         user.Id,
			WebsiteId:  user.WebsiteId,
			Email:      user.Email,
			Plan:       derefString(user.Plan),
			Created_at: user.Created_at,
		},
	})
}

func (h *AuthHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
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

	user, err := h.authService.Login(websiteId, req.Email, req.Password)
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
		http.Error(w, "internal error", http.StatusBadGateway)
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
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	secure := !config.IsDev
	var plan string
	if user.Plan != nil {
		plan = *user.Plan
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
	json.NewEncoder(w).Encode(AuthResponse{
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
	})
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(ResponseRefreshToken{Success: true, Message: "success"})
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
		Cpf:               derefString(user.Cpf),
		AvatarUrl:         derefString(user.AvatarUrl),
		PlanMaxSites:      planLimits[0],
		PlanMaxRoutes:     planLimits[1],
		DisplayName:       derefString(user.DisplayName),
		BirthDate:         formatDate(user.BirthDate),
		Gender:            derefString(user.Gender),
		Bio:               derefString(user.Bio),
		Instagram:         derefString(user.Instagram),
		WebsiteUrl:        derefString(user.WebsiteUrl),
		Whatsapp:          derefString(user.Whatsapp),
		Phone:             derefString(user.Phone),
		CreatedAt:         user.Created_at.Format(time.RFC3339),
		UpdatedAt:         user.Updated_at.Format(time.RFC3339),
	}
	if user.Company != nil {
		co := user.Company
		resp.Company = &CompanyResponse{
			ID:                co.Id,
			CompanyName:       co.CompanyName,
			TradeName:         derefString(co.TradeName),
			Cnpj:              derefString(co.Cnpj),
			Phone:             derefString(co.Phone),
			ZipCode:           derefString(co.ZipCode),
			AddressStreet:     derefString(co.AddressStreet),
			AddressNumber:     derefString(co.AddressNumber),
			AddressComplement: derefString(co.AddressComplement),
			AddressDistrict:   derefString(co.AddressDistrict),
			AddressCity:       derefString(co.AddressCity),
			AddressState:      derefString(co.AddressState),
			AddressCountry:    derefString(co.AddressCountry),
		}
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
	json.NewEncoder(w).Encode(map[string]any{"success": true, "message": "conta desativada; exclusão definitiva em 30 dias"})
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
		Cpf:               req.Cpf,
		AvatarUrl:         req.AvatarUrl,
		DisplayName:       req.DisplayName,
		BirthDate:         birthDate,
		Gender:            req.Gender,
		Bio:               req.Bio,
		Instagram:         req.Instagram,
		WebsiteUrl:        req.WebsiteUrl,
		Whatsapp:          req.Whatsapp,
		Phone:             req.Phone,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{"success": true, "message": "profile updated"})
}

type UpsertAddressRequest struct {
	Label      *string `json:"label"`
	ZipCode    *string `json:"zip_code"`
	Street     *string `json:"street"`
	Number     *string `json:"number"`
	Complement *string `json:"complement"`
	District   *string `json:"district"`
	City       *string `json:"city"`
	State      *string `json:"state"`
	Country    *string `json:"country"`
}

type AddressResponse struct {
	ID         string `json:"id"`
	Label      string `json:"label"`
	ZipCode    string `json:"zip_code"`
	Street     string `json:"street"`
	Number     string `json:"number"`
	Complement string `json:"complement"`
	District   string `json:"district"`
	City       string `json:"city"`
	State      string `json:"state"`
	Country    string `json:"country"`
	IsDefault  bool   `json:"is_default"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

func addressToResponse(a *domain.UserAddress) AddressResponse {
	return AddressResponse{
		ID:         a.Id,
		Label:      derefString(a.Label),
		ZipCode:    derefString(a.ZipCode),
		Street:     derefString(a.Street),
		Number:     derefString(a.Number),
		Complement: derefString(a.Complement),
		District:   derefString(a.District),
		City:       derefString(a.City),
		State:      derefString(a.State),
		Country:    a.Country,
		IsDefault:  a.IsDefault,
		CreatedAt:  a.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  a.UpdatedAt.Format(time.RFC3339),
	}
}

func (h *AuthHandler) ListAddresses(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	addrs, err := h.authService.ListAddresses(userID)
	if err != nil {
		http.Error(w, "erro ao listar endereços", http.StatusInternalServerError)
		return
	}
	resp := make([]AddressResponse, 0, len(addrs))
	for _, a := range addrs {
		resp = append(resp, addressToResponse(a))
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) CreateAddress(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	defer r.Body.Close()
	var req UpsertAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.authService.CreateAddress(userID, domain.UpsertAddressData{
		Label:      req.Label,
		ZipCode:    req.ZipCode,
		Street:     req.Street,
		Number:     req.Number,
		Complement: req.Complement,
		District:   req.District,
		City:       req.City,
		State:      req.State,
		Country:    req.Country,
	}); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{"success": false, "message": err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{"success": true})
}

func (h *AuthHandler) UpdateAddress(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	id := chi.URLParam(r, "id")
	defer r.Body.Close()
	var req UpsertAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.authService.UpdateAddress(id, userID, domain.UpsertAddressData{
		Label:      req.Label,
		ZipCode:    req.ZipCode,
		Street:     req.Street,
		Number:     req.Number,
		Complement: req.Complement,
		District:   req.District,
		City:       req.City,
		State:      req.State,
		Country:    req.Country,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"success": true})
}

func (h *AuthHandler) DeleteAddress(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	id := chi.URLParam(r, "id")
	if err := h.authService.DeleteAddress(id, userID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) SetDefaultAddress(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	id := chi.URLParam(r, "id")
	if err := h.authService.SetDefaultAddress(id, userID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"success": true})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	websiteId := r.Header.Get("X-Website-Id")
	if websiteId == "" {
		http.Error(w, "invalid website id", http.StatusBadRequest)
		return
	}

	http.SetCookie(w, &http.Cookie{Name: security.RefreshCookieName(websiteId), Value: "", Path: "/", MaxAge: -1})
	http.SetCookie(w, &http.Cookie{Name: security.AccessCookieName(websiteId), Value: "", Path: "/", MaxAge: -1})
	w.WriteHeader(http.StatusNoContent)
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
