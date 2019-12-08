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
	"strconv"
	"strings"
	"testing"
)

func part1(t *testing.T, start string, after int) (ret string) {
	recipes := []byte(start)
	current := []int{0, 1}

	for len(recipes) < after+10 {
		var sum byte
		for _, i := range current {
			sum += recipes[i] - '0'
		}
		recipes = append(recipes, strconv.Itoa(int(sum))...)

		for elf, at := range current {
			current[elf] = (at + 1 + int(recipes[at]-'0')) % len(recipes)
		}
	}
	return string(recipes)[after:][:10]
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name  string
		start string
		after int
		want  string
	}{
		{"part1 example 0", "37", 9, "5158916779"},
		{"part1 example 1", "37", 5, "0124515891"},
		{"part1 example 2", "37", 18, "9251071085"},
		{"part1 example 3", "37", 2018, "5941429882"},
		{"part1 answer", "37", 513401, "5371393113"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part1(t, test.start, test.after), test.want; got != want {
				t.Errorf("part1(%#v)\n = %#v, want %#v", test.start, got, want)
			}
		})
	}
}

func part2(t *testing.T, start string, in string) (ret int) {
	recipes := []byte(start)
	current := []int{0, 1}

	for {
		var sum byte
		for _, i := range current {
			sum += recipes[i] - '0'
		}
		recipes = append(recipes, strconv.Itoa(int(sum))...)

		for elf, at := range current {
			current[elf] = (at + 1 + int(recipes[at]-'0')) % len(recipes)
		}
		if len(recipes) < len(in)+2 {
			// TODO this has a small blind spot
			continue
		}

		first := len(recipes) - len(in) - 2
		if i := strings.Index(string(recipes[first:]), in); i >= 0 {
			return first + i
		}
	}
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name  string
		start string
		in    string
		want  int
	}{
		{"part2 example 0", "37", "51589", 9},
		{"part2 example 1", "37", "01245", 5},
		{"part2 example 2", "37", "92510", 18},
		{"part2 example 3", "37", "59414", 2018},
		{"part2 answer", "37", "513401", 20286858},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.start, test.in), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.start, got, want)
			}
		})
	}
}
