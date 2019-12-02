package aocday

import (
	"strings"
	"testing"
)

func part1(t *testing.T, in string) int {
	var doubles, triples int
	for _, s := range strings.Fields(in) {
		var letters [256]int
		for _, r := range s {
			letters[int(r)]++
		}
		var double, triple int
		for _, count := range letters {
			if count == 2 {
				double = 1
			}
			if count == 3 {
				triple = 1
			}
		}
		doubles += double
		triples += triple
	}
	return doubles * triples
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"example1", "abcdef bababc abbcde abcccd aabcdd abcdee ababab", 4 * 3},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part1(t, test.in), test.want; got != want {
				t.Errorf("part1(%#v) = %#v, want %#v", test.in, got, want)
			}
		})
	}
	if t.Failed() {
		return
	}

	t.Run("part1", func(t *testing.T) {
		t.Logf("part1: %v", part1(t, read(t, "input.txt")))
	})
}
