package service

import (
	"errors"
	"strings"
	"unicode"

	"github.com/ViitoJooj/Jesterx/internal/domain"
	"github.com/ViitoJooj/Jesterx/internal/repository"
)

type WebSiteService struct {
	webSiteRepo repository.WebsiteRepository
}

func NewWebSiteService(webSiteRepo repository.WebsiteRepository) *WebSiteService {
	return &WebSiteService{
		webSiteRepo: webSiteRepo,
	}
}

var acceptedTypes = [4]string{"ECOMMERCE", "LANDING_PAGE", "SOFTWARE_SELL", "COURSE"}

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

func (s *WebSiteService) CreateWebSite(Type string, Image []byte, Name string, Short_description string, Description string, Creator_id string) (*domain.WebSite, error) {
	Type = strings.ToUpper(strings.TrimSpace(Type))
	Name = strings.TrimSpace(Name)

	if !isValidType(Type) {
		return nil, errors.New("invalid type")
	}

	if len(Name) > 50 || containsInvalidChars(Name) == true {
		return nil, errors.New("Invalid name.")
	}

	existing, err := s.webSiteRepo.FindWebSiteByName(Name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("This site already exists.")
	}

	website := domain.NewWebSite(Type, Image, Name, Short_description, Description, Creator_id)
	if err := s.webSiteRepo.SaveWebSite(*website); err != nil {
		return nil, err
	}

	return website, nil
}
