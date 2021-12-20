// Copyright 2021 Kyle Lemons
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package acoday is the entrypoint for this AoC solution.
package aocday

import (
	"testing"

	"github.com/kylelemons/adventofcodesolutions/advent"
	"github.com/kylelemons/adventofcodesolutions/advent/coords"
)

type Input struct {
	Lines []Line
	X, Y  advent.RangeTracker
}

type Line struct {
	From, To coords.Coord
}

func parseInput(t *testing.T, in string) *Input {
	input := &Input{
		// ...
	}
	advent.Lines(in).Extract(t, `(\d+),(\d+) -> (\d+),(\d+)`, func(x1, y1 int, x2, y2 int) {
		input.Lines = append(input.Lines, Line{
			From: coords.XY(x1, y1),
			To:   coords.XY(x2, y2),
		})
		input.X.TrackAll(x1, x2)
		input.Y.TrackAll(y1, y2)
	})
	return input
}

func part1(t *testing.T, in string) (ret int) {
	input := parseInput(t, in)

	visited := make(map[coords.Coord]int)
	for _, line := range input.Lines {
		if line.From.X() != line.To.X() && line.From.Y() != line.To.Y() {
			continue
		}

		x0, x1 := line.From.X(), line.To.X()
		if x1 < x0 {
			x0, x1 = x1, x0
		}
		y0, y1 := line.From.Y(), line.To.Y()
		if y1 < y0 {
			y0, y1 = y1, y0
		}
		for x := x0; x <= x1; x++ {
			for y := y0; y <= y1; y++ {
				visited[coords.XY(x, y)]++
				if visited[coords.XY(x, y)] == 2 {
					ret++
				}
			}
		}
	}

	return
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part1 example 0", `0,9 -> 5,9
8,0 -> 0,8
9,4 -> 3,4
2,2 -> 2,1
7,0 -> 7,4
6,4 -> 2,0
0,9 -> 2,9
3,4 -> 1,4
0,0 -> 8,8
5,5 -> 8,2`, 5},
		{"part1 answer", advent.ReadFile(t, "input.txt"), 8622},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part1(t, test.in), test.want; got != want {
				t.Errorf("part1(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

func part2(t *testing.T, in string) (ret int) {
	input := parseInput(t, in)

	visited := make(map[coords.Coord]int)
	for _, line := range input.Lines {
		// fmt.Println(line)

		x0, x1, dx := line.From.X(), line.To.X(), 0
		if x0 < x1 {
			dx = 1
		} else if x0 > x1 {
			dx = -1
		}

		y0, y1, dy := line.From.Y(), line.To.Y(), 0
		if y0 < y1 {
			dy = 1
		} else if y0 > y1 {
			dy = -1
		}

		for x, y := x0, y0; ; x, y = x+dx, y+dy {
			visited[coords.XY(x, y)]++
			if visited[coords.XY(x, y)] == 2 {
				ret++
			}
			if x == x1 && y == y1 {
				break
			}
		}

		// for y := 0; y <= input.Y.Max; y++ {
		// 	for x := 0; x <= input.X.Max; x++ {
		// 		c := visited[coords.XY(x, y)]
		// 		if c == 0 {
		// 			fmt.Print(".")
		// 			continue
		// 		}
		// 		fmt.Print(c)
		// 	}
		// 	fmt.Println()
		// }
	}
	return
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part2 example 0", `0,9 -> 5,9
8,0 -> 0,8
9,4 -> 3,4
2,2 -> 2,1
7,0 -> 7,4
6,4 -> 2,0
0,9 -> 2,9
3,4 -> 1,4
0,0 -> 8,8
5,5 -> 8,2`, 12},
		{"part2 answer", advent.ReadFile(t, "input.txt"), 22037},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}
