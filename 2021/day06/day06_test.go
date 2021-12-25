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
	"fmt"
	"math/big"
	"testing"

	"github.com/kylelemons/adventofcodesolutions/advent"
)

type Input struct {
	FishAtCounter [9]int
}

func parseInput(t *testing.T, in string) *Input {
	input := &Input{
		// ...
	}
	advent.Split(in, ',').Scan(t, func(counter int) {
		input.FishAtCounter[counter]++
	})
	return input
}

func part1(t *testing.T, in string) (ret int) {
	input := parseInput(t, in)

	const SimulateDays = 80

	for i := 0; i < SimulateDays; i++ {
		fmt.Println(input.FishAtCounter)
		var next [9]int
		for counter, fish := range input.FishAtCounter {
			if counter == 0 {
				next[6] += fish // time to spawn a new fish
				next[8] += fish // new fish that spawned
				continue
			}
			next[counter-1] += fish
		}
		input.FishAtCounter = next
	}

	for _, fish := range input.FishAtCounter {
		ret += fish
	}
	return
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part1 example 0", "3,4,3,1,2", 5934},
		{"part1 answer", advent.ReadFile(t, "input.txt"), 355386},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part1(t, test.in), test.want; got != want {
				t.Errorf("part1(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

func part2(t *testing.T, in string) (ret *big.Int) {
	input := parseInput(t, in)

	const SimulateDays = 256

	for i := 0; i < SimulateDays; i++ {
		fmt.Println(input.FishAtCounter)
		var next [9]int
		for counter, fish := range input.FishAtCounter {
			if counter == 0 {
				next[6] += fish // time to spawn a new fish
				next[8] += fish // new fish that spawned
				continue
			}
			next[counter-1] += fish
		}
		input.FishAtCounter = next
	}

	total := new(big.Int)
	for _, fish := range input.FishAtCounter {
		total.Add(total, big.NewInt(int64(fish)))
	}
	return total
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want *big.Int
	}{
		{"part2 example 0", "3,4,3,1,2", big.NewInt(26984457539)},
		{"part2 answer", advent.ReadFile(t, "input.txt"), big.NewInt(1613415325809)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got.Cmp(want) != 0 {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}
