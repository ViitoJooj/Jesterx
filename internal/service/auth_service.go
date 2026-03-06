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

// function used for delete not verifieds users in 10 minutes
func (s *AuthService) DeleteExpiredUnverifiedUsers() error {
	return s.userRepo.DeleteExpiredUnverifiedUsers()
}

type RegisterInput struct {
	WebsiteId         string
	FirstName         string
	LastName          string
	Email             string
	Password          string
	AccountType       string
	CompanyName       *string
	TradeName         *string
	CpfCnpj           *string
	Phone             *string
	ZipCode           *string
	AddressStreet     *string
	AddressNumber     *string
	AddressComplement *string
	AddressCity       *string
	AddressState      *string
}

func (s *AuthService) Register(input RegisterInput) (*domain.User, error) {
	if input.Email == "" || len(input.Email) > 250 || len(input.Email) < 5 || !strings.Contains(input.Email, "@") || !strings.Contains(input.Email, ".") || strings.Contains(input.Email, " ") {
		return nil, errors.New("invalid email")
	}

	if input.Password == "" || len(input.Password) < 6 || len(input.Password) > 50 {
		return nil, errors.New("invalid password")
	}

	if input.AccountType != "personal" && input.AccountType != "business" {
		input.AccountType = "personal"
	}

	if input.AccountType == "business" {
		if input.CompanyName == nil || strings.TrimSpace(*input.CompanyName) == "" {
			return nil, errors.New("company name is required for business accounts")
		}
		if input.CpfCnpj == nil || strings.TrimSpace(*input.CpfCnpj) == "" {
			return nil, errors.New("CNPJ is required for business accounts")
		}
		if input.Phone == nil || strings.TrimSpace(*input.Phone) == "" {
			return nil, errors.New("phone is required for business accounts")
		}
		if input.ZipCode == nil || strings.TrimSpace(*input.ZipCode) == "" {
			return nil, errors.New("zip code is required for business accounts")
		}
		if input.AddressStreet == nil || strings.TrimSpace(*input.AddressStreet) == "" {
			return nil, errors.New("address is required for business accounts")
		}
		if input.AddressCity == nil || strings.TrimSpace(*input.AddressCity) == "" {
			return nil, errors.New("city is required for business accounts")
		}
		if input.AddressState == nil || strings.TrimSpace(*input.AddressState) == "" {
			return nil, errors.New("state is required for business accounts")
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

	user := domain.NewUser(input.WebsiteId, input.FirstName, input.LastName, input.Email, hashedPassword, input.AccountType)
	user.CpfCnpj = input.CpfCnpj
	user.CompanyName = input.CompanyName
	user.TradeName = input.TradeName
	user.Phone = input.Phone
	user.ZipCode = input.ZipCode
	user.AddressStreet = input.AddressStreet
	user.AddressNumber = input.AddressNumber
	user.AddressComplement = input.AddressComplement
	user.AddressCity = input.AddressCity
	user.AddressState = input.AddressState
	if input.AccountType == "business" {
		country := "Brasil"
		user.AddressCountry = &country
	}

	if err := s.userRepo.UserRegister(*user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) VerifyEmail(user_id string) (*domain.User, error) {
	user, err := s.userRepo.FindUserByID(user_id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	if user.Verified_email {
		return nil, errors.New("user already verified.")
	}

	err = s.userRepo.UpdateVerifiedEmailToTrue(user_id)
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

// MeWithPlan returns the user and [maxSites, maxRoutes] from their active plan.
func (s *AuthService) MeWithPlan(userID string) (*domain.User, [2]int, error) {
	user, err := s.Me(userID)
	if err != nil {
		return nil, [2]int{}, err
	}
	limits := [2]int{1, 5} // conservative defaults
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
	data.LastName  = strings.TrimSpace(data.LastName)
	if len(data.FirstName) < 1 || len(data.FirstName) > 50 {
		return errors.New("invalid first name")
	}
	if len(data.LastName) < 1 || len(data.LastName) > 50 {
		return errors.New("invalid last name")
	}
	if data.CpfCnpj != nil {
		raw := strings.TrimSpace(*data.CpfCnpj)
		if len(raw) > 18 {
			return errors.New("invalid cpf/cnpj")
		}
		if raw == "" {
			data.CpfCnpj = nil
		} else {
			data.CpfCnpj = &raw
		}
	}
	return s.userRepo.UpdateUserProfile(userID, data)
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
