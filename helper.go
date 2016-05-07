package main

import (
	"unicode"
	"unicode/utf8"
)

func lowerFirst(str string) string {
	if str == "" {
		return ""
	}
	r, size := utf8.DecodeRuneInString(str)
	return string(unicode.ToLower(r)) + str[size:]
}

func upperFirst(str string) string {
	if str == "" {
		return ""
	}
	r, size := utf8.DecodeRuneInString(str)
	return string(unicode.ToUpper(r)) + str[size:]
}
