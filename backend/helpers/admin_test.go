package helpers

import (
	"testing"

	"jesterx-core/config"
)

func TestIsPlatformAdminEmail(t *testing.T) {
	config.AdminEmails = []string{"admin@example.com", " root@Example.com "}

	if !IsPlatformAdminEmail("admin@example.com") {
		t.Fatalf("expected admin@example.com to be admin")
	}

	if !IsPlatformAdminEmail("ROOT@example.com") {
		t.Fatalf("expected case-insensitive match")
	}

	if IsPlatformAdminEmail("user@example.com") {
		t.Fatalf("non admin email detected as admin")
	}
}

func TestResolvePlatformRole(t *testing.T) {
	config.AdminEmails = []string{"boss@example.com"}

	if got := ResolvePlatformRole("boss@example.com", "platform_user"); got != "platform_admin" {
		t.Fatalf("expected platform_admin, got %s", got)
	}

	if got := ResolvePlatformRole("user@example.com", "customer"); got != "customer" {
		t.Fatalf("expected role passthrough, got %s", got)
	}

	if got := ResolvePlatformRole("user@example.com", ""); got != "platform_user" {
		t.Fatalf("expected default platform_user, got %s", got)
	}
}
