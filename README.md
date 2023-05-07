# dragon-compiler

## Lexical Analysis

## Tokens

```txt

literal     -> ".*"
number      -> [0-9]+
identifier  -> [a-zA-Z]+
operator    -> = | < | > | * | / | + | - | <= | >= | == | !=

print       -> print
```

## Syntax Analysis

### BNF Grammar

```txt

program     -> blocks

blocks      ->
    block blocks
  | block

block       ->
  { blocks }
  { statements }

statements  ->
    statement statements
  | statement

statement   ->
    identifier = expression ;
  | print expression ;
  | # expression ;

expression  ->
    ( expression operator expression )
  | number
  | literal
  | identifier

```
