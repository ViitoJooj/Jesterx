package service

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ViitoJooj/Jesterx/internal/domain"
	"github.com/ViitoJooj/Jesterx/internal/repository"
	"github.com/ViitoJooj/Jesterx/internal/security"
	"github.com/google/uuid"
)

type AuthService struct {
	userRepo    repository.UserRepository
	webSiteRepo repository.WebsiteRepository
	paymentRepo repository.PaymentRepository
}

func (s *AuthService) GetUserByID(userID string) (*domain.User, error) {
	return s.userRepo.FindUserByID(userID)
}

func NewAuthService(userRepo repository.UserRepository, webSiteRepo repository.WebsiteRepository, paymentRepo repository.PaymentRepository) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		webSiteRepo: webSiteRepo,
		paymentRepo: paymentRepo,
	}
}

func (s *AuthService) DeleteAccount(userID string) error {
	user, err := s.userRepo.FindUserByID(userID)
	if err != nil {
		return errors.New("erro interno")
	}
	if user == nil {
		return errors.New("usuário não encontrado")
	}
	deleteAfter := time.Now().Add(30 * 24 * time.Hour)
	return s.userRepo.DeactivateUserByID(userID, deleteAfter)
}

func (s *AuthService) DeleteExpiredUnverifiedUsers() error {
	return s.userRepo.DeleteExpiredUnverifiedUsers()
}

func (s *AuthService) DeleteExpiredDeactivatedUsers() error {
	return s.userRepo.DeleteExpiredDeactivatedUsers()
}

type CompanyInput struct {
	CompanyName     string
	TradeName       *string
	Cnpj            *string
	Phone           *string
	ZipCode         *string
	AddressStreet   *string
	AddressNumber   *string
	AddressComplement *string
	AddressDistrict *string
	AddressCity     *string
	AddressState    *string
	AddressCountry  *string
}

type RegisterInput struct {
	WebsiteId string
	FirstName string
	LastName  string
	Email     string
	Password  string
	Cpf       *string
	Phone     *string
	Company   *CompanyInput
}

func (s *AuthService) Register(input RegisterInput) (*domain.User, error) {
	if input.Email == "" || len(input.Email) > 250 || len(input.Email) < 5 || !strings.Contains(input.Email, "@") || !strings.Contains(input.Email, ".") || strings.Contains(input.Email, " ") {
		return nil, errors.New("invalid email")
	}

	if input.Password == "" || len(input.Password) < 6 || len(input.Password) > 50 {
		return nil, errors.New("invalid password")
	}

	if input.Company != nil {
		co := input.Company
		if strings.TrimSpace(co.CompanyName) == "" {
			return nil, errors.New("company name is required")
		}
		if co.Phone == nil || strings.TrimSpace(*co.Phone) == "" {
			return nil, errors.New("company phone is required")
		}
		if co.ZipCode == nil || strings.TrimSpace(*co.ZipCode) == "" {
			return nil, errors.New("company zip code is required")
		}
		if co.AddressStreet == nil || strings.TrimSpace(*co.AddressStreet) == "" {
			return nil, errors.New("company address is required")
		}
		if co.AddressCity == nil || strings.TrimSpace(*co.AddressCity) == "" {
			return nil, errors.New("company city is required")
		}
		if co.AddressState == nil || strings.TrimSpace(*co.AddressState) == "" {
			return nil, errors.New("company state is required")
		}
	}

	webSite, err := s.webSiteRepo.FindWebSiteByID(input.WebsiteId)
	if err != nil {
		return nil, err
	}
	if webSite == nil {
		return nil, errors.New("website does not exist")
	}

	existing, err := s.userRepo.FindUserByEmailAndWebsite(input.Email, input.WebsiteId)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	hashedPassword, err := security.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	user := domain.NewUser(input.WebsiteId, input.FirstName, input.LastName, input.Email, hashedPassword)
	user.Cpf = input.Cpf
	user.Phone = input.Phone

	if err := s.userRepo.UserRegister(*user); err != nil {
		return nil, err
	}

	if input.Company != nil {
		co := input.Company
		country := co.AddressCountry
		if country == nil {
			br := "BR"
			country = &br
		}
		company := domain.NewCompany(
			user.Id, co.CompanyName, co.TradeName, co.Cnpj, co.Phone,
			co.ZipCode, co.AddressStreet, co.AddressNumber, co.AddressComplement,
			co.AddressDistrict, co.AddressCity, co.AddressState, country,
		)
		if err := s.userRepo.CompanyRegister(*company); err != nil {
			return nil, err
		}
		user.Company = company
	}

	return user, nil
}

