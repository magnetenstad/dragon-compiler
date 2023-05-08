package main

import (
	"io"
	"strings"
	"unicode"
)

type Position struct {
	line int
}

type Token struct {
	Type     TokenType
	value    int
	lexeme   string
	position Position
}

type Lexer struct {
	line    int
	peek    rune
	lexemes map[string]TokenType
	reader  io.RuneReader
}

func newLexer(reader io.RuneReader) Lexer {
	lexer := Lexer{
		line:    1,
		peek:    ' ',
		lexemes: make(map[string]TokenType),
		reader:  reader,
	}
	lexer.reserve(Token{Type: tTypePrint, lexeme: "print"})
	return lexer
}

func (lexer *Lexer) reserve(token Token) {
	lexer.lexemes[token.lexeme] = token.Type
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
		if lexer.peek == '\n' {
			lexer.line += 1
			continue
		}
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

	token := Token{
		Type:   TokenType(lexer.peek),
		lexeme: string(lexer.peek),
		position: Position{
			line: lexer.line,
		},
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

		token.Type = tTypeLiteral
		token.lexeme = lexeme
		return &token, nil
	}

	if unicode.IsDigit(lexer.peek) {
		value := 0

		for unicode.IsDigit(lexer.peek) {
			value = value*10 + int(lexer.peek)
			lexer.peekNext()
		}

		token.Type = tTypeNumber
		token.value = value
		return &token, nil
	}

	if unicode.IsLetter(lexer.peek) {
		var sb strings.Builder

		for unicode.IsLetter(lexer.peek) {
			sb.WriteRune(lexer.peek)
			lexer.peekNext()
		}

		lexeme := sb.String()
		sType, exists := lexer.lexemes[lexeme]
		if !exists {
			sType = tTypeIdentifier
		}

		token.Type = sType
		token.lexeme = lexeme
		lexer.reserve(token)

		return &token, nil
	}

	if isOperator(lexer.peek) {
		// Copy paste of above
		var sb strings.Builder

		for isOperator(lexer.peek) {
			sb.WriteRune(lexer.peek)
			lexer.peekNext()
		}

		lexeme := sb.String()
		sType, exists := lexer.lexemes[lexeme]
		if !exists {
			sType = tTypeOperator
		}

		token.Type = sType
		token.lexeme = lexeme
		lexer.reserve(token)

		return &token, nil
	}

	lexer.peek = ' '
	return &token, nil
}

func (lexer *Lexer) scanAll() []Token {
	var tokens []Token
	for {
		token, err := lexer.scan()
		if err != nil {
			break
		}
		tokens = append(tokens, *token)
	}
	return tokens
}

func isOperator(r rune) bool {
	return r == '<' ||
		r == '>' ||
		r == '*' ||
		r == '/' ||
		r == '+' ||
		r == '-' ||
		r == '!' ||
		r == '='
}
