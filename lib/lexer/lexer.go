/*
Lexer for
Doorkey, a Monkey Derivative
by Travis Moore
By following "Writing an Interpreter in Go" by Thorsten Ball, https://interpreterbook.com/
*/

// lexer/lexer.go

// Exercise: Support Unicode, "So it's left as an exercise to the reader to fully support Unicode (and emojis!) in Monkey."

package lexer

import "github.com/tmoore2016/interpreter/lib/token"

// Lexer for input and pointers
type Lexer struct {
	input        string
	position     int  // current lexer position (points to current ch)
	readPosition int  // current reading position in input (after current ch). Enables Peek?
	ch           byte // current char being examined
}

// New calls *Lexer's readChar before NextToken is called and initializes pointers
func New(input string) *Lexer { // Call new input, prepare Lexer
	l := &Lexer{input: input} // Create Lexer instance with input
	l.readChar()              // Initialize Lexer pointer
	return l                  // when all input is lexed
}

// readChar reads each char in the input string. The read pointer's position is always one ahead of the Lexer pointer's position, unless there are 0 chars left
func (l *Lexer) readChar() {

	if l.readPosition >= len(l.input) { // If greater than 0, Lexer's read position keeps incrementing until it is beyond input length.
		l.ch = 0 // Lexer char is 0, nil?.

	} else {
		l.ch = l.input[l.readPosition] // lexer char is lexer's read position from input
	}

	l.position = l.readPosition // Lexer's char position advances to lexer's read position
	l.readPosition++            // Lexer's read pointer advances to the next input char
}

// peekChar returns the next char in the input string (the read char), but doesn't increment the position
func (l *Lexer) peekChar() byte {

	if l.readPosition >= len(l.input) { // If Lexer's read position is beyond the input length
		return 0 // No peek char

	} else {
		return l.input[l.readPosition] // Send the lexer's read position to the lexer as input
	}
}

// NextToken looks to see which is called
// Could be a Loop that calls a text file
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	// Initialize skipping whitespace
	l.skipWhitespace()

	// this can be generalized
	// Lexer's char determines the token type
	switch l.ch {

	// Assign tokens to input
	// Dual character checks could be abstracted into a function
	// '=' or '=='
	case '=': // Lexer's char is "="
		if l.peekChar() == '=' { // Call peekChar to check if the next char is another "="
			ch := l.ch                                          // Advance lexer character to the second "="
			l.readChar()                                        // Call lexer read char to advance pointer again
			literal := string(ch) + string(l.ch)                // The literal value of the char after the second "="
			tok = token.Token{Type: token.EQ, Literal: literal} // Assign the current char's literal value to EQ token
		} else {
			tok = newToken(token.ASSIGN, l.ch) // Assign the
		}
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
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	//case '':
	//	tok = newToken(token.ASSIGN, )

	// Empty or end of file
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF

	// default makes a check whenever l.ch is unrecognized, none of the cases above
	// If token is letter or digit, get type and literal value, otherwise throw error
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
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' { // \r = return
		l.readChar() // Advance lexer pointers
	}
}

/*
Advance lexer for each type, could be generalized with a loop
*/

// readIdentifier reads identifiers (names, words, chars). Advances the lexer's position until something other than a letter is encountered
func (l *Lexer) readIdentifier() string {
	position := l.position // Pointer position is lexer's position
	for isLetter(l.ch) {   // For each lexer char that is a letter,
		l.readChar() // Read and advance
	}
	return l.input[position:l.position] // Send lexer new position input
}

// advances the lexer's position until it encounters a non-number char
func (l *Lexer) readNumber() string {
	position := l.position // match indexes
	// for
	for isDigit(l.ch) { // for each lexer position that is a digit,
		l.readChar() // advance
	}
	return l.input[position:l.position] // Send lexer new position input
}

// Advances the lexer until it encounters a closing " or EOF. Previous characters are part of a string.
// Add error reporting and character escaping ("hello \"world\"")
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

/*
Booleans for token types
*/

// returns true if token is 1 byte string, _ and $ are letters for var names
// Too many ors?
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch == '$'
}

// returns true if character is a digit, 0-9
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// returns true if character is one-character token

// initialize the tokens, they are 1 byte Type string
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
