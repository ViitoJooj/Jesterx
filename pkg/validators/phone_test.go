package validators

import (
	"strings"
	"testing"
)

func TestPhone(t *testing.T) {
	t.Parallel()

	validBR := randomValidBRMobile(t)
	repeatedDigit := strings.Repeat(string(validBR[0]), len(validBR))
	validIntl := randomValidInternationalPhone(t)

	tests := []struct {
		name    string
		ddi     int
		ddd     int
		phone   string
		wantErr bool
	}{
		{name: "valid br mobile", ddi: 55, ddd: 11, phone: validBR, wantErr: false},
		{name: "invalid ddi", ddi: 0, ddd: 11, phone: validBR, wantErr: true},
		{name: "invalid br ddd", ddi: 55, ddd: 10, phone: validBR, wantErr: true},
		{name: "invalid repeated digits", ddi: 55, ddd: 11, phone: repeatedDigit, wantErr: true},
		{name: "valid international", ddi: 1, ddd: 212, phone: validIntl, wantErr: false},
		{name: "invalid international length", ddi: 1, ddd: 212, phone: randomDigits(t, 2), wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := Phone(tt.ddi, tt.ddd, tt.phone)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Phone(%d, %d, %q) error = %v, wantErr %v", tt.ddi, tt.ddd, tt.phone, err, tt.wantErr)
			}
		})
	}
}
