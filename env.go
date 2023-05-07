package main

type Symbol struct {
	lexeme     string
	symbolType SymbolType
}

type Env struct {
	table map[string]Symbol
	prev  *Env
}

func newEnv(prev *Env) Env {
	env := Env{
		table: make(map[string]Symbol),
		prev:  prev,
	}
	return env
}

func (env Env) put(symbol Symbol) {
	env.table[symbol.lexeme] = symbol
}

func (env Env) get(key string) (Symbol, bool) {
	for e := &env; e != nil; e = e.prev {
		symbol, ok := e.table[key]
		if ok {
			return symbol, true
		}
	}
	return Symbol{}, false
}
