package ast

import (
	"github.com/mikeraimondi/monkey/token"
)

// Node is a node in the AST
type Node interface {
	TokenLiteral() string
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

// LetStatement is "let"
type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

// Identifier is an expression because sometimes identifiers produce values
type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
