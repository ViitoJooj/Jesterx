package validators

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	t.Parallel()

	validFirst := randomUpperLetters(t, 5)
	validLast := randomUpperLetters(t, 7)
	hyphenLeft := randomUpperLetters(t, 4)
	hyphenRight := randomUpperLetters(t, 6)
	singleName := randomUpperLetters(t, 5)
	invalidWithNumber := fmt.Sprintf("%s %s1%s", randomUpperLetters(t, 3), randomUpperLetters(t, 2), randomUpperLetters(t, 3))
	oneLetterPart := fmt.Sprintf("A %s", randomUpperLetters(t, 6))

	tests := []struct {
		name     string
		fullName string
		wantErr  bool
	}{
		{name: "valid full name", fullName: fmt.Sprintf("%s %s", validFirst, validLast), wantErr: false},
		{name: "valid hyphen name", fullName: fmt.Sprintf("%s %s-%s", validFirst, hyphenLeft, hyphenRight), wantErr: false},
		{name: "single name", fullName: singleName, wantErr: true},
		{name: "contains number", fullName: invalidWithNumber, wantErr: true},
		{name: "one letter part", fullName: oneLetterPart, wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := Name(tt.fullName)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Name(%q) error = %v, wantErr %v", tt.fullName, err, tt.wantErr)
			}
		})
	}
}