func (s *AuthService) VerifyEmail(user_id string, websiteID string) (*domain.User, error) {
	user, err := s.userRepo.FindUserByID(user_id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	if websiteID != "" && user.WebsiteId != websiteID {
		return nil, errors.New("user not found")
	}

	if user.Verified_email {
		return nil, errors.New("user already verified.")
	}

	if websiteID != "" {
		err = s.userRepo.UpdateVerifiedEmailToTrueByWebsite(user_id, websiteID)
	} else {
		err = s.userRepo.UpdateVerifiedEmailToTrue(user_id)
	}
	if err != nil {
		return nil, errors.New("Internal error")
	}

	return user, nil
}

func (s *AuthService) Login(websiteId string, email string, password string) (*domain.User, error) {
	webSite, err := s.webSiteRepo.FindWebSiteByID(websiteId)
	if err != nil {
		log.Println(err)
		return nil, errors.New("Internal error")
	}
	if webSite == nil {
		log.Println("Website no exists")
		return nil, errors.New("website does not exist")
	}

	user, err := s.userRepo.FindUserByEmailAndWebsite(email, websiteId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		log.Println("User not exists")
		return nil, errors.New("invalid credentials")
	}

	if !security.CheckPassword(password, user.Password) {
		log.Println("Incorrect password")
		return nil, errors.New("invalid credentials")
	}

	if !user.Verified_email {
		log.Println("Email is not verified")
		return nil, errors.New("Email is not verified")
	}

	return user, nil
}

func (s *AuthService) Refresh(refreshToken string) (string, error) {
	refreshClaims, err := security.ParseRefreshToken(refreshToken)
	if err != nil {
		return "", errors.New("invalid refresh token")
	}

	user, err := s.userRepo.FindUserByID(refreshClaims.Sub)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("user not found")
	}

	if user.WebsiteId != refreshClaims.WebsiteId {
		return "", errors.New("invalid token context")
	}

	webSite, err := s.webSiteRepo.FindWebSiteByID(refreshClaims.WebsiteId)
	if err != nil {
		return "", err
	}
	if webSite == nil {
		return "", errors.New("website does not exist")
	}
	if webSite.Banned {
		return "", errors.New("website is banned")
	}

	if !user.Verified_email {
		log.Println("Email is not verified")
		return "", errors.New("Email is not verified")
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
		return "", errors.New("failed to generate access token")
	}

	return accessToken, nil
}

func (s *AuthService) Me(userID string) (*domain.User, error) {
	user, err := s.userRepo.FindUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	if !user.Verified_email {
		return nil, errors.New("Email is not verified")
	}
	return user, nil
}

func (s *AuthService) MeWithPlan(userID string) (*domain.User, [2]int, error) {
	user, err := s.Me(userID)
	if err != nil {
		return nil, [2]int{}, err
	}
	limits := [2]int{1, 5}
	if user.Plan != nil && *user.Plan != "" {
		plan, err := s.paymentRepo.FindPlanByName(*user.Plan)
		if err == nil && plan != nil {
			limits = [2]int{plan.MaxSites, plan.MaxRoutes}
		}
	}
	return user, limits, nil
}

