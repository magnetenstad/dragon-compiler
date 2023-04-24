package main

import (
	"io"
	"strings"
	"unicode"
)

type Token struct {
	tag    int
	value  int
	lexeme string
}

type Tag int

const (
	Zero   Tag = 0
	NumTag Tag = iota + 256
	IdTag
	TrueTag
	FalseTag
	StringTag
)

type Lexer struct {
	line   int
	peek   rune
	words  map[string]Token
	reader io.RuneReader
}

func newLexer(reader io.RuneReader) Lexer {
	lexer := Lexer{0, ' ', make(map[string]Token), reader}
	lexer.reserve(Token{int(TrueTag), 1, "true"})
	lexer.reserve(Token{int(FalseTag), 0, "false"})
	return lexer
}

func (lexer *Lexer) reserve(token Token) {
	lexer.words[token.lexeme] = token
}

func (lexer *Lexer) peekNext() error {
	peek, _, err := lexer.reader.ReadRune()
	if err != nil {
		lexer.peek = ' '
		return err
	}
	lexer.peek = peek
	return nil
}

func (lexer *Lexer) scan() (*Token, error) {
	var err error = nil
	for ; err == nil; err = lexer.peekNext() {
		if lexer.peek == ' ' ||
			lexer.peek == '\t' ||
			lexer.peek == '\r' {
			continue
		}
		break
	}

	if err != nil {
		return nil, err
	}

	if lexer.peek == '"' {
		var sb strings.Builder

		lexer.peekNext()
		for lexer.peek != '"' {
			sb.WriteRune(lexer.peek)
			if lexer.peekNext() != nil {
				panic("unclosed string")
			}
		}
		lexer.peekNext()

		lexeme := sb.String()
		token := Token{
			tag:    int(StringTag),
			lexeme: lexeme,
		}
		return &token, nil
	}

	if unicode.IsDigit(lexer.peek) {
		value := 0

		for unicode.IsDigit(lexer.peek) {
			value = value*10 + int(lexer.peek)
			lexer.peekNext()
		}

		token := Token{
			tag:   int(NumTag),
			value: value,
		}
		return &token, nil
	}

	if unicode.IsLetter(lexer.peek) {
		var sb strings.Builder

		for unicode.IsLetter(lexer.peek) {
			sb.WriteRune(lexer.peek)
			lexer.peekNext()
		}

		lexeme := sb.String()
		word := lexer.words[lexeme]

		if word.tag != 0 {
			return &word, nil
		}
		token := Token{
			tag:    int(IdTag),
			lexeme: string(lexeme),
		}
		lexer.words[lexeme] = token

		return &token, nil
	}

	token := Token{
		tag:    int(lexer.peek),
		lexeme: string(lexer.peek),
	}
	lexer.peek = ' '
	return &token, nil
}

func (lexer *Lexer) scanAll() []Token {
	var tokens []Token
	for true {
		token, err := lexer.scan()
		if err != nil {
			break
		}
		tokens = append(tokens, *token)
	}
	return tokens
}
