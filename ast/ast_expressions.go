package ast

import (
	"bytes"
	"monkey/token"
	"strings"
)

// Expressions

/* IdentifierNode ast node */
type IdentifierNode struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *IdentifierNode) expressionNode()      {}
func (i *IdentifierNode) TokenLiteral() string { return i.Token.Literal }
func (i *IdentifierNode) String() string       { return i.Value }

/* Identifier boolean node */
type BooleanNode struct {
	Token token.Token
	Value bool
}

func (b *BooleanNode) expressionNode()      {}
func (b *BooleanNode) TokenLiteral() string { return b.Token.Literal }
func (b *BooleanNode) String() string       { return b.Token.Literal }

/* Integer boolean node */
type IntegerLiteralNode struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteralNode) expressionNode()      {}
func (il *IntegerLiteralNode) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteralNode) String() string       { return il.Token.Literal }

/* Predix expression ast node */
type PrefixExpressionNode struct {
	Token    token.Token // The prefix token, e.g. !
	Operator string
	Right    ExpressionNode
}

func (pe *PrefixExpressionNode) expressionNode()      {}
func (pe *PrefixExpressionNode) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpressionNode) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

/* Infix expression ast node */
type InfixExpressionNode struct {
	Token    token.Token // The operator token, e.g. +
	Left     ExpressionNode
	Operator string
	Right    ExpressionNode
}

func (ie *InfixExpressionNode) expressionNode()      {}
func (ie *InfixExpressionNode) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpressionNode) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

/* If expression ast node */
type IfExpressionNode struct {
	Token       token.Token // The 'if' token
	Condition   ExpressionNode
	Consequence *BlockStatementNode
	Alternative *BlockStatementNode
}

func (ie *IfExpressionNode) expressionNode()      {}
func (ie *IfExpressionNode) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpressionNode) String() string {
	var out bytes.Buffer

	out.WriteString("if ")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString(" else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

/* Function literal ast node */
type FunctionLiteralNode struct {
	Token      token.Token // The 'fn' token
	Parameters []*IdentifierNode
	Body       *BlockStatementNode
}

func (fl *FunctionLiteralNode) expressionNode()      {}
func (fl *FunctionLiteralNode) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteralNode) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())

	return out.String()
}

/* function call expression ast node */
type CallExpressionNode struct {
	Token     token.Token    // The '(' token
	Function  ExpressionNode // Identifier or FunctionLiteral
	Arguments []ExpressionNode
}

func (ce *CallExpressionNode) expressionNode()      {}
func (ce *CallExpressionNode) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpressionNode) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}
