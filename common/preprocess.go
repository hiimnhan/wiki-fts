package common

import (
	"strings"
	"unicode"
)

func tokenize(text string) []string {
	return strings.FieldsFunc(text, func(r rune) bool {
		// return true if r is not a letter or a digit
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
}

func lower(tokens []string) []string {
	res := make([]string, len(tokens))
	for i, t := range tokens {
		res[i] = strings.ToLower(t)
	}

	return res
}

var stopwords = map[string]struct{}{
	"a":    {},
	"and":  {},
	"be":   {},
	"have": {},
	"i":    {},
	"in":   {},
	"of":   {},
	"that": {},
	"the":  {},
	"to":   {},
	"":     {},
}

func commonWordFilter(tokens []string) []string {
	res := make([]string, 0, len(tokens))
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
