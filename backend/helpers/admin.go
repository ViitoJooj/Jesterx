package helpers

import (
	"strings"

	"jesterx-core/config"
)

func IsPlatformAdminEmail(email string) bool {
	normalized := strings.ToLower(strings.TrimSpace(email))
	if normalized == "" {
		return false
	}
	for _, allowed := range config.AdminEmails {
		if strings.ToLower(strings.TrimSpace(allowed)) == normalized {
			return true
		}
	}
	return false
}

func ResolvePlatformRole(email string, current string) string {
	if IsPlatformAdminEmail(email) {
		return "platform_admin"
	}
	if strings.TrimSpace(current) != "" {
		return current
	}
	return "platform_user"
}
