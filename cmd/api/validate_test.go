package main

import (
	"fmt"
	"regexp"
	"testing"
	"unicode"

	"github.com/michalzoldak97/go-auth/internal/data"
)

var testRes bool

func Test_validateASCIIPass(t *testing.T) {

	var testApp application

	testApp.security = data.SecurityConfig{
		PassLower:    true,
		PassUpper:    true,
		PassNum:      true,
		PassSpecial:  true,
		PassMinLen:   8,
		PassMaxLen:   30,
		MaxPOSTBytes: 128,
	}

	phrases := map[string]bool{
		"cOm(P)l:exP@@5":        true,
		"           ":           false,
		"&^#*/*?<>{||` <>.?/:}": false,
		"~-=*^^eW%6)r_U+-U":     true,
		"r$9":                   false,
		"tewat488GRH	vsag{}@km7*/vknb   > feong wqe": false}

	for phrase, res := range phrases {
		if testApp.validatePassASCII(phrase) != res {
			t.Errorf("Validation failed for phrase %v", phrase)
		}
	}
}

func validatePassRegex(phrase string) bool {
	rules := []string{"[a-z]", "[A-Z]", "[0-9]", "[^\\d\\w]"}

	for _, rule := range rules {
		pass, _ := regexp.MatchString(rule, phrase)
		if !pass {
			return false
		}
	}

	return true
}

func BenchmarkRegexPass(b *testing.B) {
	phrase := "cOm(P)l:exP@@5"
	for n := 0; n < b.N; n++ {
		testRes = validatePassRegex(phrase)
	}
}

func BenchmarkASCIIPass(b *testing.B) {
	var testApp application

	testApp.security = data.SecurityConfig{
		PassLower:    true,
		PassUpper:    true,
		PassNum:      true,
		PassSpecial:  true,
		PassMinLen:   8,
		PassMaxLen:   30,
		MaxPOSTBytes: 128,
	}

	phrase := "cOm(P)l:exP@@5"
	for n := 0; n < b.N; n++ {
		testRes = testApp.validatePassASCII(phrase)
	}
	if !testRes {
		fmt.Println("Not")
	}
}

func validatePassUnicodeSwitch(phrase string) bool {
	isLower, isUpper, isNum, isSpec := false, false, false, false

	for _, char := range phrase {
		if isLower && isUpper && isNum && isSpec {
			break
		}

		switch {
		case unicode.IsLower(char):
			isLower = true
		case unicode.IsUpper(char):
			isUpper = true
		case unicode.IsNumber(char):
			isNum = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			isSpec = true
		}
	}

	return isLower && isUpper && isNum && isSpec
}

func BenchmarkUnicodeSwitchPass(b *testing.B) {
	phrase := "cOm(P)l:exP@@5"
	for n := 0; n < b.N; n++ {
		testRes = validatePassUnicodeSwitch(phrase)
	}
}

func validatePassUnicodeIf(phrase string) bool {
	isLower, isUpper, isNum, isSpec := false, false, false, false

	for _, char := range phrase {
		if !isLower && unicode.IsLower(char) {
			isLower = true
			continue
		}

		if !isUpper && unicode.IsUpper(char) {
			isUpper = true
			continue
		}

		if !isNum && unicode.IsNumber(char) {
			isNum = true
			continue
		}

		if !isSpec && (unicode.IsPunct(char) || unicode.IsSymbol(char)) {
			isSpec = true
			continue
		}
	}

	return isLower && isUpper && isNum && isSpec
}

func BenchmarkUnicodeIfPass(b *testing.B) {
	phrase := "cOm(P)l:exP@@5"
	for n := 0; n < b.N; n++ {
		testRes = validatePassUnicodeIf(phrase)
	}
}
