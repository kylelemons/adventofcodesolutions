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
	"flag"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/kylelemons/adventofcodesolutions/advent"
)

var animate = flag.Bool("animate", false, "If set, animate the output")

type Input struct {
	Area [][]byte
}

func ParseInput(t *testing.T, in string) *Input {
	input := &Input{
		Area: advent.Split2D(in),
	}
	return input
}

func (i *Input) Get(r, c int) byte {
	if r < 0 || r >= len(i.Area) || c < 0 || c >= len(i.Area[0]) {
		return '.'
	}
	return i.Area[r][c]
}

func (i *Input) Next(r, c int) byte {
	var trees, yards int
	var got []byte
	for rr := -1; rr <= 1; rr++ {
		for cc := -1; cc <= 1; cc++ {
			if rr == 0 && cc == 0 {
				continue
			}
			got = append(got, i.Get(r+rr, c+cc))
			switch i.Get(r+rr, c+cc) {
			case '|':
				trees++
			case '#':
				yards++
			}
		}
	}
	switch here := i.Get(r, c); here {
	case '.':
		if trees >= 3 {
			return '|'
		}
		return here
	case '|':
		if yards >= 3 {
			return '#'
		}
		return here
	case '#':
		if yards < 1 || trees < 1 {
			return '.'
		}
		return here
	default:
		log.Fatalf("Unknown tile %q", here)
		return '.'
	}
}

func (i *Input) Advance() {
	next := advent.Make2D(len(i.Area), len(i.Area[0]))
	for r := range next {
		for c := range next[r] {
			next[r][c] = i.Next(r, c)
		}
	}
	i.Area = next
}

func (i *Input) ResourceValue() int {
	var trees, yards int
	for r := range i.Area {
		for c := range i.Area[r] {
			switch i.Area[r][c] {
			case '|':
				trees++
			case '#':
				yards++
			}
		}
	}
	return trees * yards
}

func (i *Input) String() string {
	out := new(strings.Builder)
	for r := range i.Area {
		fmt.Fprintf(out, "%s\n", i.Area[r])
	}
	return out.String()
}

func part1(t *testing.T, in string) (ret int) {
	input := ParseInput(t, in)

	t.Logf("Initial:\n%s", input)
	for i := 0; i < 10; i++ {
		input.Advance()
		// t.Logf("After %d minute:\n%s", i+1, input)
	}
	t.Logf("Final:\n%s", input)

	return input.ResourceValue()
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part1 example 0", strings.Trim(`
.#.#...|#.
.....#|##|
.|..|...#.
..|#.....#
#.#|||#|#|
...#.||...
.|....|...
||...#|.#|
|.||||..|.
...#.|..|.
`, "\n"), 1147},
		{"part1 answer", advent.ReadFile(t, "input.txt"), 645946},
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
	input := ParseInput(t, in)

	seen := make(map[string]int)

	t.Logf("Initial:\n%s", input)
	const N = 1000000000
	for i := 0; i < N; i++ {
		key := input.String()
		if last, ok := seen[key]; ok {
			cycle := i - last
			rem := (N - i) / cycle
			i += rem * cycle
			if rem > 0 {
				t.Logf("Fast forwarding %dx %d-minute cycles", rem, cycle)
			}
		} else {
			seen[key] = i
		}

		input.Advance()

		rv := input.ResourceValue()
		// log.Printf("After %d minute: %d", i+1, rv)
		if rv == 0 {
			break
		}
	}
	t.Logf("Final:\n%s", input)

	return input.ResourceValue()
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part2 answer", advent.ReadFile(t, "input.txt"), 227688},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

func TestAnimate(t *testing.T) {
	// go test -test.run=Animate -test.timeout=10s --animate

	if !*animate {
		t.Skip("Run with --animate")
	}

	in := advent.ReadFile(t, "input.txt")
	input := ParseInput(t, in)
	seen := make(map[string]bool)

	next := time.NewTicker(1 * time.Second / 30)
	defer next.Stop()

	fmt.Printf("\033[2J") // cls
	for now := 0; ; now++ {
		if input.ResourceValue() == 0 {
			t.Fatalf("Ran out of trees!")
			return
		}
		key := input.String()
		fmt.Printf("\033[;H\nTime = %d\n%s\n\n", now, key)
		if _, ok := seen[key]; ok {
			return
		}
		input.Advance()
		<-next.C
	}
}
