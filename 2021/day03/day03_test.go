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
	"fmt"
	"strconv"
	"testing"

	"github.com/kylelemons/adventofcodesolutions/advent"
)

type Input struct {
	Lines []string
	Count int
	Len   int
}

func parseInput(t *testing.T, in string) *Input {
	input := &Input{
		Lines: advent.Lines(in).All(t),
	}
	input.Count = len(input.Lines)
	input.Len = len(input.Lines[0])
	return input
}

func part1(t *testing.T, in string) (ret int) {
	input := parseInput(t, in)

	total := len(input.Lines)
	half := total / 2

	ones := make([]int, input.Count)

	gamma := make([]byte, input.Len)
	epsilon := make([]byte, input.Len)
	for i := 0; i < input.Len; i++ {
		for _, line := range input.Lines {
			if line[i] == '1' {
				ones[i]++
			}
		}
		if ones[i] > half {
			gamma[i] = '1'
		} else {
			gamma[i] = '0'
		}
		if ones[i] > half {
			epsilon[i] = '0'
		} else {
			epsilon[i] = '1'
		}
	}

	gammaRate := int(advent.Must(strconv.ParseInt(string(gamma), 2, 64))(t))
	epsilonRate := int(advent.Must(strconv.ParseInt(string(epsilon), 2, 64))(t))

	return gammaRate * epsilonRate
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part1 example 0", `00100
11110
10110
10111
10101
01111
00111
11100
10000
11001
00010
01010`, 198},
		{"part1 answer", advent.ReadFile(t, "input.txt"), 2724524},
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

	o2 := func(lines []string, bit int) byte {
		ones := 0
		zeros := 0
		for _, line := range lines {
			if line[bit] == '1' {
				ones++
			} else {
				zeros++
			}
		}
		if ones >= zeros {
			return '1'
		}
		return '0'
	}
	co2 := func(lines []string, bit int) byte {
		if o2(lines, bit) == '1' {
			return '0'
		}
		return '1'
	}

	filter := func(current []string, bit int, pred func(lines []string, bit int) byte) []string {
		out := make([]string, 0, len(current))
		want := pred(current, bit)
		fmt.Println(string(rune(want)), "at bit", bit)
		for _, line := range current {
			if line[bit] == want {
				out = append(out, line)
			}
		}
		return out
	}

	o2reading := input.Lines
	for i := range o2reading[0] {
		if len(o2reading) == 1 {
			break
		}
		o2reading = filter(o2reading, i, o2)
		fmt.Println("o2", o2reading)
	}
	o2rating := int(advent.Must(strconv.ParseInt(string(o2reading[0]), 2, 64))(t))

	co2reading := input.Lines
	for i := range co2reading[0] {
		if len(co2reading) == 1 {
			break
		}
		co2reading = filter(co2reading, i, co2)
		fmt.Println("co2", co2reading)
	}
	co2rating := int(advent.Must(strconv.ParseInt(string(co2reading[0]), 2, 64))(t))

	return o2rating * co2rating
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part2 example 0", `00100
11110
10110
10111
10101
01111
00111
11100
10000
11001
00010
01010`, 230},
		{"part2 answer", advent.ReadFile(t, "input.txt"), 2775870},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}
