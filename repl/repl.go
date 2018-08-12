package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/tmoore2016/interpreter/lexer"
	"github.com/tmoore2016/interpreter/token"
)

// PROMPT = command prompt
const PROMPT = ">> "

// Start REPL: Read, Evaluate, Print, Loop
// Read from the input source until newline, pass the string to lexer for eval, print tokens
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
