// lexer/lexer.go

// Exercise: Support Unicode, "So it's left as an exercise to the reader to fully support Unicode (and emojis!) in Monkey."

package lexer

import "github.com/tmoore2016/interpreter/token"

// Lexer for input and pointers
type Lexer struct {
	input        string
	position     int  // current reading position in input (points to current ch)
	readPosition int  //current reading position in input (after current ch)
	ch           byte // current char being examined
}

// New calls readChar from *Lexer before NextToken is called and initializes pointers
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// find the next char in the input string
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) { // no input or end of string
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition // always points to the last char read
	l.readPosition++            // always points to the next char
}

// NextToken looks to see which is called
// Could be a Loop that calls a text file
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	switch l.ch {
	case '=':
		tok = newToken(token.ASSIGN, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	// Null, end of file
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
		//case '':
		//	tok = newToken(token.ASSIGN, )
	}

	l.readChar()
	return tok
}

// initialize the tokens
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
