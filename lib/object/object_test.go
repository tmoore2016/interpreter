/*
Object_Test package for
Doorkey, a Monkey Derivative
by Travis Moore
By following "Writing an Interpreter in Go" by Thorsten Ball, https://interpreterbook.com/
*/

// interpreter\object\object_test.go

package object

import "testing"

// TestStringHashKey tests diffs of hash keys of strings, identical values should have the same hash keys.
func TestStringHashKey(t *testing.T) {
	book1 := &String{Value: "The Sea-Wolf"}
	book2 := &String{Value: "The Sea-Wolf"}
	author1 := &String{Value: "Jack London"}
	author2 := &String{Value: "Jack London"}

	if book1.HashKey() != book2.HashKey() {
		t.Errorf("Strings with the same content have different hash keys.")
	}

	if author1.HashKey() != author2.HashKey() {
		t.Errorf("Strings with the same content have different hash keys.")
	}

	if book1.HashKey() == author1.HashKey() {
		t.Errorf("Strings with different content have the same hash keys.")
	}
}

// TestIntHashKey tests diffs of hash keys with integer values, identical values should have the same hash keys.
func TestIntHashKey(t *testing.T) {
	index1 := &Integer{Value: 001}
	index2 := &Integer{Value: 001}
	year1 := &Integer{Value: 1904}
	year2 := &Integer{Value: 1904}

	if index1.HashKey() != index2.HashKey() {
		t.Errorf("Integers with the same value have different hash keys.")
	}

	if year1.HashKey() != year2.HashKey() {
		t.Errorf("Integers with the same value have different hash keys.")
	}

	if index1.HashKey() == year1.HashKey() {
		t.Errorf("Integers with different values have the same hash keys.")
	}
}

// TestBooleanHashKey tests diffs of hash keys with boolean values, identical values should have the same hash keys.
func TestBooleanHashKey(t *testing.T) {
	true1 := &Boolean{Value: true}
	true2 := &Boolean{Value: true}
	false1 := &Boolean{Value: false}
	false2 := &Boolean{Value: false}

	if true1.HashKey() != true2.HashKey() {
		t.Errorf("Booleans of the same value have different hash keys.")
	}

	if false1.HashKey() != false2.HashKey() {
		t.Errorf("Booleans of the same value have different hash keys.")
	}

	if true1.HashKey() == false1.HashKey() {
		t.Errorf("Booleans of different values have the same hash keys.")
	}
}
