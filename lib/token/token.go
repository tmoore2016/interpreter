/*
Token package for
Doorkey, a Monkey Derivative
by Travis Moore
By following "Writing an Interpreter in Go" by Thorsten Ball, https://interpreterbook.com/
*/

// interpreter\token\token.go

package token

// TokenType Create token types
type TokenType string

// Token literal value
type Token struct {
	Type    TokenType
	Literal string
}

// Constants
const (
	ILLEGAL = "ILLEGAL" // Unknown Token/Character
	EOF     = "EOF"     // End of file

	// Identifiers and literals
	IDENT = "IDENT" // Name
	INT   = "INT"   // Integers

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	NOT      = "!"
	MULTIPLY = "*"
	DIVIDE   = "/"
	LT       = "<"
	GT       = ">"
	EQ       = "=="
	NOT_EQ   = "!="

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)

// input for keywords
var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

// LookupIdent determines whether identifier is a keyword
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
