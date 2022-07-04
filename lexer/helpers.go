package lexer

import "monkey/token"

func isLetter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func newToken[T byte | string](tokenType token.TokenType, ch T) token.Token {
	var literal string

	switch any(ch).(type) {
	case string:
		literal = any(ch).(string)
	case byte:
		literal = string(ch)
	}

	return token.Token{Type: tokenType, Literal: literal}
}
