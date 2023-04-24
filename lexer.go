package main

import (
	"io"
	"strings"
	"unicode"
)

type Token struct {
    tag int
    value int
    lexeme string
}

type Tag int
const (
    Zero Tag = 0
    NumTag Tag = iota + 256
    IdTag
    TrueTag
    FalseTag
)

type Lexer struct {
    line int
    peek rune
    words map[string]Token
    reader io.RuneReader
}

func (lexer *Lexer) scan() (*Token, error) {
    var err error = nil
    stop := false
    for !stop {
        stop, err = lexer.peekNext()
        if lexer.peek == ' ' || 
                lexer.peek == '\t' || 
                lexer.peek == '\r' { 
            continue 
        }
        if lexer.peek == '\n' { 
            lexer.line += 1
            continue
        }
        break
    }

    if (err != nil) {
        return nil, err
    }

    if unicode.IsDigit(lexer.peek) {
        value := 0
        
        for unicode.IsDigit(lexer.peek) {
            value = value * 10 + int(lexer.peek)
            lexer.peekNext()
        }
        
        return &Token { int(NumTag), value, "" }, nil
    }

    if unicode.IsLetter(lexer.peek) {
        var sb strings.Builder

        for unicode.IsLetter(lexer.peek) {
            sb.WriteRune(lexer.peek)
            lexer.peekNext()
        }
        
        lexeme := sb.String()
        word := lexer.words[lexeme]

        if (word.tag != 0) { return &word, nil }
        t := Token { int(IdTag), 0, lexeme }
        lexer.words[lexeme] = t
        
        return &t, nil
    }

    t := Token { int(lexer.peek), 0, string(lexer.peek) }
    lexer.peek = ' '
    return &t, nil
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

func (lexer *Lexer) peekNext() (bool, error) {
    peek, _, err := lexer.reader.ReadRune()
    if err != nil {
        lexer.peek = ' '
        return true, err
    }
    lexer.peek = peek
    return false, nil
}

func (lexer *Lexer) reserve(token Token) {
    lexer.words[token.lexeme] = token
}

func newLexer(reader io.RuneReader) Lexer {
    lexer := Lexer { 0, ' ', make(map[string]Token), reader }
    lexer.reserve(Token { int(TrueTag), 1, "true" })
    lexer.reserve(Token { int(FalseTag), 0, "false" })
    return lexer
}
