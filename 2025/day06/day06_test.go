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
	"strconv"
	"strings"
	"testing"

	"github.com/kylelemons/adventofcodesolutions/advent"
)

type Input struct {
	Lines [][]int
	Ops   []string
}

func parseInput(t *testing.T, in string) *Input {
	input := &Input{
		// ...
	}
	lines := strings.Split(strings.TrimSpace(in), "\n")
	lines, last := lines[:len(lines)-1], lines[len(lines)-1]

	input.Ops = strings.Fields(last)
	for _, line := range lines {
		parsed := make([]int, len(input.Ops))
		for i, num := range strings.Fields(line) {
			parsed[i], _ = strconv.Atoi(num)
		}
		input.Lines = append(input.Lines, parsed)
	}
	return input
}

func part1(t *testing.T, in string) (ret int) {
	input := parseInput(t, in)

	for i, op := range input.Ops {
		var f func(int, int) int
		var v0 int
		switch op {
		case "*":
			f, v0 = func(a, b int) int { return a * b }, 1
		case "+":
			f, v0 = func(a, b int) int { return a + b }, 0
		}

		v := v0
		for _, line := range input.Lines {
			v = f(v, line[i])
		}
		ret += v
	}

	return
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part1 example 0", "123 328  51 64 \n 45 64  387 23 \n  6 98  215 314\n*   +   *   +  ", 4277556},
		{"part1 answer", advent.ReadFile(t, "input.txt"), 6295830249262},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part1(t, test.in), test.want; got != want {
				t.Errorf("part1(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

func parseInputPart2(t *testing.T, in string) *Input {
	input := &Input{
		// ...
	}
	lines := strings.Split(in, "\n")
	lines, last := lines[:len(lines)-1], lines[len(lines)-1]

	input.Ops = strings.Fields(last)

	var currentOpNums []int
	for col := range len(lines[0]) {
		currentStr := ""
		for row := range len(lines) {
			currentStr += string(rune(lines[row][col]))
			// fmt.Printf("%d, %d, %q | %s %q\n", row, col, lines[row][col], "accumulating", currentStr)
		}
		currentStr = strings.TrimSpace(currentStr)
		if currentStr != "" {
			n, _ := strconv.Atoi(currentStr)
			currentOpNums = append(currentOpNums, n)
			// fmt.Println("appending", currentOpNums)
		}
		if currentStr == "" || col == len(lines[0])-1 {
			input.Lines = append(input.Lines, currentOpNums)
			// fmt.Printf("accumulated %v\n", currentOpNums)
			currentOpNums = nil
		}
	}

	return input
}

func part2(t *testing.T, in string) (ret int) {
	input := parseInputPart2(t, in)

	if got, want := len(input.Lines), len(input.Ops); got != want {
		t.Fatalf("got %d lines, want %d", got, want)
	}

	for i, op := range input.Ops {
		var f func(int, int) int
		var v0 int
		switch op {
		case "*":
			f, v0 = func(a, b int) int { return a * b }, 1
		case "+":
			f, v0 = func(a, b int) int { return a + b }, 0
		}

		v := v0
		for _, vn := range input.Lines[i] {
			v = f(v, vn)
		}
		ret += v
	}

	return
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part2 example 0", "123 328  51 64 \n 45 64  387 23 \n  6 98  215 314\n*   +   *   +  ", 3263827},
		{"part2 answer", advent.ReadFile(t, "input.txt"), 9194682052782},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}
