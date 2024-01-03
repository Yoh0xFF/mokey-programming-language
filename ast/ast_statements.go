package ast

import (
	"bytes"
	"monkey/token"
)

// Statements

// LetStatementNode Let statement ast node
type LetStatementNode struct {
	Token     token.Token // the 'let' token
	NameNode  IdentifierNode
	ValueNode ExpressionNode
}

func (ls *LetStatementNode) statementNode()       {}
func (ls *LetStatementNode) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatementNode) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.NameNode.String())
	out.WriteString(" = ")

	if ls.ValueNode != nil {
		out.WriteString(ls.ValueNode.String())
	}

	out.WriteString(";")

	return out.String()
}

// ReturnStatementNode Return statement ast node
type ReturnStatementNode struct {
	Token           token.Token // the 'return' token
	ReturnValueNode ExpressionNode
}

func (rs *ReturnStatementNode) statementNode()       {}
func (rs *ReturnStatementNode) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatementNode) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValueNode != nil {
		out.WriteString(rs.ReturnValueNode.String())
	}

	out.WriteString(";")

	return out.String()
}

// ExpressionStatementNode Expression statement ast node
type ExpressionStatementNode struct {
	Token          token.Token // the first token of the expression
	ExpressionNode ExpressionNode
}

func (es *ExpressionStatementNode) statementNode()       {}
func (es *ExpressionStatementNode) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatementNode) String() string {
	if es.ExpressionNode != nil {
		return es.ExpressionNode.String()
	}

	return ""
}

// BlockStatementNode Block statement ast node
type BlockStatementNode struct {
	Token          token.Token // the { token
	StatementNodes []StatementNode
}

func (bs *BlockStatementNode) statementNode()       {}
func (bs *BlockStatementNode) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatementNode) String() string {
	var out bytes.Buffer

	for _, s := range bs.StatementNodes {
		out.WriteString(s.String())
	}

	return out.String()
}
