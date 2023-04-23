package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

type TokenType string
const (
    Unknown TokenType = ""
    Quote = "\""
    NewLine = "\\n"
    Space = " "
    Dash = "-"
    Word = "\\S+"
)

var tokenTypes = []TokenType {
    Unknown,
    Quote,
    NewLine,
    Space,
    Dash,
    Word,
}

type Token struct {
    data string
    tokenType TokenType
}

func main() {

    file, err := os.ReadFile("example.txt")
    check(err)

    tokens := lex(file)

    fmt.Println(tokens)
}

func lex(bytes []byte) []Token {
    
    var tokens []Token
    var sb strings.Builder
    var token Token

    for _, b := range bytes {
        sb.WriteByte(b)

        tokenType := Unknown
        for _, t := range tokenTypes {
            matched, err := regexp.MatchString(
                fmt.Sprintf("^%s$", t), sb.String())
            check(err)

            if matched {  
                tokenType = t
                break
            }
        }
        
        if tokenType == Unknown {
            tokens = append(tokens, token)
            sb.Reset()
            continue
        }

        token.tokenType = tokenType
        token.data = sb.String()

        if (tokenType != Word) {
            tokens = append(tokens, token)
            sb.Reset()
        }
    }
    tokens = append(tokens, token)
    sb.Reset()
    
    return tokens
}
