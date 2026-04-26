package validators

import "testing"

func TestPlan(t *testing.T) {
	t.Parallel()

	unknownPlan := randomUpperLetters(t, 10)

	tests := []struct {
		name string
		plan string
	}{
		{name: "known plan", plan: "free"},
		{name: "unknown plan still returns nil", plan: unknownPlan},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := Plan(tt.plan); err != nil {
				t.Fatalf("Plan(%q) returned unexpected error: %v", tt.plan, err)
			}
		})
	}
}
