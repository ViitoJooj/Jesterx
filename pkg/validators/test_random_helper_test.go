package validators

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"testing"
)

func randomInt64(t *testing.T, max int64) int64 {
	t.Helper()

	if max <= 0 {
		t.Fatalf("max must be > 0")
	}

	v, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		t.Fatalf("failed generating random int: %v", err)
	}

	return v.Int64()
}

func randomFromCharset(t *testing.T, charset string, n int) string {
	t.Helper()

	if n <= 0 {
		t.Fatalf("n must be > 0")
	}
	if len(charset) == 0 {
		t.Fatalf("charset must not be empty")
	}

	b := make([]byte, n)
	for i := range b {
		b[i] = charset[randomInt64(t, int64(len(charset)))]
	}

	return string(b)
}

func randomUpperLetters(t *testing.T, n int) string {
	t.Helper()
	return randomFromCharset(t, "ABCDEFGHIJKLMNOPQRSTUVWXYZ", n)
}

func randomDigits(t *testing.T, n int) string {
	t.Helper()
	return randomFromCharset(t, "0123456789", n)
}

func randomUpperAlnum(t *testing.T, n int) string {
	t.Helper()
	return randomFromCharset(t, "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", n)
}

func randomAllSameDigits(t *testing.T, n int) string {
	t.Helper()
	d := randomDigits(t, 1)
	return strings.Repeat(d, n)
}

func randomValidEmail(t *testing.T) string {
	t.Helper()
	local := randomUpperLetters(t, 6) + randomDigits(t, 3)
	domain := randomUpperLetters(t, 5)
	tld := randomUpperLetters(t, 3)
	return fmt.Sprintf("%s@%s.%s", local, domain, tld)
}

func randomValidPassword(t *testing.T) string {
	t.Helper()
	return fmt.Sprintf("%s%s!", randomUpperLetters(t, 6), randomDigits(t, 4))
}

func randomValidBRMobile(t *testing.T) string {
	t.Helper()
	return "9" + randomDigits(t, 8)
}

func randomValidInternationalPhone(t *testing.T) string {
	t.Helper()
	return fmt.Sprintf("(%s) %s-%s", randomDigits(t, 3), randomDigits(t, 3), randomDigits(t, 4))
}

func randomValidCPF(t *testing.T) string {
	t.Helper()

	base := randomDigits(t, 9)
	if allDigitsEqual(base) {
		base = "123456789"
	}

	digits := make([]int, 11)
	for i := 0; i < 9; i++ {
		digits[i] = int(base[i] - '0')
	}

	sum := 0
	for i := 0; i < 9; i++ {
		sum += digits[i] * (10 - i)
	}
	d1 := (sum * 10) % 11
	if d1 == 10 {
		d1 = 0
	}
	digits[9] = d1

	sum = 0
	for i := 0; i < 10; i++ {
		sum += digits[i] * (11 - i)
	}
	d2 := (sum * 10) % 11
	if d2 == 10 {
		d2 = 0
	}
	digits[10] = d2

	return fmt.Sprintf("%d%d%d.%d%d%d.%d%d%d-%d%d",
		digits[0], digits[1], digits[2],
		digits[3], digits[4], digits[5],
		digits[6], digits[7], digits[8],
		digits[9], digits[10],
	)
}

func mutateLastDigit(s string) string {
	if len(s) == 0 {
		return s
	}
	last := s[len(s)-1]
	if last == '9' {
		return s[:len(s)-1] + "0"
	}
	if last >= '0' && last <= '8' {
		return s[:len(s)-1] + string(last+1)
	}
	return s[:len(s)-1] + "0"
}

func allDigitsEqual(s string) bool {
	if len(s) == 0 {
		return true
	}
	for i := 1; i < len(s); i++ {
		if s[i] != s[0] {
			return false
		}
	}
	return true
}
