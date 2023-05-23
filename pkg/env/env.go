package env

import (
	"github.com/magnetenstad/dragon-compiler/pkg/ast"
)

type Symbol struct {
	Lexeme     string
	SymbolType ast.NodeType
}

type Env struct {
	table map[string]Symbol
	prev  *Env
}

func NewEnv(prev *Env) Env {
	env := Env{
		table: make(map[string]Symbol),
		prev:  prev,
	}
	return env
}

func (env Env) Put(symbol Symbol) {
	env.table[symbol.Lexeme] = symbol
}

func (env Env) Get(key string) (Symbol, bool) {
	for e := &env; e != nil; e = e.prev {
		symbol, ok := e.table[key]
		if ok {
			return symbol, true
		}
	}
	return Symbol{}, false
}
