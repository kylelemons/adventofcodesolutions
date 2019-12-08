package aocday

import (
	"fmt"
	"strings"
	"testing"
)

func part1(t *testing.T, in []string) int {
	var fabric [1000][1000]int

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
		fmt.Println(id, x, y, w, h)
		for xx := x; xx < x+w; xx++ {
			for yy := y; yy < y+h; yy++ {
				fabric[xx][yy]++
			}
		}
	}

	var dups int
	for i := range fabric {
		for j := range fabric[i] {
			if fabric[i][j] > 1 {
				dups++
			}
		}
	}

	return dups
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want int
	}{
		{"example1", []string{
			"#1 @ 1,3: 4x4",
			"#2 @ 3,1: 4x4",
			"#3 @ 5,5: 2x2",
		}, 4},
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
		t.Logf("part1: %v", part1(t, strings.Split(read(t, "input.txt"), "\n")))
	})
}
