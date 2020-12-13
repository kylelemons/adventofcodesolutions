// Copyright 2020 Kyle Lemons
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
	"log"
	"testing"

	"github.com/kylelemons/adventofcodesolutions/advent"
	"github.com/kylelemons/adventofcodesolutions/advent/coords"
)

type Input struct {
	Seats [][]byte
}

func parseInput(t *testing.T, in string) *Input {
	input := &Input{
		Seats: advent.Split2D(in),
	}
	return input
}

func part1(t *testing.T, in string) (ret int) {
	input := parseInput(t, in)

	step := func() (changes int) {
		next := advent.Make2D(len(input.Seats), len(input.Seats[0]))
		defer func() { input.Seats = next }()
		for r := range input.Seats {
			for c := range input.Seats[r] {
				coord := coords.RC(r, c)
				s := coord.In2D(input.Seats)
				occupied := 0
				for _, dc := range coords.Compass {
					b, ok := coord.Add(dc).InBounds2D(input.Seats)
					if !ok {
						continue
					}
					if b == '#' {
						occupied++
					}
				}
				// log.Printf("%v %q %d", coord, s, occupied)
				switch {
				case s == 'L' && occupied == 0:
					next[r][c] = '#'
					changes++
				case s == '#' && occupied >= 4:
					next[r][c] = 'L'
					changes++
				default:
					next[r][c] = input.Seats[r][c]
				}
			}
		}
		log.Printf("%d changes", changes)
		return changes
	}

	for step() > 0 {
		// log.Printf("----")
		// for _, r := range input.Seats {
		// 	log.Printf("> %s", r)
		// }
	}

	for r := range input.Seats {
		for c := range input.Seats[r] {
			if input.Seats[r][c] == '#' {
				ret++
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
		{"part1 example 0", `L.LL.LL.LL
LLLLLLL.LL
L.L.L..L..
LLLL.LL.LL
L.LL.LL.LL
L.LLLLL.LL
..L.L.....
LLLLLLLLLL
L.LLLLLL.L
L.LLLLL.LL`, 37},
		{"part1 answer", advent.ReadFile(t, "input.txt"), 2418},
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

	step := func() (changes int) {
		next := advent.Make2D(len(input.Seats), len(input.Seats[0]))
		defer func() { input.Seats = next }()
		for r := range input.Seats {
			for c := range input.Seats[r] {
				coord := coords.RC(r, c)
				s := coord.In2D(input.Seats)
				occupied := 0
			nextDir:
				for _, dc := range coords.Compass {
					consider := coord
					for {
						consider = consider.Add(dc)
						b, ok := consider.InBounds2D(input.Seats)
						if !ok {
							continue nextDir
						}
						switch b {
						case '#':
							occupied++
							continue nextDir
						case 'L':
							continue nextDir
						}
					}
				}
				// log.Printf("%v %q %d", coord, s, occupied)
				switch {
				case s == 'L' && occupied == 0:
					next[r][c] = '#'
					changes++
				case s == '#' && occupied >= 5:
					next[r][c] = 'L'
					changes++
				default:
					next[r][c] = input.Seats[r][c]
				}
			}
		}
		log.Printf("%d changes", changes)
		return changes
	}

	for step() > 0 {
		// log.Printf("----")
		// for _, r := range input.Seats {
		// 	log.Printf("> %s", r)
		// }
	}

	for r := range input.Seats {
		for c := range input.Seats[r] {
			if input.Seats[r][c] == '#' {
				ret++
			}
		}
	}
	return
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part2 example 0", `L.LL.LL.LL
LLLLLLL.LL
L.L.L..L..
LLLL.LL.LL
L.LL.LL.LL
L.LLLLL.LL
..L.L.....
LLLLLLLLLL
L.LLLLLL.L
L.LLLLL.LL`, 26},
		{"part2 answer", advent.ReadFile(t, "input.txt"), 2144},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}
