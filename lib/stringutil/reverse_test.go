// This is a test for reverse.go using the Go testing package
package stringutil

import "testing"

func TestReverse(t *testing.T) {
    cases := []struct {
        in, want string
    }{
        {"Welcome to Doorkey", "yekrooD ot emocleW"},
        {"Hello, World", "dlroW ,olleH"},
        {"", ""},
    }
    for _, c:= range cases {
        got := Reverse(c.in)
        if got != c.want {
            t.Errorf("Reverse(%q) == %q, want %q", c.in, got, c.want)
        }
    }
}