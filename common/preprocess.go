package common

import (
	"strings"
	"unicode"
)

func tokenize(text string) []string {
	return strings.FieldsFunc(text, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
}

func lower(tokens []string) []string {
	res := make([]string, len(tokens))
	for i, t := range tokens {
		res[i] = strings.ToLower(t)
	}

	return res
}

var stopwords = map[string]bool{
	"a":    true,
	"and":  true,
	"be":   true,
	"have": true,
	"i":    true,
	"in":   true,
	"of":   true,
	"that": true,
	"the":  true,
	"to":   true,
}

func commonWordFilter(tokens []string) []string {
	res := make([]string, len(tokens))
	for _, t := range tokens {
		if _, ok := stopwords[t]; !ok {
			res = append(res, t)
		}
	}

	return res
}

func TokenizeAndFilter(text string) []string {
	tokens := tokenize(text)
	tokens = lower(tokens)
	tokens = commonWordFilter(tokens)
	return tokens
}