func (s *AuthService) UpdateProfile(userID string, data domain.UpdateProfileData) error {
	data.FirstName = strings.TrimSpace(data.FirstName)
	data.LastName = strings.TrimSpace(data.LastName)
	if len(data.FirstName) < 1 || len(data.FirstName) > 50 {
		return errors.New("invalid first name")
	}
	if len(data.LastName) < 1 || len(data.LastName) > 50 {
		return errors.New("invalid last name")
	}
	if data.Cpf != nil {
		raw := strings.TrimSpace(*data.Cpf)
		if len(raw) > 14 {
			return errors.New("CPF inválido")
		}
		if raw == "" {
			data.Cpf = nil
		} else {
			data.Cpf = &raw
		}
	}
	if data.DisplayName != nil {
		raw := strings.TrimSpace(*data.DisplayName)
		if raw == "" {
			data.DisplayName = nil
		} else if len(raw) > 100 {
			return errors.New("display_name muito longo")
		} else {
			data.DisplayName = &raw
		}
	}
	if data.Bio != nil {
		raw := strings.TrimSpace(*data.Bio)
		if raw == "" {
			data.Bio = nil
		} else if len(raw) > 500 {
			return errors.New("bio muito longa")
		} else {
			data.Bio = &raw
		}
	}
	if data.Gender != nil {
		raw := strings.ToLower(strings.TrimSpace(*data.Gender))
		if raw == "" {
			data.Gender = nil
		} else if raw != "male" && raw != "female" && raw != "other" && raw != "prefer_not" {
			return errors.New("gênero inválido")
		} else {
			data.Gender = &raw
		}
	}
	if data.BirthDate != nil && data.BirthDate.After(time.Now()) {
		return errors.New("data de nascimento inválida")
	}
	if data.Instagram != nil {
		raw := strings.TrimSpace(*data.Instagram)
		if raw == "" {
			data.Instagram = nil
		} else if len(raw) > 100 {
			return errors.New("instagram inválido")
		} else {
			data.Instagram = &raw
		}
	}
	if data.WebsiteUrl != nil {
		raw := strings.TrimSpace(*data.WebsiteUrl)
		if raw == "" {
			data.WebsiteUrl = nil
		} else if len(raw) > 200 {
			return errors.New("site inválido")
		} else {
			data.WebsiteUrl = &raw
		}
	}
	if data.Whatsapp != nil {
		raw := strings.TrimSpace(*data.Whatsapp)
		if raw == "" {
			data.Whatsapp = nil
		} else if len(raw) > 20 {
			return errors.New("whatsapp inválido")
		} else {
			data.Whatsapp = &raw
		}
	}
	return s.userRepo.UpdateUserProfile(userID, data)
}

func (s *AuthService) ListAddresses(userID string) ([]*domain.UserAddress, error) {
	return s.userRepo.ListUserAddresses(userID)
}

func (s *AuthService) CreateAddress(userID string, data domain.UpsertAddressData) error {
	addrs, err := s.userRepo.ListUserAddresses(userID)
	if err != nil {
		return errors.New("erro interno")
	}
	if len(addrs) >= 10 {
		return errors.New("limite de 10 endereços atingido")
	}
	country := "BR"
	if data.Country != nil && strings.TrimSpace(*data.Country) != "" {
		country = strings.TrimSpace(*data.Country)
	}
	id, _ := uuid.NewV7()
	addr := domain.UserAddress{
		Id:         id.String(),
		UserId:     userID,
		Label:      data.Label,
		ZipCode:    data.ZipCode,
		Street:     data.Street,
		Number:     data.Number,
		Complement: data.Complement,
		District:   data.District,
		City:       data.City,
		State:      data.State,
		Country:    country,
		IsDefault:  len(addrs) == 0,
	}
	return s.userRepo.CreateUserAddress(addr)
}

func (s *AuthService) UpdateAddress(id, userID string, data domain.UpsertAddressData) error {
	return s.userRepo.UpdateUserAddress(id, userID, data)
}

func (s *AuthService) DeleteAddress(id, userID string) error {
	return s.userRepo.DeleteUserAddress(id, userID)
}

func (s *AuthService) SetDefaultAddress(id, userID string) error {
	return s.userRepo.SetDefaultAddress(id, userID)
}

func (s *AuthService) Logout(w http.ResponseWriter) error {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	return nil
}
