package helpers

import (
	"testing"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{"Valid email", "test@example.com", false},
		{"Valid email with subdomain", "user@mail.example.com", false},
		{"Empty email", "", true},
		{"Too short", "a@b.c", true},
		{"Too long", "a" + string(make([]byte, 200)) + "@example.com", true},
		{"Missing @", "testexample.com", true},
		{"Multiple @", "test@@example.com", true},
		{"No domain", "test@", true},
		{"No local part", "@example.com", true},
		{"Domain starts with dot", "test@.example.com", true},
		{"Domain ends with dot", "test@example.com.", true},
		{"No dot in domain", "test@example", true},
		{"Invalid character", "test@exam ple.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"Valid password", "Password123!", false},
		{"Valid password complex", "MyP@ssw0rd", false},
		{"Empty password", "", true},
		{"Too short", "Pass1!", true},
		{"Too long", "P@ssw0rd" + string(make([]byte, 100)), true},
		{"No uppercase", "password123!", true},
		{"No lowercase", "PASSWORD123!", true},
		{"No digit", "Password!@#", true},
		{"No special char", "Password123", true},
		{"Invalid character", "Password123!ñ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
