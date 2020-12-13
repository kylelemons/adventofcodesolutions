// Copyright 2019 Kyle Lemons
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
	"github.com/kylelemons/adventofcodesolutions/advent/coords"
)

type Input struct {
	State [][]byte
}

func (input *Input) Key() string {
	out := new(strings.Builder)
	for _, line := range input.State {
		fmt.Fprintf(out, "%s\n", line)
	}
	return out.String()
}

func (input *Input) Biodiv() (res int) {
	power := 1
	for _, row := range input.State {
		for _, col := range row {
			if col == '#' {
				res += power
			}
			power *= 2
		}
	}
	return
}

func parseInput(t *testing.T, in string) *Input {
	input := &Input{
		State: advent.Split2D(in),
	}
	return input
}

func part1(t *testing.T, in string) (ret int) {
	input := parseInput(t, in)

	seen := make(map[string]int)

	for i := 0; i < 10000000; i++ {
		k := input.Key()
		if _, ok := seen[k]; ok {
			return input.Biodiv()
		}
		seen[k] = i
		// fmt.Println(k)

		next := advent.Make2D(len(input.State), len(input.State[0]))
		for r := range input.State {
			for c := range input.State[r] {
				cur := input.State[r][c]
				count := 0
				for _, d := range coords.Cardinals {
					nrc := coords.RC(r, c).Add(d)

					found := func() byte {
						defer func() { recover() }()
						if input.State[nrc.R()][nrc.C()] == '#' {
							return '#'
						}
						return '.'
					}()
					if found == '#' {
						count++
					}
				}
				if cur == '#' {
					if count == 1 {
						next[r][c] = '#'
					} else {
						next[r][c] = '.'
					}
				} else {
					if count == 1 || count == 2 {
						next[r][c] = '#'
					} else {
						next[r][c] = '.'
					}
				}
			}
		}
		input.State = next
	}
	return -1
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		// {"part1 example 0", "...", 0},
		{"part1 answer", advent.ReadFile(t, "input.txt"), 3186366},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part1(t, test.in), test.want; got != want {
				t.Errorf("part1(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

func part2(t *testing.T, in string, minutes int) (ret int) {
	space := [][][]byte{advent.Split2D(in)}

	for i := 0; i < minutes; i++ {
		prev := space

		// fmt.Println("---------------")
		// for li, l := range space {
		// 	fmt.Printf("Depth %d:\n", li-len(space)/2)
		// 	for _, r := range l {
		// 		fmt.Printf("  %s\n", r)
		// 	}
		// 	fmt.Println()
		// }
		// fmt.Println("---------------")

		space = make([][][]byte, 0, 2*len(space)+1)
		space = append(space, advent.Make2D(len(prev[0]), len(prev[0][0])))
		space = append(space, prev...)
		space = append(space, advent.Make2D(len(prev[0]), len(prev[0][0])))

		var next [][][]byte
		for range space {
			next = append(next, advent.Make2D(len(space[0]), len(space[0][0])))
		}

		// -1 is outer, +1 is inner
		var countFrom func(l, r, c int, looking coords.Vector) int
		countFrom = func(l, r, c int, looking coords.Vector) (cnt int) {
			// if l == 0 && r == 2 && c == 2 && looking == coords.West {
			// 	defer func() {
			// 		fmt.Println("countFrom", l, r, c, looking, cnt)
			// 	}()
			// }
			switch {
			case l < 0:
				return 0
			case l >= len(space):
				return 0
			case r < 0:
				return countFrom(l-1, 1, 2, coords.North)
			case c < 0:
				return countFrom(l-1, 2, 1, coords.West)
			case r >= len(space[l]):
				return countFrom(l-1, 3, 2, coords.South)
			case c >= len(space[l][r]):
				return countFrom(l-1, 2, 3, coords.East)
			case r == 2 && c == 2:
				switch looking {
				case coords.North:
					for i := 0; i < 5; i++ {
						cnt += countFrom(l+1, 4, i, looking)
					}
				case coords.East:
					for i := 0; i < 5; i++ {
						cnt += countFrom(l+1, i, 0, looking)
					}
				case coords.South:
					for i := 0; i < 5; i++ {
						cnt += countFrom(l+1, 0, i, looking)
					}
				case coords.West:
					for i := 0; i < 5; i++ {
						cnt += countFrom(l+1, i, 4, looking)
					}
				}
				return cnt
			case space[l][r][c] == '#':
				return 1
			default:
				return 0
			}
		}

		for l := range space {
			for r := range space[l] {
				for c := range space[l][r] {
					if r == 2 && c == 2 {
						next[l][r][c] = '?'
						continue
					}

					cur := space[l][r][c]
					count := 0

					for _, d := range coords.Cardinals {
						nrc := coords.RC(r, c).Add(d)
						count += countFrom(l, nrc.R(), nrc.C(), d)
						// if l == 2 && r == 0 {
						// 	fmt.Printf("(%d,%d)+%v : countFrom(%d, %d, %d, %d) = %v\n",
						// 		r, c,
						// 		d,
						// 		l-len(space)/2, nrc.R(), nrc.C(), d,
						// 		countFrom(l, nrc.R(), nrc.C(), d))
						// }
					}

					if cur == '#' {
						if count == 1 {
							next[l][r][c] = '#'
						} else {
							next[l][r][c] = '.'
						}
					} else {
						if count == 1 || count == 2 {
							next[l][r][c] = '#'
						} else {
							next[l][r][c] = '.'
						}
					}
				}
			}
		}
		space = next
	}

	var bugs int
	for li, l := range space {
		_ = li
		// fmt.Printf("Depth %d:\n", li-len(space)/2)
		for _, r := range l {
			// fmt.Printf("  %s\n", r)
			for _, c := range r {
				if c == '#' {
					bugs++
				}
			}
		}
		// fmt.Println()
	}
	return bugs
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name string
		in   string
		min  int
		want int
	}{
		{"part2 example 0", `....#
#..#.
#.?##
..#..
#....`, 10, 99},
		{"part2 answer", advent.ReadFile(t, "input.txt"), 200, 2031},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in, test.min), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}
