package parser

import (
	"strings"

	"github.com/alecthomas/participle"
	"github.com/kaschnit/csvquery/frontend/lexer"
)

type Boolean bool

func (b *Boolean) Capture(values []string) error {
	*b = strings.ToUpper(values[0]) == "TRUE"
	return nil
}

type SelectQuery struct {
	SelectClause *SelectClause `@@`
	FromClause   *FromClause   `@@`
}

type SelectClause struct {
	Expression *SelectExpression `"SELECT" @@`
}

type SelectExpression struct {
	All         bool                 `@"*"`
	Expressions []*AliasedExpression `| @@ ( "," @@ )*`
}

type AliasedExpression struct {
	Expression *Expression `@@`
	As         string      `( "AS" @Ident )?`
}

type FromClause struct {
	FromExpression *FromExpression `"FROM" @@`
}

type FromExpression struct {
	Target      string       `@Ident`
	WhereClause *WhereClause `( @@ )?`
}

type WhereClause struct {
	Expression *Expression `"WHERE" @@`
}

type Expression struct {
	Or []*OrCondition `@@ ( "OR" @@ )*`
}

type OrCondition struct {
	And []*CompareConditon `@@ ( "AND" @@ )*`
}

type CompareConditon struct {
	LHS *Operand `@@`
	Op  string   `( @( "<>" | "<=" | ">=" | "=" | "<" | ">" | "!=" )`
	RHS *Operand `@@ )?`
}

type Operand struct {
	Summand *Summand `@@`
}

type Summand struct {
	LHS *Factor `@@`
	Op  string  `( @("+" | "-")`
	RHS *Factor `  @@ )?`
}

type Factor struct {
	LHS *Term  `@@`
	Op  string `( @("*" | "/" | "%")`
	RHS *Term  `@@ )?`
}

type Term struct {
	ConstantValue *ConstantValue `@@`
	SymbolRef     *SymbolRef     `| @@`
	SubExpression *Expression    `| "(" @@ ")"`
}

type SymbolRef struct {
	Symbol     string        `@Ident`
	Parameters []*Expression `( "(" @@ { "," @@ } ")" )?`
}

type ConstantValue struct {
	Number  *float64 `@Number`
	String  *string  `| @String`
	Boolean *Boolean `| @("TRUE" | "FALSE")`
}

func BuildParser(grammar interface{}) *participle.Parser {
	return participle.MustBuild(
		grammar,
		participle.Lexer(lexer.QueryLexer),
		participle.Unquote("String"),
		participle.CaseInsensitive("Keyword"),
	)
}
