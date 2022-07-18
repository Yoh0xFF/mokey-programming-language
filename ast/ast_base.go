package ast

// The base Node interface
type Node interface {
	TokenLiteral() string
	String() string
}

// All statement nodes implement this interface
type StatementNode interface {
	Node
	statementNode()
}

// All expression nodes implement this interface
type ExpressionNode interface {
	Node
	expressionNode()
}
