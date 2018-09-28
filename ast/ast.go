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

// BlockStatement is a series of statements
type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}

// TokenLiteral is used for debugging
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	out := StringBuilder{}

	for _, s := range bs.Statements {
		out.MustWrite(s.String())
	}

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

// IfExpression is an expression with 'if'
type IfExpression struct {
	Token       token.Token // the 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode() {}

// TokenLiteral is used for debugging
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	out := StringBuilder{}

	out.MustWrite("if")
	out.MustWrite(ie.Condition.String())
	out.MustWrite(" ")
	out.MustWrite(ie.Consequence.String())

	if ie.Alternative != nil {
		out.MustWrite("else ")
		out.MustWrite(ie.Alternative.String())
	}

	return out.String()
}

// CallExpression is a function invocation
type CallExpression struct {
	Token     token.Token // The '(' token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}

// TokenLiteral is used for debugging
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	out := StringBuilder{}

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.MustWrite(ce.Function.String())
	out.MustWrite("(")
	out.MustWrite(strings.Join(args, ", "))
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

// FunctionLiteral is fn
type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode() {}

// TokenLiteral is used for debugging
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	out := StringBuilder{}

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.MustWrite(fl.TokenLiteral())
	out.MustWrite("(")
	out.MustWrite(strings.Join(params, ", "))
	out.MustWrite(")")
	out.MustWrite(fl.Body.String())

	return out.String()
}

// StringLiteral is a string
type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode() {}

// TokenLiteral is used for debugging
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

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
