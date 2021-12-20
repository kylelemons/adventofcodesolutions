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
	"testing"

	"github.com/kylelemons/adventofcodesolutions/advent"
)

type Input struct {
	Draws     []int
	Boards    [][][]int
	Locations map[int][]loc
}

type loc struct {
	board int
	id    int // 0-4 row, 5+ col
}

func parseInput(t *testing.T, in string) *Input {
	input := &Input{
		// ...
	}
	records := advent.Records(in).All(t)

	advent.Split(records[0], ',').Scan(t, func(i int) {
		input.Draws = append(input.Draws, i)
	})
	for _, rec := range records[1:] {
		var board [][]int
		advent.Lines(rec).Scan(t, func(a, b, c, d, e int) {
			board = append(board, []int{a, b, c, d, e})
		})
		input.Boards = append(input.Boards, board)
	}

	nums := make(map[int][]loc)
	for i, b := range input.Boards {
		for r := range b {
			for c := range b[r] {
				num := b[r][c]
				nums[num] = append(nums[num], loc{
					board: i,
					id:    r,
				})
				nums[num] = append(nums[num], loc{
					board: i,
					id:    5 + c,
				})
			}
		}
	}
	input.Locations = nums
	return input
}

func part1(t *testing.T, in string) (ret int) {
	input := parseInput(t, in)

	counts := map[loc]int{}
	marked := map[int]bool{}
	for _, draw := range input.Draws {
		marked[draw] = true
		for _, l := range input.Locations[draw] {
			counts[l]++
			if counts[l] == 5 {
				var unmarked int
				for _, row := range input.Boards[l.board] {
					for _, num := range row {
						if marked[num] {
							continue
						}
						unmarked += num
					}
				}
				return unmarked * draw
			}
		}
	}

	return -1
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part1 example 0", `7,4,9,5,11,17,23,2,0,14,21,24,10,16,13,6,15,25,12,22,18,20,8,19,3,26,1

22 13 17 11  0
8  2 23  4 24
21  9 14 16  7
6 10  3 18  5
1 12 20 15 19

3 15  0  2 22
9 18 13 17  5
19  8  7 25 23
20 11 10 24  4
14 21 16 12  6

14 21 17 24  4
10 16 15  9 19
18  8 23 26 20
22 11 13  6  5
2  0 12  3  7`, 4512},
		{"part1 answer", advent.ReadFile(t, "input.txt"), 65325},
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

	counts := map[loc]int{}
	marked := map[int]bool{}
	alive := map[int]bool{}
	for i := range input.Boards {
		alive[i] = true
	}
	for _, draw := range input.Draws {
		marked[draw] = true
		for _, l := range input.Locations[draw] {
			if !alive[l.board] {
				continue
			}
			counts[l]++
			if counts[l] == 5 {
				delete(alive, l.board)
				if len(alive) > 0 {
					continue
				}
				var unmarked int
				for _, row := range input.Boards[l.board] {
					for _, num := range row {
						if marked[num] {
							continue
						}
						unmarked += num
					}
				}
				return unmarked * draw
			}
		}
	}

	return -1
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part2 example 0", `7,4,9,5,11,17,23,2,0,14,21,24,10,16,13,6,15,25,12,22,18,20,8,19,3,26,1

22 13 17 11  0
8  2 23  4 24
21  9 14 16  7
6 10  3 18  5
1 12 20 15 19

3 15  0  2 22
9 18 13 17  5
19  8  7 25 23
20 11 10 24  4
14 21 16 12  6

14 21 17 24  4
10 16 15  9 19
18  8 23 26 20
22 11 13  6  5
2  0 12  3  7`, 1924},
		{"part2 answer", advent.ReadFile(t, "input.txt"), 4624},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}
