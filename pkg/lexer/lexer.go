package lexer

import (
	"io"
	"strings"
	"unicode"
)

type TokenType int

const (
	TypeZero       TokenType = 0
	TypeIdentifier           = iota + 256
	TypeLiteral
	TypeNumber
	TypeBoolean
	TypeOperator
	TypePrint
	TypeRequire
	TypeIs
	TypeNot
)

type Position struct {
	Line int
}

type Token struct {
	Type     TokenType
	Value    int
	Lexeme   string
	Position Position
}

type Lexer struct {
	line    int
	peek    rune
	Lexemes map[string]Token
	reader  io.RuneReader
}

func NewLexer(reader io.RuneReader) Lexer {
	lexer := Lexer{
		line:    1,
		peek:    ' ',
		Lexemes: make(map[string]Token),
		reader:  reader,
	}
	lexer.reserve(Token{Type: TypeBoolean, Value: 1, Lexeme: "true"})
	lexer.reserve(Token{Type: TypeBoolean, Value: 0, Lexeme: "false"})
	lexer.reserve(Token{Type: TypePrint, Lexeme: "print"})
	lexer.reserve(Token{Type: TypeRequire, Lexeme: "require"})
	lexer.reserve(Token{Type: TypeIs, Lexeme: "is"})
	lexer.reserve(Token{Type: TypeNot, Lexeme: "not"})
	return lexer
}

func (lexer *Lexer) reserve(token Token) {
	lexer.Lexemes[token.Lexeme] = token
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
		Lexeme: string(lexer.peek),
		Position: Position{
			Line: lexer.line,
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

		token.Type = TypeLiteral
		token.Lexeme = lexeme
		return &token, nil
	}

	if unicode.IsDigit(lexer.peek) {
		value := 0

		for unicode.IsDigit(lexer.peek) {
			value = value*10 + (int(lexer.peek) - '0')
			lexer.peekNext()
		}

		token.Type = TypeNumber
		token.Value = value
		return &token, nil
	}

	if unicode.IsLetter(lexer.peek) {
		var sb strings.Builder

		for unicode.IsLetter(lexer.peek) {
			sb.WriteRune(lexer.peek)
			lexer.peekNext()
		}

		lexeme := sb.String()
		token.Type = TypeIdentifier
		token.Lexeme = lexeme

		existingToken, exists := lexer.Lexemes[lexeme]
		if exists {
			token.Type = existingToken.Type
			token.Lexeme = existingToken.Lexeme
			token.Value = existingToken.Value
		}

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
		token.Type = TypeOperator
		token.Lexeme = lexeme
		lexer.reserve(token)

		return &token, nil
	}

	lexer.peek = ' '
	return &token, nil
}

func (lexer *Lexer) ScanAll() []Token {
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
