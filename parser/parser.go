package parser

import (
	"github.com/alecthomas/participle"
)

type SelectQuery struct {
	SelectClause *SelectClause `SELECT @@`
}

type SelectClause struct {
	Expression *SelectExpression
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
	FromExpression *FromExpression `FROM @@`
}

type FromExpression struct {
}

type Expression struct {
	Or []*OrCondition `@@ ( "OR" @@ )*`
}

type OrCondition struct {
	And *[]CompareConditon `@@ ( "AND" @@ )*`
}

type CompareConditon struct {
	LHS          *Operand `@@`
	ConditionRHS string   `( @( "<>" | "<=" | ">=" | "=" | "<" | ">" | "!=" )`
	RHS          *Operand `@@ )?`
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
}

type SymbolRef struct {
	Symbol     string        `@Ident`
	Parameters []*Expression `( "(" @@ { "," @@ } ")" )?`
}

type ConstantValue struct {
	Number  *float64 `@Number`
	String  *string  `| @String`
	Boolean *bool    `| @("TRUE" | "FALSE")`
}

var queryParser = participle.MustBuild(
	&SelectQuery{},
	participle.Lexer(),
	participle.Unquote("String"),
	participle.CaseInsensitive("Keyword"),
)
