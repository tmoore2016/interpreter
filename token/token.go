// interpreter\token\token.go

package token

// TokenType Create a type of token
type TokenType string

// Token TokenType is a string
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
	ASSIGN = "="
	PLUS   = "+"

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
)

// input for keywords
var keywords = map[string]TokenType{
	"fn":  FUNCTION,
	"let": LET,
}

// LookupIdent determines whether identifier is a keyword
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
