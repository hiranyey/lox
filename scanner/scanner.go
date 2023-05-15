package scanner

import (
	"jlox/errors"
	"strconv"
)

type Scanner struct {
	source  string
	start   int
	current int
	line    int
}

func NewScanner(source string) *Scanner {
	return &Scanner{source: source, start: 0, current: 0, line: 1}
}

func (s *Scanner) ScanTokens() []Token {
	var tokens []Token
	for !s.isAtEnd() {
		s.start = s.current
		nextToken := s.scanToken()
		if nextToken.Type != UNKNOWN {
			tokens = append(tokens, nextToken)
		}
	}
	tokens = append(tokens, Token{Type: EOF, Lexeme: "", Literal: nil, Line: s.line})
	return tokens
}

func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != expected {
		return false
	}
	s.current++
	return true
}
func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) advance() byte {
	s.current++
	return s.source[s.current-1]
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return '\000'
	}
	return s.source[s.current]
}

func (s *Scanner) string() Token {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}
	if s.isAtEnd() {
		errors.Error(s.line, "Unterminated string.")
		panic("Unterminated string")
	}
	s.advance()
	return Token{Type: STRING, Lexeme: s.source[s.start+1 : s.current-1], Literal: s.source[s.start+1 : s.current-1], Line: s.line}
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return '\000'
	}
	return s.source[s.current+1]
}

func (s *Scanner) number() Token {
	for isDigit(s.peek()) {
		s.advance()
	}
	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()
	}
	for isDigit(s.peek()) {
		s.advance()
	}
	number, err := strconv.ParseFloat(s.source[s.start:s.current], 64)
	if err != nil {
		errors.Error(s.line, "Error parsing float")
		panic("Error parsing float")
	}
	return Token{Type: NUMBER, Lexeme: s.source[s.start:s.current], Literal: number, Line: s.line}
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isAlphaNumeric(c byte) bool {
	return isAlpha(c) || isDigit(c)
}

func (s *Scanner) identifier() Token {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}
	text := s.source[s.start:s.current]
	if tokenType, ok := keywords[text]; ok {
		return Token{Type: tokenType, Lexeme: text, Literal: nil, Line: s.line}
	} else {
		return Token{Type: IDENTIFIER, Lexeme: text, Literal: nil, Line: s.line}
	}
}

func (s *Scanner) scanToken() Token {
	next := s.advance()
	switch next {
	case '(':
		return Token{Type: LEFT_PAREN, Lexeme: "(", Literal: nil, Line: s.line}
	case ')':
		return Token{Type: RIGHT_PAREN, Lexeme: ")", Literal: nil, Line: s.line}
	case '{':
		return Token{Type: LEFT_BRACE, Lexeme: "{", Literal: nil, Line: s.line}
	case '}':
		return Token{Type: RIGHT_BRACE, Lexeme: "}", Literal: nil, Line: s.line}
	case ',':
		return Token{Type: COMMA, Lexeme: ",", Literal: nil, Line: s.line}
	case '.':
		return Token{Type: DOT, Lexeme: ".", Literal: nil, Line: s.line}
	case '-':
		return Token{Type: MINUS, Lexeme: "-", Literal: nil, Line: s.line}
	case '+':
		return Token{Type: PLUS, Lexeme: "+", Literal: nil, Line: s.line}
	case ';':
		return Token{Type: SEMICOLON, Lexeme: ";", Literal: nil, Line: s.line}
	case '*':
		return Token{Type: STAR, Lexeme: "*", Literal: nil, Line: s.line}

	case '!':
		if s.match('=') {
			return Token{Type: BANG_EQUAL, Lexeme: "!=", Literal: nil, Line: s.line}
		} else {
			return Token{Type: BANG, Lexeme: "!", Literal: nil, Line: s.line}
		}
	case '=':
		if s.match('=') {
			return Token{Type: EQUAL_EQUAL, Lexeme: "==", Literal: nil, Line: s.line}
		} else {
			return Token{Type: EQUAL, Lexeme: "=", Literal: nil, Line: s.line}
		}
	case '<':
		if s.match('=') {
			return Token{Type: LESS_EQUAL, Lexeme: "<=", Literal: nil, Line: s.line}
		} else {
			return Token{Type: LESS, Lexeme: "<", Literal: nil, Line: s.line}
		}
	case '>':
		if s.match('=') {
			return Token{Type: GREATER_EQUAL, Lexeme: ">=", Literal: nil, Line: s.line}
		} else {
			return Token{Type: GREATER, Lexeme: ">", Literal: nil, Line: s.line}
		}

	case '/':
		if s.match('/') {
			for !s.isAtEnd() && s.peek() != '\n' {
				s.advance()
			}
			return Token{Type: UNKNOWN, Lexeme: "", Literal: nil, Line: s.line}
		} else if s.match('*') {
			for !s.isAtEnd() && !(s.peek() == '*' && s.peekNext() == '/') {
				if s.peek() == '\n' {
					s.line++
				}
				s.advance()
			}
			if s.isAtEnd() {
				errors.Error(s.line, "Unterminated block comment")
				panic("Unterminated block comment")
			}
			s.advance()
			s.advance()
			return Token{Type: UNKNOWN, Lexeme: "", Literal: nil, Line: s.line}
		} else {
			return Token{Type: SLASH, Lexeme: "/", Literal: nil, Line: s.line}
		}

	case ' ':
		return Token{Type: UNKNOWN, Lexeme: "", Literal: nil, Line: s.line}
	case '\r':
		return Token{Type: UNKNOWN, Lexeme: "", Literal: nil, Line: s.line}
	case '\t':
		return Token{Type: UNKNOWN, Lexeme: "", Literal: nil, Line: s.line}
	case '\n':
		s.line++
		return Token{Type: UNKNOWN, Lexeme: "", Literal: nil, Line: s.line}

	case '"':
		return s.string()

	default:
		if isDigit(next) {
			return s.number()
		} else if isAlpha(next) {
			return s.identifier()
		} else {
			errors.Error(s.line, "Unexpected character.")
			panic("Unexpected character")
		}
	}
}
