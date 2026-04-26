package validators

import "testing"

func TestPassword(t *testing.T) {
	t.Parallel()

	validPassword := randomValidPassword(t)
	tooSmallPassword := randomUpperAlnum(t, 7)
	tooLargePassword := randomUpperAlnum(t, 51)
	invalidCharPassword := randomUpperAlnum(t, 5) + "(" + randomUpperAlnum(t, 3)

	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{name: "valid password", password: validPassword, wantErr: false},
		{name: "too small", password: tooSmallPassword, wantErr: true},
		{name: "too large", password: tooLargePassword, wantErr: true},
		{name: "contains invalid char", password: invalidCharPassword, wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := Password(tt.password)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Password(%q) error = %v, wantErr %v", tt.password, err, tt.wantErr)
			}
		})
	}
}
