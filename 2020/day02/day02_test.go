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
)

type Input struct {
	Passwords []Password
}

func parseInput(t *testing.T, in string) *Input {
	input := &Input{
		// ...
	}
	advent.Lines(in).Extract(t, `(\d+)-(\d+) (.): (.*)`, func(min, max int, letter, password string) {
		input.Passwords = append(input.Passwords, Password{
			Min:      min,
			Max:      max,
			Letter:   letter[0],
			Password: password,
		})
	})
	return input
}

type Password struct {
	Min, Max int
	Letter   byte
	Password string
}

func (p Password) IsValidOldCompany() bool {
	var count int
	for _, c := range []byte(p.Password) {
		if c == p.Letter {
			count++
		}
	}
	return count >= p.Min && count <= p.Max
}

func part1(t *testing.T, in string) (ret int) {
	input := parseInput(t, in)

	for _, p := range input.Passwords {
		if p.IsValidOldCompany() {
			ret++
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
		{"part1 example 0", `1-3 a: abcde
1-3 b: cdefg
2-9 c: ccccccccc`, 2},
		{"part1 answer", advent.ReadFile(t, "input.txt"), 460},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part1(t, test.in), test.want; got != want {
				t.Errorf("part1(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

func (p Password) IsValidNow() bool {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("%+v: %v", p, r)
		}
	}()

	var count int
	for _, i := range []int{p.Min, p.Max} {
		if p.Password[i-1] == p.Letter {
			count++
		}
	}
	return count == 1
}

func part2(t *testing.T, in string) (ret int) {
	input := parseInput(t, in)

	for _, p := range input.Passwords {
		// log.Printf("%d-%d %q %q: %v", p.Min, p.Max, p.Letter, p.Password, p.IsValidNow())
		if p.IsValidNow() {
			ret++
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
		{"part2 example 0", `1-3 a: abcde
1-3 b: cdefg
2-9 c: ccccccccc`, 1},
		{"part2 answer", advent.ReadFile(t, "input.txt"), 251},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}
