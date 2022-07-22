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
	Token     token.Token // The prefix token, e.g. !
	Operator  string
	RightNode ExpressionNode
}

func (pe *PrefixExpressionNode) expressionNode()      {}
func (pe *PrefixExpressionNode) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpressionNode) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.RightNode.String())
	out.WriteString(")")

	return out.String()
}

/* Infix expression ast node */
type InfixExpressionNode struct {
	Token     token.Token // The operator token, e.g. +
	LeftNode  ExpressionNode
	Operator  string
	RightNode ExpressionNode
}

func (ie *InfixExpressionNode) expressionNode()      {}
func (ie *InfixExpressionNode) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpressionNode) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.LeftNode.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.RightNode.String())
	out.WriteString(")")

	return out.String()
}

/* If expression ast node */
type IfExpressionNode struct {
	Token           token.Token // The 'if' token
	ConditionNode   ExpressionNode
	ConsequenceNode *BlockStatementNode
	AlternativeNode *BlockStatementNode
}

func (ie *IfExpressionNode) expressionNode()      {}
func (ie *IfExpressionNode) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpressionNode) String() string {
	var out bytes.Buffer

	out.WriteString("if ")
	out.WriteString(ie.ConditionNode.String())
	out.WriteString(" ")
	out.WriteString(ie.ConsequenceNode.String())

	if ie.AlternativeNode != nil {
		out.WriteString(" else ")
		out.WriteString(ie.AlternativeNode.String())
	}

	return out.String()
}

/* Function literal ast node */
type FunctionLiteralNode struct {
	Token      token.Token // The 'fn' token
	ParamNodes []*IdentifierNode
	BodyNode   *BlockStatementNode
}

func (fl *FunctionLiteralNode) expressionNode()      {}
func (fl *FunctionLiteralNode) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteralNode) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fl.ParamNodes {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.BodyNode.String())

	return out.String()
}

/* Function call expression ast node */
type CallExpressionNode struct {
	Token    token.Token    // The '(' token
	FnNode   ExpressionNode // Identifier or FunctionLiteral
	ArgNodes []ExpressionNode
}

func (ce *CallExpressionNode) expressionNode()      {}
func (ce *CallExpressionNode) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpressionNode) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.ArgNodes {
		args = append(args, a.String())
	}

	out.WriteString(ce.FnNode.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

/* String literal expression ast node */
type StringLiteralNode struct {
	Token token.Token
	Value string
}

func (sl *StringLiteralNode) expressionNode()      {}
func (sl *StringLiteralNode) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteralNode) String() string       { return sl.Value }

/* Array literal expression ast node */
type ArrayLiteralNode struct {
	Token    token.Token // The '[' token
	Elements []ExpressionNode
}

func (al *ArrayLiteralNode) expressionNode()      {}
func (al *ArrayLiteralNode) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteralNode) String() string {
	var out bytes.Buffer

	strElements := []string{}
	for _, el := range al.Elements {
		strElements = append(strElements, el.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(strElements, ", "))
	out.WriteString("]")

	return out.String()
}

/* Index expression ast node */
type IndexExpressionNode struct {
	Token token.Token // The '[' token
	Left  ExpressionNode
	Index ExpressionNode
}

func (ie *IndexExpressionNode) expressionNode()      {}
func (ie *IndexExpressionNode) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpressionNode) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")

	return out.String()
}

/* Hash literal ast node */
type HashLiteralNode struct {
	Token token.Token // The '{' token
	Pairs map[ExpressionNode]ExpressionNode
}

func (hl *HashLiteralNode) expressionNode()      {}
func (hl *HashLiteralNode) TokenLiteral() string { return hl.Token.Literal }
func (hl *HashLiteralNode) String() string {
	var out bytes.Buffer

	strPairs := []string{}
	for key, val := range hl.Pairs {
		strPairs = append(strPairs, key.String()+" : "+val.String())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(strPairs, ", "))
	out.WriteString("}")

	return out.String()
}
