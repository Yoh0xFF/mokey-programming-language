package ast

import (
	"bytes"
	"monkey/token"
)

// Statements

/* Let statement ast node */
type LetStatementNode struct {
	Token token.Token // the 'let' token
	Name  Identifier
	Value ExpressionNode
}

func (ls *LetStatementNode) statementNode()       {}
func (ls *LetStatementNode) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatementNode) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

/* Return statement ast node */
type ReturnStatementNode struct {
	Token       token.Token // the 'return' token
	ReturnValue ExpressionNode
}

func (rs *ReturnStatementNode) statementNode()       {}
func (rs *ReturnStatementNode) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatementNode) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

/* Expession statement ast node */
type ExpressionStatementNode struct {
	Token      token.Token // the first token of the expression
	Expression ExpressionNode
}

func (es *ExpressionStatementNode) statementNode()       {}
func (es *ExpressionStatementNode) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatementNode) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}

/* Block statement ast node */
type BlockStatementNode struct {
	Token      token.Token // the { token
	Statements []StatementNode
}

func (bs *BlockStatementNode) statementNode()       {}
func (bs *BlockStatementNode) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatementNode) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}
