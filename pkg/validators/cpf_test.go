package validators

import "testing"

func TestCpf(t *testing.T) {
	t.Parallel()

	validCPF := randomValidCPF(t)
	invalidCheckDigits := mutateLastDigit(validCPF)
	invalidLength := randomDigits(t, 3)
	allSame := randomAllSameDigits(t, 11)

	tests := []struct {
		name    string
		cpf     string
		wantErr bool
	}{
		{name: "valid cpf with punctuation", cpf: validCPF, wantErr: false},
		{name: "invalid length", cpf: invalidLength, wantErr: true},
		{name: "all same digits", cpf: allSame, wantErr: true},
		{name: "invalid check digits", cpf: invalidCheckDigits, wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := Cpf(tt.cpf)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Cpf(%q) error = %v, wantErr %v", tt.cpf, err, tt.wantErr)
			}
		})
	}
}
