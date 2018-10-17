/*
Parser tracer for
Doorkey, a Monkey Derivative
by Travis Moore
By following "Writing an Interpreter in Go" by Thorsten Ball, https://interpreterbook.com/
*/

// Parser_tracer traces parser functions
// go test -v -run TestOperatorPrecedenceParsing ./lib/parser

package parser

import (
	"fmt"
	"strings"
)

var traceLevel int = 0

// placeholder string for identLevel
const traceIdentPlaceholder string = "\t"

// go through traceLevel until it is nil
func identLevel() string {
	return strings.Repeat(traceIdentPlaceholder, traceLevel-1)
}

// print parser strings, level #
func tracePrint(fs string) {
	fmt.Printf("%s%s\n", identLevel(), fs)
}

// increment tracelevel
func incIdent() {
	traceLevel = traceLevel + 1
}

// decrement tracelevel
func decIdent() {
	traceLevel = traceLevel - 1
}

func trace(msg string) string {
	incIdent()
	tracePrint("BEGIN " + msg)
	return msg
}

func untrace(msg string) {
	tracePrint("END" + msg)
	decIdent()
}
