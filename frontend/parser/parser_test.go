package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func getLeftmostTerm(expression *Expression) *Term {
	return expression.Or[0].And[0].LHS.Summand.LHS.LHS
}

func TestParseSelectAllQueryBasicPasses(t *testing.T) {
	parser := BuildParser(&SelectQuery{})
	result := &SelectQuery{}
	err := parser.ParseString("SELECT * FROM sometable", result)
	if assert.NoError(t, err) {
		select_expr := result.SelectClause.Expression
		assert.True(t, select_expr.All)
		assert.Nil(t, select_expr.Expressions)

		assert.Equal(t, "sometable", result.FromClause.FromExpression.Target)
		assert.Nil(t, result.FromClause.FromExpression.WhereClause)
	}
}

func TestParseSelectFieldsQueryBasicPasses(t *testing.T) {
	parser := BuildParser(&SelectQuery{})
	result := &SelectQuery{}
	err := parser.ParseString("SELECT abc FROM target1", result)
	if assert.NoError(t, err) {
		select_expr := result.SelectClause.Expression
		assert.False(t, select_expr.All)

		assert.Len(t, select_expr.Expressions, 1)
		assert.Equal(t, "abc", getLeftmostTerm(select_expr.Expressions[0].Expression).SymbolRef.Symbol)

		assert.Equal(t, "target1", result.FromClause.FromExpression.Target)
		assert.Nil(t, result.FromClause.FromExpression.WhereClause)
	}

	result = &SelectQuery{}
	err = parser.ParseString("SELECT abc, def1, g__hi FROM abcd_2", result)
	if assert.NoError(t, err) {
		select_expr := result.SelectClause.Expression
		assert.False(t, select_expr.All)

		assert.Len(t, select_expr.Expressions, 3)
		assert.Equal(t, "abc", getLeftmostTerm(select_expr.Expressions[0].Expression).SymbolRef.Symbol)
		assert.Equal(t, "def1", getLeftmostTerm(select_expr.Expressions[1].Expression).SymbolRef.Symbol)
		assert.Equal(t, "g__hi", getLeftmostTerm(select_expr.Expressions[2].Expression).SymbolRef.Symbol)

		assert.Equal(t, "abcd_2", result.FromClause.FromExpression.Target)
		assert.Nil(t, result.FromClause.FromExpression.WhereClause)
	}
}

func TestParseIntConstantValue(t *testing.T) {
	parser := BuildParser(&ConstantValue{})
	result := &ConstantValue{}
	err := parser.ParseString("1", result)
	if assert.NoError(t, err) {
		assert.Equal(t, float64(1), *result.Number)
		assert.Nil(t, result.String)
		assert.Nil(t, result.Boolean)
	}

	result = &ConstantValue{}
	err = parser.ParseString("0.56e5", result)
	if assert.NoError(t, err) {
		assert.Equal(t, 0.56e5, *result.Number)
		assert.Nil(t, result.String)
		assert.Nil(t, result.Boolean)
	}
}

func TestParseStringConstantValue(t *testing.T) {
	parser := BuildParser(&ConstantValue{})
	result := &ConstantValue{}
	err := parser.ParseString("\"hello\"", result)
	if assert.NoError(t, err) {
		assert.Equal(t, "hello", *result.String)
		assert.Nil(t, result.Number)
		assert.Nil(t, result.Boolean)
	}

	result = &ConstantValue{}
	err = parser.ParseString("\"\"", result)
	if assert.NoError(t, err) {
		assert.Equal(t, "", *result.String)
		assert.Nil(t, result.Number)
		assert.Nil(t, result.Boolean)
	}

	result = &ConstantValue{}
	err = parser.ParseString("\"a\"", result)
	if assert.NoError(t, err) {
		assert.Equal(t, "a", *result.String)
		assert.Nil(t, result.Number)
		assert.Nil(t, result.Boolean)
	}

	result = &ConstantValue{}
	err = parser.ParseString("\"sgasg\nd\"", result)
	if assert.NoError(t, err) {
		assert.Equal(t, "sgasg\nd", *result.String)
		assert.Nil(t, result.Number)
		assert.Nil(t, result.Boolean)
	}
}

func TestParseBoolConstantValue(t *testing.T) {
	parser := BuildParser(&ConstantValue{})
	result := &ConstantValue{}
	err := parser.ParseString("true", result)
	if assert.NoError(t, err) {
		assert.True(t, bool(*result.Boolean))
		assert.Nil(t, result.Number)
		assert.Nil(t, result.String)
	}

	result = &ConstantValue{}
	err = parser.ParseString("TRue", result)
	if assert.NoError(t, err) {
		assert.True(t, bool(*result.Boolean))
		assert.Nil(t, result.Number)
		assert.Nil(t, result.String)
	}

	result = &ConstantValue{}
	err = parser.ParseString("TRUE", result)
	if assert.NoError(t, err) {
		assert.True(t, bool(*result.Boolean))
		assert.Nil(t, result.Number)
		assert.Nil(t, result.String)
	}

	result = &ConstantValue{}
	err = parser.ParseString("FALSE", result)
	if assert.NoError(t, err) {
		assert.False(t, bool(*result.Boolean))
		assert.Nil(t, result.Number)
		assert.Nil(t, result.String)
	}

	result = &ConstantValue{}
	err = parser.ParseString("faLse", result)
	if assert.NoError(t, err) {
		assert.False(t, bool(*result.Boolean))
		assert.Nil(t, result.Number)
		assert.Nil(t, result.String)
	}
}

func TestParseConstantValueFails(t *testing.T) {
	parser := BuildParser(&ConstantValue{})
	result := &ConstantValue{}
	err := parser.ParseString("\"", result)
	assert.Error(t, err)
	assert.Nil(t, result.String)
	assert.Nil(t, result.Number)
	assert.Nil(t, result.Boolean)

	result = &ConstantValue{}
	err = parser.ParseString("abjvd", result)
	assert.Error(t, err)
	assert.Nil(t, result.String)
	assert.Nil(t, result.Number)
	assert.Nil(t, result.Boolean)
}
