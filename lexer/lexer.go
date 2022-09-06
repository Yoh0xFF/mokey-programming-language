package lexer

import "monkey/token"

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
}

func New(input string) *Lexer {
	l := &Lexer{input: input}

	l.readChar()

	return l
}

func (l *Lexer) NextToken() token.Token {
	var tk token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		tokenType := token.ASSIGN
		literal := string(l.ch)

		if l.peekChar() == '=' {
			tokenType = token.EQ
			literal = literal + string(l.readChar())
		}

		tk = newToken(token.TokenType(tokenType), literal)
	case '+':
		tk = newToken(token.PLUS, l.ch)
	case '-':
		tk = newToken(token.MINUS, l.ch)
	case '!':
		tokenType := token.BANG
		literal := string(l.ch)

		if l.peekChar() == '=' {
			tokenType = token.NOT_EQ
			literal = literal + string(l.readChar())
		}

		tk = newToken(token.TokenType(tokenType), literal)
	case '/':
		tk = newToken(token.SLASH, l.ch)
	case '*':
		tk = newToken(token.ASTERISK, l.ch)
	case '<':
		tk = newToken(token.LT, l.ch)
	case '>':
		tk = newToken(token.GT, l.ch)
	case ';':
		tk = newToken(token.SEMICOLON, l.ch)
	case ':':
		tk = newToken(token.COLON, l.ch)
	case ',':
		tk = newToken(token.COMMA, l.ch)
	case '{':
		tk = newToken(token.LBRACE, l.ch)
	case '}':
		tk = newToken(token.RBRACE, l.ch)
	case '(':
		tk = newToken(token.LPAREN, l.ch)
	case ')':
		tk = newToken(token.RPAREN, l.ch)
	case '[':
		tk = newToken(token.LBRACKET, l.ch)
	case ']':
		tk = newToken(token.RBRACKET, l.ch)
	case '"':
		tk = newToken(token.STRING, l.readString())
	case 0:
		tk = newToken(token.EOF, "")
	default:
		if isLetter(l.ch) {
			literal := l.readIdentifier()
			tk = newToken(token.CheckIsKeyword(literal), literal)
			return tk
		} else if isDigit(l.ch) {
			literal := l.readNumber()
			tk = newToken(token.INT, literal)
			return tk
		} else {
			tk = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tk
}

func (l *Lexer) readChar() byte {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition += 1

	return l.ch
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position

	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position

	for isDigit(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) readString() string {
	position := l.position + 1

	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}

	return l.input[position:l.position]
}
