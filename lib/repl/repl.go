/*
Read, Evaluate, Print, Loop (REPL) for
Doorkey, a Monkey Derivative
by Travis Moore
By following "Writing an Interpreter in Go" by Thorsten Ball, https://interpreterbook.com/
*/

package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/tmoore2016/interpreter/lib/lexer"
	"github.com/tmoore2016/interpreter/lib/parser"
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
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		io.WriteString(out, program.String())
		io.WriteString(out, "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Uh oh, parser error(s) detected:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
