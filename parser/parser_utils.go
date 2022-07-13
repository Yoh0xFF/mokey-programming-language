package parser

import (
	"fmt"
	"monkey/token"
)

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("Expected next token to be %s, but got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixParserFnFound(t token.TokenType) {
	msg := fmt.Sprintf("No prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) registerPrefixParserFn(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFnMap[tokenType] = fn
}

func (p *Parser) registerInfixParserFn(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFnMap[tokenType] = fn
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedenceMap[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedenceMap[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}
