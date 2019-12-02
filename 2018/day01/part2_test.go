package aocday

import (
	"strconv"
	"strings"
	"testing"
)

func part2(t *testing.T, in string) int {
	var offsets []int
	for _, s := range strings.Fields(in) {
		o, err := strconv.Atoi(s)
		if err != nil {
			t.Fatalf("strconv(%q): %s", s, err)
		}
		offsets = append(offsets, o)
	}

	var freq int
	seen := map[int]bool{}
	for i := 0; i < 1000; i++ {
		for _, o := range offsets {
			freq += o
			if seen[freq] {
				return freq
			}
			seen[freq] = true
		}
	}
	panic("loop range exceeded")
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"example1", "+1 -1", 0},
		{"example2", "+3 +3 +4 -2 -4", 10},
		{"example3", "-6 +3 +8 +5 -6", 5},
		{"example3", "+7 +7 -2 -7 -4", 14},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got != want {
				t.Errorf("part2(%#v) = %#v, want %#v", test.in, got, want)
			}
		})
	}

	t.Run("part2", func(t *testing.T) {
		t.Logf("part2: %v", part2(t, read(t, "input.txt")))
	})
}
