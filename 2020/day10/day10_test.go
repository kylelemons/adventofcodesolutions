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
	"sort"
	"testing"

	"github.com/kylelemons/adventofcodesolutions/advent"
)

type Input struct {
	joltages []int
}

func parseInput(t *testing.T, in string) *Input {
	input := &Input{
		// ...
	}
	advent.Lines(in).Scan(t, func(j int) {
		input.joltages = append(input.joltages, j)
	})
	return input
}

func part1(t *testing.T, in string) (ret int) {
	input := parseInput(t, in)

	input.joltages = append(input.joltages, 0) // arm rest
	sort.Ints(input.joltages)
	input.joltages = append(input.joltages, input.joltages[len(input.joltages)-1]+3) // device

	var ones, threes int
	for i := range input.joltages[1:] {
		delta := input.joltages[i+1] - input.joltages[i]
		switch delta {
		case 1:
			ones++
		case 3:
			threes++
		}
	}
	return ones * threes
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part1 example 0", `16
10
15
5
1
11
7
19
6
12
4`, 7 * 5},
		{"part1 example 1", `28
33
18
42
31
14
46
20
48
47
24
23
49
45
19
38
39
11
1
32
25
35
8
17
7
9
4
2
34
10
3`, 22 * 10},
		{"part1 answer", advent.ReadFile(t, "input.txt"), 2310},
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

	input.joltages = append(input.joltages, 0) // arm rest
	sort.Ints(input.joltages)
	input.joltages = append(input.joltages, input.joltages[len(input.joltages)-1]+3) // device

	var arrange func(int) int

	// Naive recursive implementation
	arrange = func(start int) (count int) {
		// Base case: we're at the end!
		if start == len(input.joltages)-1 {
			return 1
		}

		// Recursive cases that move closer to the conclusion
		joltage := input.joltages[start]
		max := joltage + 3

		for skip, nextJoltage := range input.joltages[start+1:] {
			if nextJoltage > max {
				break
			}
			count += arrange(start + 1 + skip)
		}
		return
	}

	// Memoize! (poof, dynamic programming)
	memo := make(map[int]int) // memo[start] = count
	real := arrange
	arrange = func(start int) int {
		if count, ok := memo[start]; ok {
			return count
		}
		count := real(start)
		memo[start] = count
		return count
	}

	return arrange(0)
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part2 example 0", `16
10
15
5
1
11
7
19
6
12
4`, 8},
		{"part2 example 1", `28
33
18
42
31
14
46
20
48
47
24
23
49
45
19
38
39
11
1
32
25
35
8
17
7
9
4
2
34
10
3`, 19208},
		{"part2 answer", advent.ReadFile(t, "input.txt"), 64793042714624},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}
