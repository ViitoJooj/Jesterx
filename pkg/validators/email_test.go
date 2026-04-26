package validators

import (
	"strings"
	"testing"
)

func TestEmail(t *testing.T) {
	t.Parallel()

	validEmail := randomValidEmail(t)
	tooLargeEmail := strings.Repeat(randomUpperLetters(t, 1), 245) + "@A.COM"
	containsLower := strings.ToLower(validEmail[:1]) + validEmail[1:]
	containsSpace := strings.Replace(validEmail, "@", " @", 1)
	missingAt := strings.Replace(validEmail, "@", "", 1)
	missingDot := strings.Replace(validEmail, ".", "", 1)
	tooSmallEmail := randomUpperLetters(t, 1) + "@" + randomUpperLetters(t, 1)

	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{name: "valid uppercase email", email: validEmail, wantErr: false},
		{name: "too large", email: tooLargeEmail, wantErr: true},
		{name: "too small", email: tooSmallEmail, wantErr: true},
		{name: "missing at", email: missingAt, wantErr: true},
		{name: "missing dot", email: missingDot, wantErr: true},
		{name: "contains spaces", email: containsSpace, wantErr: true},
		{name: "contains lowercase", email: containsLower, wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := Email(tt.email)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Email(%q) error = %v, wantErr %v", tt.email, err, tt.wantErr)
			}
		})
	}
}

func TestLoadEmbedded(t *testing.T) {
	t.Parallel()

	blocklist, err := LoadEmbedded()
	if err != nil {
		t.Fatalf("LoadEmbedded() error = %v", err)
	}

	if len(blocklist.domains) == 0 {
		t.Fatalf("LoadEmbedded() returned empty domain blocklist")
	}
}

func TestContainsCapitalLetters(t *testing.T) {
	t.Parallel()

	allUpper := randomUpperLetters(t, 8)
	upperWithDigitsAndSymbols := randomUpperLetters(t, 1) + randomDigits(t, 2) + "@" + randomUpperLetters(t, 1)
	containsLower := strings.ToLower(allUpper[:1]) + allUpper[1:]

	tests := []struct {
		name string
		s    string
		want bool
	}{
		{name: "all letters uppercase", s: allUpper, want: true},
		{name: "uppercase with digits and symbols", s: upperWithDigitsAndSymbols, want: true},
		{name: "contains lowercase", s: containsLower, want: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := containsCapitalLetters(tt.s)
			if got != tt.want {
				t.Fatalf("containsCapitalLetters(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}
