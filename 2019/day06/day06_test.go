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
	"strings"
	"testing"

	"github.com/kylelemons/adventofcodesolutions/2019/advent"
)

type orbit struct {
	inner, outer string
}

func part1(t *testing.T, in string) (ret int) {
	orbits := make(map[string][]string)
	advent.Lines(in).Extract(t, `([^)]+)\)(.+)`, func(inner, outer string) {
		orbits[inner] = append(orbits[inner], outer)
	})

	var branches func(from string, depth int) int
	branches = func(from string, depth int) int {
		var total int
		for _, next := range orbits[from] {
			// fmt.Println(from, next)
			total += depth + 1 + branches(next, depth+1)
		}
		return total
	}

	return branches("COM", 0)
}

func part2(t *testing.T, in string) (ret int) {
	orbits := make(map[string][]string)
	orbiting := make(map[string]string)
	advent.Lines(in).Extract(t, `([^)]+)\)(.+)`, func(inner, outer string) {
		orbits[inner] = append(orbits[inner], outer)
		orbiting[outer] = inner
	})

	path := make(map[string]bool)

	// Count the distance to the center.
	for planet := "YOU"; planet != ""; planet = orbiting[planet] {
		path[planet] = true
	}

	// Count the distance to the first common planet.
	toCommon := 0
	for planet := "SAN"; planet != ""; planet = orbiting[planet] {
		if path[planet] {
			delete(path, planet) // common path, subtract
		} else {
			toCommon++ // haven't reached common yet
		}
	}

	return toCommon + len(path) - 2
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part1 example 0", strings.TrimSpace(`
COM)B
B)C
C)D
D)E
E)F
B)G
G)H
D)I
E)J
J)K
K)L
`), 42},
		{"part1 answer", advent.ReadFile(t, "input.txt"), 249308},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part1(t, test.in), test.want; got != want {
				t.Errorf("part1(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part2 example 0", strings.TrimSpace(`
COM)B
B)C
C)D
D)E
E)F
B)G
G)H
D)I
E)J
J)K
K)L
K)YOU
I)SAN`), 4},
		{"part2 answer", advent.ReadFile(t, "input.txt"), 349},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}
