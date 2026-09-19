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
	"bytes"
	"strconv"
	"strings"
	"testing"

	"github.com/kylelemons/adventofcodesolutions/advent"
)

type Range struct {
	Lo, Hi int
}
type Input struct {
	Ranges []Range
	MaxLen int
}

func parseInput(t *testing.T, in string) *Input {
	input := &Input{
		// ...
	}
	for rng := range strings.SplitSeq(in, ",") {
		low, hi, ok := strings.Cut(rng, "-")
		if !ok {
			t.Fatalf("invalid range: %q", rng)
		}
		input.Ranges = append(input.Ranges, Range{
			Lo: must(strconv.Atoi(low))(t),
			Hi: must(strconv.Atoi(hi))(t),
		})
		if len(low) > input.MaxLen {
			input.MaxLen = len(low)
		}
		if len(hi) > input.MaxLen {
			input.MaxLen = len(hi)
		}
	}
	return input
}

func must[T any](v T, err error) func(t *testing.T) T {
	return func(t *testing.T) T {
		if err != nil {
			t.Fatalf("creating %T: error: %v", v, err)
		}
		return v
	}
}

func part1(t *testing.T, in string) (ret int) {
	input := parseInput(t, in)

	buf := make([]byte, 0, input.MaxLen)
	for _, r := range input.Ranges {
		for i := r.Lo; i <= r.Hi; i++ {
			v := strconv.AppendInt(buf[:0], int64(i), 10)
			half := len(v) / 2
			left, right := v[:half], v[half:]
			if bytes.Equal(left, right) {
				ret += i
			}
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
		{"part1 example 0", "11-22,95-115,998-1012,1188511880-1188511890,222220-222224,1698522-1698528,446443-446449,38593856-38593862,565653-565659,824824821-824824827,2121212118-2121212124", 1227775554},
		{"part1 answer", advent.ReadFile(t, "input.txt"), 44487518055},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part1(t, test.in), test.want; got != want {
				t.Errorf("part1(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}
