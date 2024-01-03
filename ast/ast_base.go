package ast

// Node The base interface
type Node interface {
	TokenLiteral() string
	String() string
}

// StatementNode All statement nodes implement this interface
type StatementNode interface {
	Node
	statementNode()
}

// ExpressionNode All expression nodes implement this interface
type ExpressionNode interface {
	Node
	expressionNode()
}
