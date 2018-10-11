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

	fmt.Printf("Hello %s! Welcome to Doorkey a Monkey derivative!\n", user.Username)
	fmt.Printf("I will lex your input, parse it, and return the Abstract Syntax Tree. I'm currently learning to evaluate the AST, and I'll give back what I know.")
	repl.Start(os.Stdin, os.Stdout)
}
