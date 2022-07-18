package parser

import (
	"monkey/ast"
	"monkey/token"
)

func (p *Parser) parseExpression(precedence int) ast.ExpressionNode {
	prefixParseFn := p.prefixParseFnMap[p.curToken.Type]
	if prefixParseFn == nil {
		p.noPrefixParserFnFound(p.curToken.Type)
		return nil
	}
	leftExr := prefixParseFn()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infixParseFn := p.infixParseFnMap[p.peekToken.Type]
		if infixParseFn == nil {
			return leftExr
		}

		p.nextToken()
		leftExr = infixParseFn(leftExr)
	}

	return leftExr
}
