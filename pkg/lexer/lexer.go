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
	TypeNot
	TypeStruct
	TypeTypeHint
	TypeSkip
	TypeSkipIf
)

func (e TokenType) String() string {
	switch e {
	case TypeIdentifier:
		return "TypeIdentifier"
	case TypeLiteral:
		return "TypeLiteral"
	case TypeNumber:
		return "TypeNumber"
	case TypeBoolean:
		return "TypeBoolean"
	case TypeOperator:
		return "TypeOperator"
	case TypePrint:
		return "TypePrint"
	case TypeNot:
		return "TypeNot"
	case TypeStruct:
		return "TypeStruct"
	case TypeTypeHint:
		return "TypeTypeHint"
	default:
		return string(rune(e))
	}
}

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
	lexer.reserve(Token{Type: TypeNot, Lexeme: "!"})
	lexer.reserve(Token{Type: TypeStruct, Lexeme: "struct"})
	lexer.reserve(Token{Type: TypeSkip, Lexeme: "skip"})
	lexer.reserve(Token{Type: TypeSkipIf, Lexeme: "skip_if"})
	lexer.reserve(Token{Type: TypeTypeHint, Lexeme: "Int"})
	lexer.reserve(Token{Type: TypeTypeHint, Lexeme: "Float"})
	lexer.reserve(Token{Type: TypeTypeHint, Lexeme: "String"})
	lexer.reserve(Token{Type: TypeTypeHint, Lexeme: "Bool"})
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

	err := lexer.scanWhiteSpace()

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
		return lexer.scanStringLiteral(token)
	}

	if unicode.IsDigit(lexer.peek) {
		return lexer.scanNumber(token)
	}

	if unicode.IsLetter(lexer.peek) {
		return lexer.scanWord(token)
	}

	if isOperator(lexer.peek) {
		return lexer.scanOperator(token)
	}

	lexer.peek = ' '
	return &token, nil
}

func (lexer *Lexer) scanWhiteSpace() error {
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
	return err
}

func (lexer *Lexer) scanStringLiteral(token Token) (*Token, error) {
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

func (lexer *Lexer) scanNumber(token Token) (*Token, error) {
	value := 0

	for unicode.IsDigit(lexer.peek) {
		value = value*10 + (int(lexer.peek) - '0')
		lexer.peekNext()
	}

	token.Type = TypeNumber
	token.Value = value
	return &token, nil
}

func (lexer *Lexer) scanWord(token Token) (*Token, error) {
	var sb strings.Builder

	tokenType := TypeIdentifier
	if unicode.IsUpper(lexer.peek) { // TODO
		tokenType = TypeTypeHint
	}

	for unicode.IsLetter(lexer.peek) || lexer.peek == '_' {
		sb.WriteRune(lexer.peek)
		lexer.peekNext()
	}

	lexeme := sb.String()
	token.Type = TokenType(tokenType)
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

func (lexer *Lexer) scanOperator(token Token) (*Token, error) {
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
		r == '-'
}
