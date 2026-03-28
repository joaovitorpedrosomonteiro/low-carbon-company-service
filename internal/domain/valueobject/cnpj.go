package valueobject

import (
	"errors"
	"strings"
)

var (
	ErrInvalidCNPJ = errors.New("invalid CNPJ format")
)

type CNPJ struct {
	value string
}

func NewCNPJ(raw string) (CNPJ, error) {
	digits := stripNonDigits(raw)
	if len(digits) != 14 {
		return CNPJ{}, ErrInvalidCNPJ
	}
	if isAllSameDigits(digits) {
		return CNPJ{}, ErrInvalidCNPJ
	}
	if !validateCNPJDigits(digits) {
		return CNPJ{}, ErrInvalidCNPJ
	}
	return CNPJ{value: digits}, nil
}

func NewCNPJFromDB(digits string) CNPJ {
	return CNPJ{value: digits}
}

func (c CNPJ) String() string {
	return c.value
}

func (c CNPJ) Formatted() string {
	if len(c.value) != 14 {
		return c.value
	}
	return c.value[0:2] + "." + c.value[2:5] + "." + c.value[5:8] + "/" + c.value[8:12] + "-" + c.value[12:14]
}

func stripNonDigits(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func isAllSameDigits(s string) bool {
	for i := 1; i < len(s); i++ {
		if s[i] != s[0] {
			return false
		}
	}
	return true
}

func validateCNPJDigits(digits string) bool {
	weights1 := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	sum := 0
	for i := 0; i < 12; i++ {
		sum += int(digits[i]-'0') * weights1[i]
	}
	remainder := sum % 11
	digit1 := 0
	if remainder >= 2 {
		digit1 = 11 - remainder
	}
	if int(digits[12]-'0') != digit1 {
		return false
	}

	weights2 := []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	sum = 0
	for i := 0; i < 13; i++ {
		sum += int(digits[i]-'0') * weights2[i]
	}
	remainder = sum % 11
	digit2 := 0
	if remainder >= 2 {
		digit2 = 11 - remainder
	}
	return int(digits[13]-'0') == digit2
}
