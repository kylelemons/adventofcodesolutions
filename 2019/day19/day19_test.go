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
	"log"
	"testing"
	"time"

	"github.com/kylelemons/adventofcodesolutions/2019/intcode"
	"github.com/kylelemons/adventofcodesolutions/advent"
	"github.com/kylelemons/adventofcodesolutions/advent/coords"
)

type Input struct {
	prog *intcode.Program
}

func parseInput(t *testing.T, in string) *Input {
	input := &Input{
		prog: intcode.Compile(t, in),
	}
	return input
}

func part1(t *testing.T, in string) (ret int) {
	// input := parseInput(t, in)

	for r := 0; r < 50; r++ {
		for c := 0; c < 50; c++ {
			prog := intcode.Compile(t, in)

			input := make(chan int, 2)
			input <- r
			input <- c
			prog.Input = func() (v int) {
				select {
				case v = <-input:
				default:
					panic("no input")
				}
				return
			}

			prog.Output = func(v int) {
				if v == 1 {
					ret++
				}
			}
			prog.Run(t)
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
		// {"part1 example 0", "...", 0},
		{"part1 answer", advent.ReadFile(t, "input.txt"), 211},
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
	status := time.NewTicker(3 * time.Second)

	var rows []advent.RangeTracker

	var lastC advent.RangeTracker
	lastC.Track(0)

	// TODO(kevlar): Follow the lower-left edge and look for -100x+100 only

nextRow:
	for r := 0; ; r++ {
		var lastPulled bool
		var row advent.RangeTracker
		for c := lastC.Min; ; c++ {
			cur := coords.RC(r, c)

			select {
			case <-status.C:
				log.Println("Current: ", cur, "Last row: ", rows[len(rows)-1])
			default:
			}

			prog := intcode.Compile(t, in)
			input := make(chan int, 2)
			prog.Input = func() (v int) {
				select {
				case v = <-input:
				default:
					panic("no input")
				}
				return
			}

			input <- cur.X()
			input <- cur.Y()

			var pulled bool
			prog.Output = func(v int) {
				pulled = v == 1
			}
			prog.Run(t)

			if pulled {
				row.Track(c)
			}
			if !pulled && lastPulled {
				break
			}
			if c-lastC.Min > 10 && !pulled {
				break
			}
			lastPulled = pulled
		}
		lastC = row
		rows = append(rows, row)

		if !row.Valid {
			continue
		}
		if row.Delta() < 100 {
			continue
		}
		left := row.Min + 0
		right := left + 99

		for i := 1; i <= 100; i++ {
			r := rows[len(rows)-i]
			if r.Max < right {
				log.Printf("Only %d rows", i)
				continue nextRow
			}
		}
		topLeft := coords.RC(r-100+1, left)
		log.Println("Found ", topLeft)

		// buf := new(strings.Builder)
		// for r, rt := range rows {
		// 	for c := 0; c <= rt.Max; c++ {
		// 		switch {
		// 		case r == topLeft.R() && c == topLeft.C():
		// 			buf.WriteByte('O')
		// 		case c >= rt.Min:
		// 			buf.WriteByte('#')
		// 		default:
		// 			buf.WriteByte('.')
		// 		}
		// 	}
		// 	buf.WriteByte('\n')
		// }
		// log.Printf("Got:\n%s", buf)

		return topLeft.X()*10000 + topLeft.Y()
	}

	return
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		// {"part2 example 0", "...", 0},
		{"part2 answer", advent.ReadFile(t, "input.txt"), 8071006},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}
