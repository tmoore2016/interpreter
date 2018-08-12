/*
Lexer for
Doorkey, The Monkey Derivative
by Travis Moore
By following "Writing an Interpreter in Go" by Thorsten Ball, https://interpreterbook.com/
*/

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

// readChar returns the last char in the input string and increments to the next one until there are none left.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) { // ! = no input or end of string
		l.ch = 0 // nill
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition // always points to the length char read
	l.readPosition++            // always points to the next char
}

// peekChar returns the next char in the input string, but doesn't increment the position
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

// NextToken looks to see which is called
// Could be a Loop that calls a text file
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	// Initialize skipping whitespace
	l.skipWhitespace()

	// this can be generalized
	// the char determines the token type
	switch l.ch {
	// Two character checks could be abstracted into a function
	// '=' or '=='
	case '=':
		// call peekChar() to check for a second '='
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '/':
		tok = newToken(token.DIVIDE, l.ch)
	case '*':
		tok = newToken(token.MULTIPLY, l.ch)
	// Two character checks could be abstracted into a function
	// '!' or '!='
	case '!':
		// call peekChar() to check for a '='
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			tok = newToken(token.NOT, l.ch)
		}
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)

	//case '':
	//	tok = newToken(token.ASSIGN, )

	// Nill, end of file
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF

	// default makes a check whenever l.ch is unrecognized
	// If token is letter, get its literal and type (could be a keyword), otherwise throw error
	default:
		if isLetter(l.ch) { // if length character is letter
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

// skipWhitespace ignores white space
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar() // \r = return
	}
}

/*
Advance lexer for each type, could be generalized with a loop
*/

// advances the lexer's position until it encounters a non-letter char
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// advances the lexer's position until it encounters a non-number char
func (l *Lexer) readNumber() string {
	position := l.position // match indexes
	// for
	for isDigit(l.ch) { // if index  is a digit move to next char
		l.readChar() // send length non-digit to readChar()
	}
	return l.input[position:l.position]
}

/*
Booleans for token types
*/

// returns true if arg (token) is a letter, _ and $ are letters for var names
// Too many ors?
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch == '$'
}

// returns true if character is a digit, 0-9
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// returns true if character is one-character token

// initialize the tokens
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
