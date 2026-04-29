package domain

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestNewUserCurrentBehavior(t *testing.T) {
	t.Parallel()

	suffix := strings.ToUpper(strings.ReplaceAll(uuid.NewString(), "-", ""))
	name := fmt.Sprintf("%s %s", randomUpperLetters(t, 6), randomUpperLetters(t, 8))
	email := fmt.Sprintf("USER%s@TEST.COM", suffix)
	password := fmt.Sprintf("STRONG%s!", suffix[:8])
	cpf := "52998224725"

	user, err := NewUser(name, email, password, "admin", cpf)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}

	if user == nil {
		t.Fatal("expected user, got nil")
	}

	if user.Email != email {
		t.Fatalf("expected email %q, got %q", email, user.Email)
	}
}

func randomUpperLetters(t *testing.T, n int) string {
	t.Helper()

	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		v, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			t.Fatalf("failed generating random letters: %v", err)
		}
		b[i] = letters[v.Int64()]
	}

	return string(b)
}
