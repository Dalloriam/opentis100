package compiler

import (
	"unicode"
)

const (
	linebreak = "LINEBREAK"
	at        = "AT"
	colon     = "COLON"
	word      = "WORD"
	digit     = "DIGIT"
)

type Token struct {
	tokType string
	value   string
}

func Tokenize(source string) (<-chan Token, <-chan error) {
	outChan := make(chan Token)
	errChan := make(chan error)

	go func() {
		defer close(outChan)
		defer close(errChan)

		for i := 0; i < len(source); i++ {
			head := source[i]

			if head == '\n' {
				outChan <- Token{linebreak, "\n"}
			} else if head == '@' {
				outChan <- Token{at, "@"}
			} else if head == ':' {
				outChan <- Token{colon, ":"}
			} else if unicode.IsLetter(rune(head)) {
				var tokValue []uint8

				for i < len(source) && (unicode.IsDigit(rune(source[i])) || unicode.IsLetter(rune(source[i]))) {
					tokValue = append(tokValue, source[i])
					i++
				}
				i--
				outChan <- Token{word, string(tokValue)}
			} else if unicode.IsDigit(rune(head)) {
				var tokValue []uint8

				for i < len(source) && unicode.IsDigit(rune(source[i])) {
					tokValue = append(tokValue, source[i])
					i++
				}
				i--

				outChan <- Token{digit, string(tokValue)}
			}
		}

	}()

	return outChan, errChan
}
