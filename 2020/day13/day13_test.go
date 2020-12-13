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
package aocday

import (
	"testing"

	"github.com/kylelemons/adventofcodesolutions/advent"
)

type idOffset struct {
	id     int
	offset int
}

type Input struct {
	t0        int
	ids       []int
	idOffsets []idOffset
}

func parseInput(t *testing.T, in string) *Input {
	input := &Input{
		// ...
	}
	lines := advent.Lines(in).All(t)

	advent.Scanner(lines[0]).Scan(t, &input.t0)
	advent.Split(lines[1], ',').Each(func(i int, token advent.Scanner) {
		if token == "x" {
			return
		}
		var id int
		token.Scan(t, &id)
		input.ids = append(input.ids, id)
		input.idOffsets = append(input.idOffsets, idOffset{
			id:     id,
			offset: i,
		})
	})
	return input
}

func part1(t *testing.T, in string) (ret int) {
	input := parseInput(t, in)

	for i := input.t0; ; i++ {
		for _, id := range input.ids {
			if i%id == 0 {
				return (i - input.t0) * id
			}
		}
	}
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part1 example 0", `939
7,13,x,x,59,x,31,19`, 295},
		{"part1 answer", advent.ReadFile(t, "input.txt"), 259},
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

	tt := 0

	for {
		// log.Printf("t=%d", tt)
		seq := 0
		d := 1
		for _, ido := range input.idOffsets {
			if (tt+ido.offset)%ido.id == 0 {
				seq++
				d *= ido.id
				continue
			}
			break
		}
		// if seq > 1 {
		// 	log.Printf("At t=%d, %d in sequence (delta: %d)", tt, seq, d)
		// 	for _, ido := range input.idOffsets[:seq] {
		// 		log.Printf(" ... t=%d %d (%d cycles)", tt+ido.offset, ido.id, (tt+ido.offset)/ido.id)
		// 	}
		// }
		if seq == len(input.idOffsets) {
			return tt
		}
		tt += d
	}
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part2 example 0", `939
7,13,x,x,59,x,31,19`, 1068781},
		{"part2 answer", advent.ReadFile(t, "input.txt"), 210612924879242},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}
