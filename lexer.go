package main

import (
	"io"
	"strings"
	"unicode"
)

type Token struct {
	tag    Tag
	value  int
	lexeme string
}

type Tag int

const (
	TagZero Tag = 0
	TagNum  Tag = iota + 256
	TagId
	TagTrue
	TagFalse
	TagString
	TagExpr
	TagIf
	TagFor
)

type Lexer struct {
	line   int
	peek   rune
	words  map[string]Token
	reader io.RuneReader
}

func newLexer(reader io.RuneReader) Lexer {
	lexer := Lexer{
		line:   0,
		peek:   ' ',
		words:  make(map[string]Token),
		reader: reader,
	}
	lexer.reserve(Token{TagTrue, 1, "true"})
	lexer.reserve(Token{TagFalse, 0, "false"})
	lexer.reserve(Token{TagExpr, 0, "expr"}) // TODO: Remove this
	lexer.reserve(Token{TagIf, 0, "if"})
	lexer.reserve(Token{TagFor, 0, "for"})
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
			lexer.peek == '\n' ||
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
			tag:    TagString,
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
			tag:   TagNum,
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
		word, ok := lexer.words[lexeme]

		if ok {
			return &word, nil
		}
		token := Token{
			tag:    TagId,
			lexeme: lexeme,
		}
		lexer.reserve(token)

		return &token, nil
	}

	token := Token{
		tag:    Tag(lexer.peek),
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
