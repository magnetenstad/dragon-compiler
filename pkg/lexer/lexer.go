package lexer

import (
	"io"
	"strings"
	"unicode"
)

type TokeNType int

const (
	TTypeZero       TokeNType = 0
	TTypeIdentifier           = iota + 256
	TTypeLiteral
	TTypeNumber
	TTypeBoolean
	TTypeOperator
	TTypePrint
	TTypeRequire
	TTypeIs
	TTypeNot
)

type Position struct {
	Line int
}

type Token struct {
	Type     TokeNType
	Value    int
	Lexeme   string
	Position Position
}

type Lexer struct {
	Line    int
	Peek    rune
	Lexemes map[string]Token
	Reader  io.RuneReader
}

func NewLexer(reader io.RuneReader) Lexer {
	lexer := Lexer{
		Line:    1,
		Peek:    ' ',
		Lexemes: make(map[string]Token),
		Reader:  reader,
	}
	lexer.reserve(Token{Type: TTypeBoolean, Value: 1, Lexeme: "true"})
	lexer.reserve(Token{Type: TTypeBoolean, Value: 0, Lexeme: "false"})
	lexer.reserve(Token{Type: TTypePrint, Lexeme: "print"})
	lexer.reserve(Token{Type: TTypeRequire, Lexeme: "require"})
	lexer.reserve(Token{Type: TTypeIs, Lexeme: "is"})
	lexer.reserve(Token{Type: TTypeNot, Lexeme: "not"})
	return lexer
}

func (lexer *Lexer) reserve(token Token) {
	lexer.Lexemes[token.Lexeme] = token
}

func (lexer *Lexer) peekNext() error {
	peek, _, err := lexer.Reader.ReadRune()
	if err != nil {
		lexer.Peek = ' '
		return err
	}
	lexer.Peek = peek
	return nil
}

func (lexer *Lexer) scan() (*Token, error) {
	var err error = nil
	for ; err == nil; err = lexer.peekNext() {
		if lexer.Peek == '\n' {
			lexer.Line += 1
			continue
		}
		if lexer.Peek == ' ' ||
			lexer.Peek == '\t' ||
			lexer.Peek == '\r' {
			continue
		}
		break
	}

	if err != nil {
		return nil, err
	}

	token := Token{
		Type:   TokeNType(lexer.Peek),
		Lexeme: string(lexer.Peek),
		Position: Position{
			Line: lexer.Line,
		},
	}

	if lexer.Peek == '"' {
		var sb strings.Builder

		lexer.peekNext()
		for lexer.Peek != '"' {
			sb.WriteRune(lexer.Peek)
			if lexer.peekNext() != nil {
				panic("unclosed string")
			}
		}
		lexer.peekNext()

		lexeme := sb.String()

		token.Type = TTypeLiteral
		token.Lexeme = lexeme
		return &token, nil
	}

	if unicode.IsDigit(lexer.Peek) {
		value := 0

		for unicode.IsDigit(lexer.Peek) {
			value = value*10 + (int(lexer.Peek) - '0')
			lexer.peekNext()
		}

		token.Type = TTypeNumber
		token.Value = value
		return &token, nil
	}

	if unicode.IsLetter(lexer.Peek) {
		var sb strings.Builder

		for unicode.IsLetter(lexer.Peek) {
			sb.WriteRune(lexer.Peek)
			lexer.peekNext()
		}

		lexeme := sb.String()
		token.Type = TTypeIdentifier
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

	if isOperator(lexer.Peek) {
		// Copy paste of above
		var sb strings.Builder

		for isOperator(lexer.Peek) {
			sb.WriteRune(lexer.Peek)
			lexer.peekNext()
		}

		lexeme := sb.String()
		token.Type = TTypeOperator
		token.Lexeme = lexeme
		lexer.reserve(token)

		return &token, nil
	}

	lexer.Peek = ' '
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
