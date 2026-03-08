package validate

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

var emailRe = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

type Error struct {
	Field   string
	Message string
}

func (e *Error) Error() string {
	if e.Field != "" {
		return e.Field + ": " + e.Message
	}
	return e.Message
}

type Errors []Error

func (e Errors) Error() string {
	msgs := make([]string, len(e))
	for i, err := range e {
		msgs[i] = err.Error()
	}
	return strings.Join(msgs, "; ")
}

func (e Errors) First() string {
	if len(e) == 0 {
		return ""
	}
	return e[0].Error()
}

// V is a validator that accumulates errors.
type V struct {
	errs Errors
}

func New() *V { return &V{} }

func (v *V) Required(field, value string) *V {
	if strings.TrimSpace(value) == "" {
		v.errs = append(v.errs, Error{field, "é obrigatório"})
	}
	return v
}

func (v *V) MaxLen(field, value string, max int) *V {
	if utf8.RuneCountInString(value) > max {
		v.errs = append(v.errs, Error{field, fmt.Sprintf("máximo %d caracteres", max)})
	}
	return v
}

func (v *V) MinLen(field, value string, minLen int) *V {
	if utf8.RuneCountInString(strings.TrimSpace(value)) < minLen {
		v.errs = append(v.errs, Error{field, fmt.Sprintf("mínimo %d caracteres", minLen)})
	}
	return v
}

func (v *V) Email(field, value string) *V {
	if value != "" && !emailRe.MatchString(value) {
		v.errs = append(v.errs, Error{field, "email inválido"})
	}
	return v
}

func (v *V) MinFloat(field string, value, minVal float64) *V {
	if value < minVal {
		v.errs = append(v.errs, Error{field, fmt.Sprintf("deve ser >= %.2f", minVal)})
	}
	return v
}

func (v *V) MinInt(field string, value, minVal int) *V {
	if value < minVal {
		v.errs = append(v.errs, Error{field, fmt.Sprintf("deve ser >= %d", minVal)})
	}
	return v
}

func (v *V) Err() error {
	if len(v.errs) == 0 {
		return nil
	}
	return v.errs
}

func (v *V) First() string {
	return v.errs.First()
}
