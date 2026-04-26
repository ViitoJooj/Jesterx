package validators

import (
	"testing"

	"github.com/google/uuid"
)

func TestUuid(t *testing.T) {
	t.Parallel()

	validV7, err := uuid.NewV7()
	if err != nil {
		t.Fatalf("uuid.NewV7() error = %v", err)
	}

	v4 := uuid.New()

	tests := []struct {
		name    string
		value   any
		wantErr bool
	}{
		{name: "valid v7 uuid string", value: validV7.String(), wantErr: false},
		{name: "invalid type", value: 123, wantErr: true},
		{name: "invalid format", value: "not-an-uuid", wantErr: true},
		{name: "wrong version", value: v4.String(), wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := Uuid(tt.value)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Uuid(%v) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
		})
	}
}
