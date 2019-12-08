package aocday

import (
	"strings"
	"testing"
)

func part2(t *testing.T, in []string) int {
	var fabric [1000][1000][]int

	for _, claim := range in {
		// #1212 @ 708,382: 18x19
		var (
			id   int
			x, y int
			w, h int
		)
		if !scanner(claim).scan(t, `#(\d+) @ (\d+),(\d+): (\d+)x(\d+)`, &id, &x, &y, &w, &h) {
			continue
		}
		for xx := x; xx < x+w; xx++ {
			for yy := y; yy < y+h; yy++ {
				fabric[xx][yy] = append(fabric[xx][yy], id)
			}
		}
	}

nextClaim:
	for _, claim := range in {
		// #1212 @ 708,382: 18x19
		var (
			id   int
			x, y int
			w, h int
		)
		if !scanner(claim).scan(t, `#(\d+) @ (\d+),(\d+): (\d+)x(\d+)`, &id, &x, &y, &w, &h) {
			continue
		}
		for xx := x; xx < x+w; xx++ {
			for yy := y; yy < y+h; yy++ {
				if len(fabric[xx][yy]) > 1 {
					continue nextClaim
				}
			}
		}
		return id
	}

	panic("none found")
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want int
	}{
		{"example1", []string{
			"#1 @ 1,3: 4x4",
			"#2 @ 3,1: 4x4",
			"#3 @ 5,5: 2x2",
		}, 3},
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
		t.Logf("part2: %v", part2(t, strings.Split(read(t, "input.txt"), "\n")))
	})
}
