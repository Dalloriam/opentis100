package compiler

import (
	"unicode"
)

const (
	LINEBREAK = "LINEBREAK"
	AT        = "AT"
	COLON     = "COLON"
	WORD      = "WORD"
	DIGIT     = "DIGIT"
	EOF       = "EOF"
	LABEL     = "LABEL"
)

type Token struct {
	TokType string
	Value   string
}

func Tokenize(source string) (<-chan *Token, <-chan error) {
	outChan := make(chan *Token)
	errChan := make(chan error)

	go func() {
		defer close(outChan)
		defer close(errChan)

		for i := 0; i < len(source); i++ {
			head := source[i]

			if head == '\n' {
				outChan <- &Token{LINEBREAK, "\n"}
			} else if head == '@' {
				outChan <- &Token{AT, "@"}
			} else if unicode.IsLetter(rune(head)) {
				var tokValue []uint8

				for i < len(source) && (unicode.IsDigit(rune(source[i])) || unicode.IsLetter(rune(source[i]))) {
					tokValue = append(tokValue, source[i])
					i++
				}

				var tokType string
				if i < len(source) && source[i] == ':' {
					// Current word is a label
					tokType = LABEL
				} else {
					// TODO: Distinct. keywords & arguments
					tokType = WORD
					i--
				}
				outChan <- &Token{tokType, string(tokValue)}
			} else if unicode.IsDigit(rune(head)) {
				var tokValue []uint8

				for i < len(source) && unicode.IsDigit(rune(source[i])) {
					tokValue = append(tokValue, source[i])
					i++
				}
				i--

				outChan <- &Token{DIGIT, string(tokValue)}
			}
		}

		outChan <- &Token{EOF, "EOF"}

	}()

	return outChan, errChan
}
