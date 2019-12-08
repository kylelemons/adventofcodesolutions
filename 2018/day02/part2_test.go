package aocday

import (
	"strings"
	"testing"
)

func part2(t *testing.T, in string) string {
	boxes := strings.Fields(in)
	for i, box1 := range boxes {
	nextBox:
		for _, box2 := range boxes[i+1:] {
			if len(box1) != len(box2) {
				continue
			}
			var diffs, pos int
			for i := range box1 {
				if c1, c2 := box1[i], box2[i]; c1 != c2 {
					diffs++
					if diffs > 1 {
						continue nextBox
					}
					pos = i
				}
			}
			if diffs != 1 {
				continue
			}
			return box1[:pos] + box1[pos+1:]
		}
	}
	return ""
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"example1", "abcde fghij klmno pqrst fguij axcye wvxyz", "fgij"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got != want {
				t.Errorf("part2(%#v) = %#v, want %#v", test.in, got, want)
			}
		})
	}
	if t.Failed() {
		return
	}

	t.Run("part2", func(t *testing.T) {
		t.Logf("part2: %v", part2(t, read(t, "input.txt")))
	})
}
