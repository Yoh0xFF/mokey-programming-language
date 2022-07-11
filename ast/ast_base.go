package ast

// The base Node interface
type Node interface {
	TokenLiteral() string
	String() string
}

// All statement nodes implement this interface
type Statement interface {
	Node
	statementNode()
}

// All expression nodes implement this interface
type Expression interface {
	Node
	expressionNode()
}

// Statements

// Expressions
