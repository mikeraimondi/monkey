package parser

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/mikeraimondi/monkey/ast"
	"github.com/mikeraimondi/monkey/lexer"
)

func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 83838;
`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatal("ParseProgram() returned nil")
	}
	if l := len(program.Statements); l != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got %d", l)
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("'let' statement %s", tt.expectedIdentifier), func(t *testing.T) {
			stmt := program.Statements[i]
			if !testLetStatement(t, stmt, tt.expectedIdentifier) {
				return
			}
		})
	}
}

func TestReturnStatements(t *testing.T) {
	input := `
return 5;
return 10;
return 993322;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if l := len(program.Statements); l != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got %d", l)
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.returnStatement. got %T", stmt)
			continue
		}
		if r := returnStmt.TokenLiteral(); r != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q", r)
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	expected := "foobar"
	l := lexer.New(expected + ";")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if l := len(program.Statements); l != 1 {
		t.Fatalf("program does not have enough statements. got %d", l)
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got %T",
			program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("expression is not ast.Identifier. got %T", stmt.Expression)
	}
	if v := ident.Value; v != expected {
		t.Errorf("ident.Value not %s. got %s", expected, v)
	}
	if tl := ident.TokenLiteral(); tl != expected {
		t.Errorf("ident.TokenLiteral not %s. got %s", expected, tl)
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	expected := "5"

	l := lexer.New(expected + ";")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if l := len(program.Statements); l != 1 {
		t.Fatalf("program does not have enough statements. got %d", l)
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got %T",
			program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral. got %T", stmt.Expression)
	}
	expectedInt, err := strconv.ParseInt(expected, 0, 64)
	if err != nil {
		t.Fatalf("cannot parse expected value into integer")
	}
	if literal.Value != expectedInt {
		t.Fatalf("literal.Value not %d. got %d", expectedInt, literal.Value)
	}
	if tl := literal.TokenLiteral(); tl != expected {
		t.Errorf("literal.TokenLiteral not %s. got %s", expected, tl)
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if l := len(program.Statements); l != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got %d", 1, l)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got %T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got %T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got %s",
				tt.operator, exp.Operator)
		}
		if !testIntegerLiteral(t, exp.Right, tt.integerValue) {
			return
		}
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got %q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got %T", s)
		return false
	}

	if n := letStmt.Name.Value; n != name {
		t.Errorf("letStmt.Name.Value not %s. got %s", name, n)
		return false
	}

	if n := letStmt.Name.TokenLiteral(); n != name {
		t.Errorf("letStmt.Name.TokenLiteral() not %s. got %s", name, n)
		return false
	}

	return true
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got %T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d. got %d", value, integ.Value)
		return false
	}

	if tl := integ.TokenLiteral(); tl != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got %s", value, tl)
		return false
	}

	return true
}
