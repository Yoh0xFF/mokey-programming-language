package ast

import (
	"monkey/token"
	"testing"
)

func TestString(t *testing.T) {
	program := ProgramNode{
		StatementNodes: []StatementNode{
			&LetStatementNode{
				Token: token.Token{Type: token.LET, Literal: "let"},

				NameNode: IdentifierNode{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},

				ValueNode: &IdentifierNode{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if program.String() != "let myVar = anotherVar;" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}
