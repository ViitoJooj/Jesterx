package helpers

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "TestPassword123!"
	
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}
	
	if hash == "" {
		t.Error("HashPassword() returned empty hash")
	}
	
	if hash == password {
		t.Error("HashPassword() returned unhashed password")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "TestPassword123!"
	wrongPassword := "WrongPassword123!"
	
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}
	
	tests := []struct {
		name     string
		password string
		hash     string
		want     bool
	}{
		{"Correct password", password, hash, true},
		{"Wrong password", wrongPassword, hash, false},
		{"Empty password", "", hash, false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckPasswordHash(tt.password, tt.hash); got != tt.want {
				t.Errorf("CheckPasswordHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHashPasswordConsistency(t *testing.T) {
	password := "TestPassword123!"
	
	hash1, err1 := HashPassword(password)
	if err1 != nil {
		t.Fatalf("HashPassword() error = %v", err1)
	}
	
	hash2, err2 := HashPassword(password)
	if err2 != nil {
		t.Fatalf("HashPassword() error = %v", err2)
	}
	
	// Hashes should be different (bcrypt includes random salt)
	if hash1 == hash2 {
		t.Error("HashPassword() returned identical hashes for same password (should be different due to salt)")
	}
	
	// But both should validate correctly
	if !CheckPasswordHash(password, hash1) {
		t.Error("CheckPasswordHash() failed for hash1")
	}
	
	if !CheckPasswordHash(password, hash2) {
		t.Error("CheckPasswordHash() failed for hash2")
	}
}
