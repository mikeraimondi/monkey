package ast

import (
	"strings"

	"github.com/mikeraimondi/monkey/token"
)

// Node is a node in the AST
type Node interface {
	TokenLiteral() string
	String() string
}

// Statement is evaluated but produces no value
type Statement interface {
	Node
	statementNode()
}

// Expression is evaluated and produces a value
type Expression interface {
	Node
	expressionNode()
}

// Program is the root node of the AST
type Program struct {
	Statements []Statement
}

// TokenLiteral is used for debugging
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	out := StringBuilder{}

	for _, s := range p.Statements {
		out.MustWrite(s.String())
	}

	return out.String()
}

// LetStatement is "let"
type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}

// TokenLiteral is used for debugging
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

func (ls *LetStatement) String() string {
	out := StringBuilder{}

	out.MustWrite(ls.TokenLiteral() + " ")
	out.MustWrite(ls.Name.String())
	out.MustWrite(" = ")

	if ls.Value != nil {
		out.MustWrite(ls.Value.String())
	}

	out.MustWrite(";")

	return out.String()
}

// Identifier is an expression because sometimes identifiers produce values
type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode() {}

// TokenLiteral is used for debugging
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

// ReturnStatement is "return"
type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}

// TokenLiteral is used for debugging
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

func (rs *ReturnStatement) String() string {
	out := StringBuilder{}

	out.MustWrite(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.MustWrite(rs.ReturnValue.String())
	}

	out.MustWrite(";")

	return out.String()
}

// ExpressionStatement is a statement that produces a value
type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}

// TokenLiteral is used for debugging
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// PrefixExpression is an expression with an operator that precedes an operand
type PrefixExpression struct {
	Token    token.Token // the prefix token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}

// TokenLiteral is used for debugging
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	out := StringBuilder{}

	out.MustWrite("(")
	out.MustWrite(pe.Operator)
	out.MustWrite(pe.Right.String())
	out.MustWrite(")")

	return out.String()
}

// InfixExpression is an expression with an operator between two operands
type InfixExpression struct {
	Token    token.Token // the operator token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}

// TokenLiteral is used for debugging
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	out := StringBuilder{}

	out.MustWrite("(")
	out.MustWrite(ie.Left.String())
	out.MustWrite(" " + ie.Operator + " ")
	out.MustWrite(ie.Right.String())
	out.MustWrite(")")

	return out.String()
}

// IntegerLiteral is an integer
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}

// TokenLiteral is used for debugging
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

// Boolean is a bool
type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode() {}

// TokenLiteral is used for debugging
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

// StringBuilder is a convenience wrapper around strings.Builder
type StringBuilder struct {
	strings.Builder
}

// MustWrite panics if WriteString is unsuccesful. This should "never happen".
func (builder *StringBuilder) MustWrite(in string) {
	if _, err := builder.WriteString(in); err != nil {
		panic("error writing string to builder: " + err.Error())
	}
}
