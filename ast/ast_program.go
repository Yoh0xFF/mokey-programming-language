package ast

import "bytes"

// Top level program node
type ProgramNode struct {
	Statements []StatementNode
}

func (p *ProgramNode) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *ProgramNode) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}
