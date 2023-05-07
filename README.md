# dragon-compiler

## Lexical Analysis

## Tokens

```txt

literal     -> ".*"
number      -> [0-9]+
identifier  -> [a-zA-Z]+
operator    -> = | < | > | * | / | + | - | <= | >= | == | !=

octothorpe  -> #
print       -> print
```

## Syntax Analysis

### BNF Grammar

```txt

program     -> blocks
blocks      -> blocks block | block
block       -> { statements blocks }
statements  -> statements statement | statement
statement   ->
      identifier = expression ;
    | print expression ;
    | # expression ;
expression  -> ( expression operator expression ) | number | literal | identifier

```
