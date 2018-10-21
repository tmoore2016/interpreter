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
		panic(err) // End program with stack trace
	}

	fmt.Printf(" Hello %s,", user.Username)
	fmt.Printf(" Welcome to Doorkey a Monkey derivative!\n I can evaluate your input, go ahead and give me a try.\n")
	repl.Start(os.Stdin, os.Stdout)
}
