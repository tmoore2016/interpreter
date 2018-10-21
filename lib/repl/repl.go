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

	"github.com/tmoore2016/interpreter/lib/evaluator"
	"github.com/tmoore2016/interpreter/lib/lexer"
	"github.com/tmoore2016/interpreter/lib/object"
	"github.com/tmoore2016/interpreter/lib/parser"
)

// PROMPT = command prompt
const PROMPT = ">> "

// Start REPL: Read, Evaluate, Print, Loop
// Read from the input source until newline, pass the string to lexer, parse the lexer output, print the AST, evaluate the AST and print the eval.
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		// Lex the input and write the parsed output line by line
		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		// If there are parser errors, print the errors
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		// Evaluate the input and write as output
		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, program.String())
			io.WriteString(out, "\n")
		}

		io.WriteString(out, evaluated.Inspect())
		io.WriteString(out, "\n")
	}
}

// printParserErrors writes any parser errors found
func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Uh oh, parser error(s) detected:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
