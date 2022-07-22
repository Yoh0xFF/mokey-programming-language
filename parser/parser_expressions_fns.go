package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/token"
	"strconv"
)

func (p *Parser) parseIdentifier() ast.ExpressionNode {
	return &ast.IdentifierNode{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.ExpressionNode {
	il := &ast.IntegerLiteralNode{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("Could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	il.Value = value

	return il
}

func (p *Parser) parseStringLiteral() ast.ExpressionNode {
	return &ast.StringLiteralNode{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseBooleanLiteral() ast.ExpressionNode {
	return &ast.BooleanNode{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parsePrefixExpression() ast.ExpressionNode {
	expr := &ast.PrefixExpressionNode{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expr.RightNode = p.parseExpression(PREFIX)

	return expr
}

func (p *Parser) parseInfixExpression(leftNode ast.ExpressionNode) ast.ExpressionNode {
	expr := &ast.InfixExpressionNode{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		LeftNode: leftNode,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expr.RightNode = p.parseExpression(precedence)

	return expr
}

func (p *Parser) parseGroupedExpression() ast.ExpressionNode {
	p.nextToken()

	expr := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return expr
}

func (p *Parser) parseIfExpression() ast.ExpressionNode {
	expr := &ast.IfExpressionNode{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	expr.ConditionNode = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expr.ConsequenceNode = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		expr.AlternativeNode = p.parseBlockStatement()
	}

	return expr
}

func (p *Parser) parseFunctionLiteral() ast.ExpressionNode {
	lit := &ast.FunctionLiteralNode{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	lit.ParamNodes = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	lit.BodyNode = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.IdentifierNode {
	identifiers := []*ast.IdentifierNode{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident := &ast.IdentifierNode{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &ast.IdentifierNode{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseCallExpression(fnNode ast.ExpressionNode) ast.ExpressionNode {
	expr := &ast.CallExpressionNode{Token: p.curToken, FnNode: fnNode}

	expr.ArgNodes = p.parseExpressionList(token.RPAREN)

	return expr
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.ExpressionNode {
	exprNodes := []ast.ExpressionNode{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return exprNodes
	}

	p.nextToken()
	exprNodes = append(exprNodes, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		exprNodes = append(exprNodes, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return exprNodes
}

func (p *Parser) parseArrayLiteral() ast.ExpressionNode {
	array := &ast.ArrayLiteralNode{Token: p.curToken}

	array.Elements = p.parseExpressionList(token.RBRACKET)

	return array
}

func (p *Parser) parseHashLiteral() ast.ExpressionNode {
	hash := &ast.HashLiteralNode{Token: p.curToken}
	hash.Pairs = make(map[ast.ExpressionNode]ast.ExpressionNode)

	for !p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		key := p.parseExpression(LOWEST)

		if !p.expectPeek(token.COLON) {
			return nil
		}

		p.nextToken()
		value := p.parseExpression(LOWEST)

		hash.Pairs[key] = value

		if !p.peekTokenIs(token.RBRACE) && !p.expectPeek(token.COMMA) {
			return nil
		}
	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	return hash
}

func (p *Parser) parseIndexExpression(left ast.ExpressionNode) ast.ExpressionNode {
	expr := &ast.IndexExpressionNode{Token: p.curToken, Left: left}

	p.nextToken()
	expr.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return expr
}
