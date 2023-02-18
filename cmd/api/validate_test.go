package main

import (
	"fmt"
	"regexp"
	"testing"
	"unicode"
)

var testRes bool

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
	phrase := "cOm(P)l:exP@@5"
	for n := 0; n < b.N; n++ {
		testRes = validatePassASCII(phrase)
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
