package ast

import "bytes"

// ProgramNode Top level program node
type ProgramNode struct {
	StatementNodes []StatementNode
}

func (p *ProgramNode) TokenLiteral() string {
	if len(p.StatementNodes) > 0 {
		return p.StatementNodes[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *ProgramNode) String() string {
	var out bytes.Buffer

	for _, s := range p.StatementNodes {
		out.WriteString(s.String())
	}

	return out.String()
}
