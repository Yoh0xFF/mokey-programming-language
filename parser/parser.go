package parser

import (
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

const (
	_ int = iota
	LOWEST
	EQUALS       // ==
	LESS_GREATER // > or <
	SUM          // +
	PRODUCT      // *
	PREFIX       // -x or !x
	CALL         // myFn(x)
)

var precedenceMap = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESS_GREATER,
	token.GT:       LESS_GREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
}

type (
	prefixParseFn func() ast.ExpressionNode
	infixParseFn  func(ast.ExpressionNode) ast.ExpressionNode
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFnMap map[token.TokenType]prefixParseFn
	infixParseFnMap  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFnMap = make(map[token.TokenType]prefixParseFn)
	p.registerPrefixParserFn(token.IDENT, p.parseIdentifier)
	p.registerPrefixParserFn(token.INT, p.parseIntegerLiteral)
	p.registerPrefixParserFn(token.BANG, p.parsePrefixExpression)
	p.registerPrefixParserFn(token.MINUS, p.parsePrefixExpression)
	p.registerPrefixParserFn(token.TRUE, p.parseBoolean)
	p.registerPrefixParserFn(token.FALSE, p.parseBoolean)
	p.registerPrefixParserFn(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefixParserFn(token.IF, p.parseIfExpression)
	p.registerPrefixParserFn(token.FUNCTION, p.parseFunctionLiteral)

	p.infixParseFnMap = make(map[token.TokenType]infixParseFn)
	p.registerInfixParserFn(token.PLUS, p.parseInfixExpression)
	p.registerInfixParserFn(token.MINUS, p.parseInfixExpression)
	p.registerInfixParserFn(token.SLASH, p.parseInfixExpression)
	p.registerInfixParserFn(token.ASTERISK, p.parseInfixExpression)
	p.registerInfixParserFn(token.EQ, p.parseInfixExpression)
	p.registerInfixParserFn(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfixParserFn(token.LT, p.parseInfixExpression)
	p.registerInfixParserFn(token.GT, p.parseInfixExpression)
	p.registerInfixParserFn(token.LPAREN, p.parseCallExpression)

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) ParseProgram() *ast.ProgramNode {
	program := &ast.ProgramNode{}
	program.StatementNodes = []ast.StatementNode{}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.StatementNodes = append(program.StatementNodes, stmt)
		}
		p.nextToken()
	}

	return program
}
