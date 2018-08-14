package parser

import (
	"fmt"
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
