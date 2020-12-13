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
	"strings"
	"testing"

	"github.com/kylelemons/adventofcodesolutions/advent"
)

func part1(t *testing.T, in string) (total int) {
	for _, module := range strings.Fields(in) {
		var mass int
		advent.Scanner(module).Scan(t, &mass)
		total += (mass / 3) - 2
	}
	return
}

func fuel(mass int) int {
	f := (mass / 3) - 2
	if f <= 0 {
		return 0
	}
	return f + fuel(f)
}

func part2(t *testing.T, in string) (total int) {
	advent.Lines(in).Scan(t, func(mass int) {
		total += fuel(mass)
	})
	return
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part1 example 1", "12", 2},
		{"part1 example 2", "14", 2},
		{"part1 example 3", "1969", 654},
		{"part1 example 4", "100756", 33583},
		{"combined", "12 14", 4},
		{"part1 answer", advent.ReadFile(t, "input.txt"), 3256794},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part1(t, test.in), test.want; got != want {
				t.Errorf("part1(%#v) = %#v, want %#v", test.in, got, want)
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
		{"part2 example 1", "12", 2},
		{"part2 example 2", "1969", 966},
		{"part2 example 4", "100756", 50346},
		{"part2 answer", advent.ReadFile(t, "input.txt"), 4882337},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got != want {
				t.Errorf("part2(%#v) = %#v, want %#v", test.in, got, want)
			}
		})
	}
}
