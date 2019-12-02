package aocday

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func part1(t *testing.T, in string) string {
	var offsets []int
	for _, s := range strings.Fields(in) {
		o, err := strconv.Atoi(s)
		if err != nil {
			t.Fatalf("strconv(%q): %s", s, err)
		}
		offsets = append(offsets, o)
	}

	var freq int
	for _, o := range offsets {
		freq += o
	}
	return fmt.Sprint(freq)
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"example1", "+1 +1 +1", "3"},
		{"example2", "+1 +1 -2", "0"},
		{"example3", "-1 -2 -3", "-6"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part1(t, test.in), test.want; got != want {
				t.Errorf("part1(%#v) = %#v, want %#v", test.in, got, want)
			}
		})
	}

	t.Run("part1", func(t *testing.T) {
		t.Logf("part1: %s", part1(t, read(t, "input.txt")))
	})
}
