package validators

import (
	"bufio"
	"embed"
	"errors"
	"fmt"
	"log"
	"strings"
	"unicode"
)

type Blocklist struct {
	domains map[string]struct{}
}

func Email(email string) error {
	if len(email) > 250 {
		return errors.New("Email is too large.")
	}

	if len(email) < 5 {
		return errors.New("Email is too small.")
	}

	if !strings.Contains(email, "@") {
		return errors.New("Email is not valid")
	}

	if !strings.Contains(email, ".") {
		return errors.New("Email is not valid")
	}

	if strings.Contains(email, " ") {
		return errors.New("Email cannot contain spaces.")
	}

	blocklist, err := LoadEmbedded()
	if err != nil {
		log.Fatal(err)
	}

	domain := strings.Split(email, "@")[1]
	if blocklist.isBlockedDomain(domain) {
		fmt.Println("Email bloqueado")
	}

	if !containsCapitalLetters(email) {
		return errors.New("Email addresses cannot contain capital letters.")
	}

	return nil
}

//go:embed utils/bannedDomains.txt
var domainFile embed.FS

func LoadEmbedded() (*Blocklist, error) {
	f, err := domainFile.Open("utils/bannedDomains.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	domains := make(map[string]struct{}, 6000)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		domains[strings.ToLower(line)] = struct{}{}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &Blocklist{domains: domains}, nil
}

func (b *Blocklist) isBlockedDomain(domain string) bool {
	domain = strings.ToLower(strings.TrimSpace(domain))

	for {
		if _, ok := b.domains[domain]; ok {
			return true
		}

		i := strings.Index(domain, ".")
		if i == -1 {
			break
		}

		domain = domain[i+1:]
	}

	return false
}

func containsCapitalLetters(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) && !unicode.IsUpper(r) {
			return false
		}
	}
	return true
}
