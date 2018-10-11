/*
Doorkey, a Monkey Derivative
by Travis Moore
By following "Writing an Interpreter in Go" by Thorsten Ball, https://interpreterbook.com/
*/

package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/tmoore2016/interpreter/lib/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf(" Hello %s,", user.Username)
	fmt.Printf("\n Welcome to Doorkey a Monkey derivative!\n I will lex your input, parse it, and return the Abstract Syntax Tree.\n I'm learning to evaluate your input, go ahead give it a try.\n")
	repl.Start(os.Stdin, os.Stdout)
}
