package validators

import "testing"

func TestRole(t *testing.T) {
	t.Parallel()

	invalidRole := randomUpperLetters(t, 8)

	tests := []struct {
		name    string
		role    string
		wantErr bool
	}{
		{name: "valid admin", role: "admin", wantErr: false},
		{name: "valid support", role: "support", wantErr: false},
		{name: "invalid role", role: invalidRole, wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := Role(tt.role)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Role(%q) error = %v, wantErr %v", tt.role, err, tt.wantErr)
			}
		})
	}
}
