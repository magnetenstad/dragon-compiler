package env

import (
	"github.com/magnetenstad/dragon-compiler/pkg/ast"
)

type Symbol struct {
	Lexeme     string
	SymbolType ast.NodeType
}

type Env struct {
	Table map[string]Symbol
	Prev  *Env
}

func NewEnv(prev *Env) Env {
	env := Env{
		Table: make(map[string]Symbol),
		Prev:  prev,
	}
	return env
}

func (env Env) Put(symbol Symbol) {
	env.Table[symbol.Lexeme] = symbol
}

func (env Env) Get(key string) (Symbol, bool) {
	for e := &env; e != nil; e = e.Prev {
		symbol, ok := e.Table[key]
		if ok {
			return symbol, true
		}
	}
	return Symbol{}, false
}
