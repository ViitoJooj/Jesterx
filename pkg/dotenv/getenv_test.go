package dotenv

import "testing"

func TestGet(t *testing.T) {
	t.Run("returns env value", func(t *testing.T) {
		t.Setenv("TEST_ENV_NAME", "value")

		got := get("TEST_ENV_NAME")
		if got != "value" {
			t.Fatalf("get returned %q, want %q", got, "value")
		}
	})

	t.Run("panics when env is missing", func(t *testing.T) {
		t.Setenv("TEST_ENV_MISSING", "")

		defer func() {
			if recover() == nil {
				t.Fatalf("expected panic for missing env")
			}
		}()

		_ = get("TEST_ENV_MISSING")
	})
}
