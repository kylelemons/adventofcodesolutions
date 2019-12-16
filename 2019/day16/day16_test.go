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
	"log"
	"strconv"
	"testing"

	"github.com/kylelemons/adventofcodesolutions/advent"
)

type Input struct {
	Values  []int8
	Scratch []int8
	Offset  int
}

func parseInput(t *testing.T, in string, N int) *Input {
	input := &Input{
		Values:  make([]int8, len(in)*N),
		Scratch: make([]int8, len(in)*N),
	}
	for i := range input.Values {
		input.Values[i] = int8(in[i%len(in)] - '0')
	}
	input.Offset, _ = strconv.Atoi(in[:7])
	return input
}

var base = []int{0, 1, 0, -1}

func pat(i, j int) int {
	offset := (j + 1) / (i + 1)
	return base[(offset)%len(base)]
}

func (input *Input) FastFFT(t *testing.T) {
	current, next := input.Values, input.Scratch

	for i := range next {
		var sum int

		start := i
		repeats := i + 1
		sign := 1

		for start < len(next) {
			end := start + repeats
			if end > len(next) {
				end = len(next)
			}
			for _, n := range current[start:end] {
				sum += sign * int(n)
			}
			start += 2 * repeats
			sign *= -1
		}

		v := sum % 10
		if v < 0 {
			v = -v
		}
		next[i] = int8(v)
	}
	input.Values, input.Scratch = input.Scratch, input.Values
}

func finalize(v []int8) string {
	out := make([]byte, len(v))
	for i := range out {
		out[i] = '0' + byte(v[i])
	}
	return string(out)
}

func part1(t *testing.T, in string) (ret string) {
	input := parseInput(t, in, 1)

	for i := 0; i < 100; i++ {
		input.FastFFT(t)
	}

	return finalize(input.Values[0:][:8])
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"part1 example 0", "80871224585914546619083218645595", "24176176"},
		{"part1 example 1", "80871224585914546619083218645595", "24176176"},
		{"part1 example 2", "19617804207202209144916044189917", "73745418"},
		{"part1 answer", advent.ReadFile(t, "input.txt"), "23135243"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part1(t, test.in), test.want; fmt.Sprint(got) != want {
				t.Errorf("part1(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

func (input *Input) CheatingFFT(t *testing.T) {
	if input.Offset < len(input.Values)/2 {
		t.Fatalf("Can't cheat, offset too early")
	}

	var sum int
	for i := len(input.Values) - 1; i >= input.Offset; i-- {
		sum += int(input.Values[i])

		v := sum % 10
		if v < 0 {
			v = -v
		}
		input.Values[i] = int8(v)
	}
}

func part2(t *testing.T, in string) (ret string) {
	input := parseInput(t, in, 10000)
	log.Printf("Input size: %d", len(input.Values))
	log.Printf("Offset:     %d", input.Offset)
	for i := 0; i < 100; i++ {
		input.CheatingFFT(t)
	}

	return finalize(input.Values[input.Offset:][:8])
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"part2 example 0", "03036732577212944063491565474664", "84462026"},
		{"part2 example 1", "02935109699940807407585447034323", "78725270"},
		{"part2 example 2", "03081770884921959731165446850517", "53553731"},
		{"part2 answer", advent.ReadFile(t, "input.txt"), "21130597"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; fmt.Sprint(got) != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}
