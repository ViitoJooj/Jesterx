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

func (s *AuthService) Register(websiteId string, first_name string, last_name string, email string, password string) (*domain.User, error) {
	if email == "" || len(email) > 250 || len(email) < 5 || !strings.Contains(email, "@") || !strings.Contains(email, ".") || strings.Contains(email, " ") {
		return nil, errors.New("invalid email")
	}

	if password == "" || len(password) < 6 || len(password) > 50 {
		return nil, errors.New("invalid password")
	}

	webSite, err := s.webSiteRepo.FindWebSiteByID(websiteId)
	if err != nil {
		return nil, err
	}
	if webSite == nil {
		return nil, errors.New("website does not exist")
	}

	existing, err := s.userRepo.FindUserByEmailAndWebsite(email, websiteId)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	hashedPassword, err := security.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := domain.NewUser(websiteId, first_name, last_name, email, hashedPassword)

	if err := s.userRepo.UserRegister(*user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) VerifyEmail(user_id string) error {
	user, err := s.userRepo.FindUserByID(user_id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	if user.Verified_email {
		return errors.New("user already verified.")
	}

	err = s.userRepo.UpdateVerifiedEmailToTrue(user_id)
	if err != nil {
		return errors.New("Internal error")
	}

	return nil
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

func (s *AuthService) UpdateProfile(userID string, firstName string, lastName string, cpfCnpj *string, avatarUrl *string) error {
	firstName = strings.TrimSpace(firstName)
	lastName  = strings.TrimSpace(lastName)
	if len(firstName) < 1 || len(firstName) > 50 {
		return errors.New("invalid first name")
	}
	if len(lastName) < 1 || len(lastName) > 50 {
		return errors.New("invalid last name")
	}
	if cpfCnpj != nil {
		raw := strings.TrimSpace(*cpfCnpj)
		if len(raw) > 18 {
			return errors.New("invalid cpf/cnpj")
		}
		if raw == "" {
			cpfCnpj = nil
		} else {
			cpfCnpj = &raw
		}
	}
	return s.userRepo.UpdateUserProfile(userID, firstName, lastName, cpfCnpj, avatarUrl)
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
