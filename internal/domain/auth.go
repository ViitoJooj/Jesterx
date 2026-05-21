package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id             string
	WebsiteId      string
	First_name     string
	Last_name      string
	Email          string
	Verified_email bool
	Password       string
	Role           string
	Plan           *string
	Cpf            *string
	AvatarUrl      *string
	Phone          *string
	DisplayName    *string
	BirthDate      *time.Time
	Gender         *string
	Bio            *string
	Instagram      *string
	WebsiteUrl     *string
	Whatsapp       *string
	IsActive       bool
	DeactivatedAt  *time.Time
	DeleteAfter    *time.Time
	Company        *Company
	Updated_at     time.Time
	Created_at     time.Time
}

type UserAddress struct {
	Id         string
	UserId     string
	Label      *string
	ZipCode    *string
	Street     *string
	Number     *string
	Complement *string
	District   *string
	City       *string
	State      *string
	Country    string
	IsDefault  bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Company struct {
	Id                string
	OwnerUserId       string
	CompanyName       string
	TradeName         *string
	Cnpj              *string
	Phone             *string
	ZipCode           *string
	AddressStreet     *string
	AddressNumber     *string
	AddressComplement *string
	AddressDistrict   *string
	AddressCity       *string
	AddressState      *string
	AddressCountry    *string
	IsActive          bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type UpdateProfileData struct {
	FirstName   string
	LastName    string
	Cpf         *string
	AvatarUrl   *string
	Phone       *string
	DisplayName *string
	BirthDate   *time.Time
	Gender      *string
	Bio         *string
	Instagram   *string
	WebsiteUrl  *string
	Whatsapp    *string
}

type UpsertAddressData struct {
	Label      *string
	ZipCode    *string
	Street     *string
	Number     *string
	Complement *string
	District   *string
	City       *string
	State      *string
	Country    *string
}

func NewUser(websiteId, firstName, lastName, email, passwordHash string) *User {
	id, _ := uuid.NewV7()
	now := time.Now()
	return &User{
		Id:             id.String(),
		WebsiteId:      websiteId,
		First_name:     firstName,
		Last_name:      lastName,
		Email:          email,
		Verified_email: false,
		Password:       passwordHash,
		Role:           "user",
		IsActive:       true,
		Updated_at:     now,
		Created_at:     now,
	}
}

func NewCompany(ownerUserId, companyName string, tradeName, cnpj, phone, zipCode, addressStreet, addressNumber, addressComplement, addressDistrict, addressCity, addressState, addressCountry *string) *Company {
	id, _ := uuid.NewV7()
	now := time.Now()
	return &Company{
		Id:                id.String(),
		OwnerUserId:       ownerUserId,
		CompanyName:       companyName,
		TradeName:         tradeName,
		Cnpj:              cnpj,
		Phone:             phone,
		ZipCode:           zipCode,
		AddressStreet:     addressStreet,
		AddressNumber:     addressNumber,
		AddressComplement: addressComplement,
		AddressDistrict:   addressDistrict,
		AddressCity:       addressCity,
		AddressState:      addressState,
		AddressCountry:    addressCountry,
		IsActive:          true,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}
