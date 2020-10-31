package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var selectQueryParser = BuildParser(&SelectQuery{})

func exprLeftmostTerm(expression *Expression) *Term {
	return expression.Or[0].And[0].LHS.Summand.LHS.LHS
}

func TestSelectNoSelectFieldsFails(t *testing.T) {
	result := &SelectQuery{}
	err := selectQueryParser.ParseString("SELECT FROM sometable", result)
	assert.Error(t, err)
}

func TestSelectAllAndSelectFieldsFails(t *testing.T) {
	result := &SelectQuery{}
	err := selectQueryParser.ParseString("SELECT *, a FROM sometable", result)
	assert.Error(t, err)

	result = &SelectQuery{}
	err = selectQueryParser.ParseString("SELECT a, * FROM sometable", result)
	assert.Error(t, err)

	result = &SelectQuery{}
	err = selectQueryParser.ParseString("SELECT a, *, b FROM sometable", result)
	assert.Error(t, err)
}

func TestSelectAllAliasedFails(t *testing.T) {
	result := &SelectQuery{}
	err := selectQueryParser.ParseString("SELECT * as vals FROM sometable", result)
	assert.Error(t, err)
}

func TestSelectNoFromTargetFails(t *testing.T) {
	result := &SelectQuery{}
	err := selectQueryParser.ParseString("SELECT a FROM ", result)
	assert.Error(t, err)

	result = &SelectQuery{}
	err = selectQueryParser.ParseString("SELECT a FROM WHERE a < 3", result)
	assert.Error(t, err)
}

func TestSelectAllQueryBasicPasses(t *testing.T) {
	result := &SelectQuery{}
	err := selectQueryParser.ParseString("SELECT * FROM sometable", result)
	if assert.NoError(t, err) {
		selectExpr := result.SelectClause.Expression
		assert.True(t, selectExpr.All)
		assert.Nil(t, selectExpr.Expressions)

		assert.Equal(t, "sometable", result.FromClause.FromExpression.Target)
		assert.Nil(t, result.FromClause.FromExpression.WhereClause)
	}
}

func TestSelectWhereBasicPasses(t *testing.T) {
	result := &SelectQuery{}
	err := selectQueryParser.ParseString("SELECT * FROM sometable WHERE x < 3", result)
	if assert.NoError(t, err) {
		selectExpr := result.SelectClause.Expression
		assert.True(t, selectExpr.All)
		assert.Nil(t, selectExpr.Expressions)

		assert.Equal(t, "sometable", result.FromClause.FromExpression.Target)

		whereClause := result.FromClause.FromExpression.WhereClause
		assert.Equal(t, "x", exprLeftmostTerm(whereClause.Expression).SymbolRef.Symbol)

		andExpr := whereClause.Expression.Or[0].And[0]
		assert.Equal(t, "<", andExpr.Op)
		assert.Equal(t, float64(3), *andExpr.RHS.Summand.LHS.LHS.ConstantValue.Number)
	}
}

func TestSelectFieldsWithAliasQueryBasicPasses(t *testing.T) {
	result := &SelectQuery{}
	err := selectQueryParser.ParseString("SELECT A as B FROM target1", result)
	if assert.NoError(t, err) {
		selectExpr := result.SelectClause.Expression
		assert.False(t, selectExpr.All)
		assert.Len(t, selectExpr.Expressions, 1)
		assert.Equal(t, "A", exprLeftmostTerm(selectExpr.Expressions[0].Expression).SymbolRef.Symbol)
		assert.Equal(t, "B", selectExpr.Expressions[0].As)
	}

	result = &SelectQuery{}
	err = selectQueryParser.ParseString("SELECT X as Y, noAlia, longername AS short FROM target1", result)
	if assert.NoError(t, err) {
		selectExpr := result.SelectClause.Expression
		assert.False(t, selectExpr.All)
		assert.Len(t, selectExpr.Expressions, 3)
		assert.Equal(t, "X", exprLeftmostTerm(selectExpr.Expressions[0].Expression).SymbolRef.Symbol)
		assert.Equal(t, "Y", selectExpr.Expressions[0].As)

		assert.Equal(t, "noAlia", exprLeftmostTerm(selectExpr.Expressions[1].Expression).SymbolRef.Symbol)
		assert.Equal(t, "", selectExpr.Expressions[1].As)

		assert.Equal(t, "longername", exprLeftmostTerm(selectExpr.Expressions[2].Expression).SymbolRef.Symbol)
		assert.Equal(t, "short", selectExpr.Expressions[2].As)
	}
}

func TestSelectFieldsQueryBasicPasses(t *testing.T) {
	result := &SelectQuery{}
	err := selectQueryParser.ParseString("SELECT abc FROM target1", result)
	if assert.NoError(t, err) {
		selectExpr := result.SelectClause.Expression
		assert.False(t, selectExpr.All)

		assert.Len(t, selectExpr.Expressions, 1)
		assert.Equal(t, "abc", exprLeftmostTerm(selectExpr.Expressions[0].Expression).SymbolRef.Symbol)
		assert.Equal(t, "", selectExpr.Expressions[0].As)

		assert.Equal(t, "target1", result.FromClause.FromExpression.Target)
		assert.Nil(t, result.FromClause.FromExpression.WhereClause)
	}

	result = &SelectQuery{}
	err = selectQueryParser.ParseString("SELECT abc, def1, g__hi FROM abcd_2", result)
	if assert.NoError(t, err) {
		selectExpr := result.SelectClause.Expression
		assert.False(t, selectExpr.All)

		assert.Len(t, selectExpr.Expressions, 3)
		assert.Equal(t, "abc", exprLeftmostTerm(selectExpr.Expressions[0].Expression).SymbolRef.Symbol)
		assert.Equal(t, "def1", exprLeftmostTerm(selectExpr.Expressions[1].Expression).SymbolRef.Symbol)
		assert.Equal(t, "g__hi", exprLeftmostTerm(selectExpr.Expressions[2].Expression).SymbolRef.Symbol)

		assert.Equal(t, "abcd_2", result.FromClause.FromExpression.Target)
		assert.Nil(t, result.FromClause.FromExpression.WhereClause)
	}
}

func TestIntConstantValue(t *testing.T) {
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

func TestStringConstantValue(t *testing.T) {
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

func TestBoolConstantValue(t *testing.T) {
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

func TestConstantValueFails(t *testing.T) {
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
