// Copyright 2018 Kyle Lemons
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
package acoday

import (
	"fmt"
	"strings"
	"testing"

	"github.com/kylelemons/adventofcodesolutions/advent"
)

type Point struct{ X, Y int }

func (p Point) Add(dx, dy int) Point {
	return Point{p.X + dx, p.Y + dy}
}

type Input struct {
	At     map[Point]byte
	X, Y   advent.RangeTracker
	Spring Point
}

func (in *Input) Map() string {
	out := new(strings.Builder)
	for y := in.Y.Min; y <= in.Y.Max; y++ {
		for x := in.X.Min; x <= in.X.Max; x++ {
			// fmt.Fprintf(out, " [%v,%v]", x, y)
			b, ok := in.At[Point{x, y}]
			if !ok || b == 0 {
				b = '.'
			}
			fmt.Fprintf(out, "%c", b)
		}
		fmt.Fprintln(out)
	}
	return out.String()
}

func parseInput(t *testing.T, in string) *Input {
	input := &Input{
		At: make(map[Point]byte),
	}
	advent.Lines(in).Extract(t, `([xy])=(-?\d+), [xy]=(-?\d+)..(-?\d+)`, func(dimA string, valA, minB, maxB int) {
		var x0, x1, y0, y1 int
		switch dimA {
		case "x":
			x0, x1, y0, y1 = valA, valA, minB, maxB
		case "y":
			y0, y1, x0, x1 = valA, valA, minB, maxB
		}

		for x := x0; x <= x1; x++ {
			input.X.Track(x)
			for y := y0; y <= y1; y++ {
				input.Y.Track(y)
				input.At[Point{x, y}] = '#'
			}
		}
	})
	input.Spring = Point{500, input.Y.Min}
	return input
}

type FillResult struct {
	reachable int
	retained  int
}

func (in *Input) Fill(t *testing.T, from Point) (ret FillResult) {
	in.X.Track(from.X) // in case we go wide

	// Mark visited, just in case.
	in.At[from] = '|'
	ret.reachable++

	// Water flows downward.
	floor := from
flowDown:
	for {
		floor = floor.Add(0, 1)
		if floor.Y > in.Y.Max {
			ret.reachable += 0 // Don't count this one, we're out of bounds
			return
		}
		switch ch := in.At[floor]; ch {
		case 0:
			// Keep flowing through empty space.
			in.At[floor] = '|'
			ret.reachable++
			continue
		case '#', '~':
			// Reached a floor that we can fill.
			break flowDown
		case '|':
			// Reached already-flowing water.
			return
		default:
			t.Fatalf("Unrecognized value %q flowing down:", ch)
		}
	}

	// The water is on the level above the floor.
	water := floor.Add(0, -1) // this tile has been counted as reachable

	scan := func(dx int) (p Point, wall bool) {
		p = water.Add(dx, 0)
		for {
			// See if we hit a wall that will hold the water in.
			switch in.At[p] {
			case '#':
				// Water stops at a wall to the p.
				wall = true
				return
			}

			// Check if the water at p is supported by something.
			below := p.Add(0, 1)
			switch in.At[below] {
			case '#', '~':
				// Water can flow through this block.
				in.At[p] = '|'
				ret.reachable++
				p = p.Add(dx, 0)
				continue
			}

			// Water is unsupported, and flows down here.
			in.At[p] = '|'
			res := in.Fill(t, p)
			ret.reachable += res.reachable
			ret.retained += res.retained
			switch in.At[below] {
			case '#', '~':
				// It's filled in now!
				p = p.Add(dx, 0)
				continue
			}
			return
		}
	}

	// Figure out how far out we can extend the water.
waterLevel:
	left, leftWall := scan(-1)
	right, rightWall := scan(1)

	// Check if we filled anything in.
	if leftWall && rightWall {
		for w := left.Add(1, 0); w != right; w = w.Add(1, 0) {
			in.At[w] = '~'
			ret.retained++
		}
		water = water.Add(0, -1)
		if water.Y > from.Y {
			goto waterLevel
		}
	}
	return
}

func part1(t *testing.T, in string) int {
	input := parseInput(t, in)
	// t.Logf("Decoded input:\n%s", input.Map())
	ret := input.Fill(t, input.Spring)
	t.Logf("Filled input:\n%s", input.Map())
	return ret.reachable
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part1 example 0", strings.Trim(`
x=495, y=2..7
y=7, x=495..501
x=501, y=3..7
x=498, y=2..4
x=506, y=1..2
x=498, y=10..13
x=504, y=10..13
y=13, x=498..504
`, "\n"), 57},
		{"part1 answer", advent.ReadFile(t, "input.txt"), 32552},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part1(t, test.in), test.want; got != want {
				t.Errorf("part1(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

func part2(t *testing.T, in string) int {
	input := parseInput(t, in)
	ret := input.Fill(t, input.Spring)
	return ret.retained
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part2 example 0", strings.Trim(`
x=495, y=2..7
y=7, x=495..501
x=501, y=3..7
x=498, y=2..4
x=506, y=1..2
x=498, y=10..13
x=504, y=10..13
y=13, x=498..504
`, "\n"), 29},
		{"part2 answer", advent.ReadFile(t, "input.txt"), 26405},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}
