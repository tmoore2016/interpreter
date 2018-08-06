/*
Reverse.go, a utility for the Go library
by Travis Moore
By following "How to Write Go Code", Golang documentation
*/

// Package stringutil contains the utility functions for working with strings
package stringutil

// Reverse returns its argument string reversed rune-wise left to right.
func Reverse(s string) string {
	r := []rune(s)

	for i, j := 0, len(r)-1; // Counting backwards  := for non-typed vars

	i < len(r)/2; // i can't go below 0

	i, j = i+1, j-1 { // i = j's opposite
		r[i], r[j] = r[j], r[i]
	}

	return string(r)
}
