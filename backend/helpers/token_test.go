package helpers

import (
	"gen-you-ecommerce/config"
	"testing"
)

func init() {
	// Set a test JWT secret for testing
	config.JwtSecret = "test-secret-key-for-unit-tests"
}

func TestGenerateToken(t *testing.T) {
	user := UserData{
		Id:          "test-user-id",
		Profile_img: "test.jpg",
		First_name:  "John",
		Last_name:   "Doe",
		Email:       "john@example.com",
		Role:        "user",
		Plan:        "free",
	}
	
	token, err := GenerateToken(user)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}
	
	if token == "" {
		t.Error("GenerateToken() returned empty token")
	}
}

func TestValidateToken(t *testing.T) {
	user := UserData{
		Id:          "test-user-id",
		Profile_img: "test.jpg",
		First_name:  "John",
		Last_name:   "Doe",
		Email:       "john@example.com",
		Role:        "user",
		Plan:        "free",
	}
	
	token, err := GenerateToken(user)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}
	
	claims, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken() error = %v", err)
	}
	
	if claims["sub"] != user.Id {
		t.Errorf("ValidateToken() sub = %v, want %v", claims["sub"], user.Id)
	}
	
	if claims["email"] != user.Email {
		t.Errorf("ValidateToken() email = %v, want %v", claims["email"], user.Email)
	}
	
	if claims["role"] != user.Role {
		t.Errorf("ValidateToken() role = %v, want %v", claims["role"], user.Role)
	}
}

func TestValidateTokenInvalid(t *testing.T) {
	tests := []struct {
		name  string
		token string
	}{
		{"Empty token", ""},
		{"Invalid token", "invalid.token.here"},
		{"Malformed token", "not-a-jwt-token"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ValidateToken(tt.token)
			if err == nil {
				t.Error("ValidateToken() expected error for invalid token")
			}
		})
	}
}

func TestGetLoginDuration(t *testing.T) {
	tests := []struct {
		name           string
		keepMeLoggedIn bool
		want           int
	}{
		{"Temporary login", false, 24},
		{"Persistent login", true, 744},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetLoginDuration(tt.keepMeLoggedIn); got != tt.want {
				t.Errorf("GetLoginDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}
