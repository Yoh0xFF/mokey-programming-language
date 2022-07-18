package parser

import (
	"monkey/ast"
	"monkey/token"
)

func (p *Parser) parseStatement() ast.StatementNode {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatementNode {
	stmt := &ast.LetStatementNode{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.NameNode = ast.IdentifierNode{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.ValueNode = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatementNode {
	stmt := &ast.ReturnStatementNode{Token: p.curToken}

	p.nextToken()

	stmt.ReturnValueNode = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatementNode {
	stmt := &ast.ExpressionStatementNode{Token: p.curToken}

	stmt.ExpressionNode = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStatementNode {
	block := &ast.BlockStatementNode{Token: p.curToken}
	block.StatementNodes = []ast.StatementNode{}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		statement := p.parseStatement()
		block.StatementNodes = append(block.StatementNodes, statement)
		p.nextToken()
	}

	return block
}
