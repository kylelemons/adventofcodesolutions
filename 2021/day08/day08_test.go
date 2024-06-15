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
	"strings"
	"testing"

	"github.com/kylelemons/adventofcodesolutions/advent"
)

type Input struct {
	Lines []Line
}

type Line struct {
	Combos []string
	Digits []string
}

func parseInput(t *testing.T, in string) *Input {
	input := &Input{
		// ...
	}
	advent.Lines(in).Each(func(i int, v advent.Scanner) {
		combos, digits, _ := strings.Cut(string(v), " | ")
		input.Lines = append(input.Lines, Line{
			Combos: strings.Fields(combos),
			Digits: strings.Fields(digits),
		})
	})
	return input
}

func part1(t *testing.T, in string) (ret int) {
	input := parseInput(t, in)

	for _, line := range input.Lines {
		fmt.Println(line)
		for _, digit := range line.Digits {
			switch len(digit) {
			case 2: // 1
				ret++
			case 4: // 4
				ret++
			case 3: // 7
				ret++
			case 7: // 8
				ret++
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
		{"part1 example 0", `be cfbegad cbdgef fgaecd cgeb fdcge agebfd fecdb fabcd edb | fdgacbe cefdb cefbgd gcbe
		edbfga begcd cbg gc gcadebf fbgde acbgfd abcde gfcbed gfec | fcgedb cgb dgebacf gc
		fgaebd cg bdaec gdafb agbcfd gdcbef bgcad gfac gcb cdgabef | cg cg fdcagb cbg
		fbegcd cbd adcefb dageb afcb bc aefdc ecdab fgdeca fcdbega | efabcd cedba gadfec cb
		aecbfdg fbg gf bafeg dbefa fcge gcbea fcaegb dgceab fcbdga | gecf egdcabf bgf bfgea
		fgeab ca afcebg bdacfeg cfaedg gcfdb baec bfadeg bafgc acf | gebdcfa ecba ca fadegcb
		dbcfg fgd bdegcaf fgec aegbdf ecdfab fbedc dacgb gdcebf gf | cefg dcbef fcge gbcadfe
		bdfegc cbegaf gecbf dfcage bdacg ed bedf ced adcbefg gebcd | ed bcgafe cdgba cbgef
		egadfb cdbfeg cegd fecab cgb gbdefca cg fgcdab egfdb bfceg | gbdfcae bgc cg cgb
		gcafb gcf dcaebfg ecagb gf abcdeg gaef cafbge fdbac fegbdc | fgae cfgab fg bagce`, 26},
		{"part1 answer", advent.ReadFile(t, "input.txt"), 548},
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

	digitSegmentStrings := [10]string{
		0: "abcefg",
		1: "cf",
		2: "acdeg",
		3: "acdfg",
		4: "bcdf",
		5: "abdfg",
		6: "abdefg",
		7: "acf",
		8: "abcdefg",
		9: "abcdfg",
	}

	segmentsToDigit := make(map[[7]bool]int)
	for digit, segmentString := range digitSegmentStrings {
		var segs [7]bool
		for _, c := range segmentString {
			segs[c-'a'] = true
		}
		segmentsToDigit[segs] = digit
	}

	for _, line := range input.Lines {
		fmt.Println(line)

		possible := [7][7]bool{}
		for wire := range possible {
			for seg := range possible[wire] {
				possible[wire][seg] = true
			}
		}
		mustLightDigit := func(digitWires string, actualDigit int) {
			fmt.Println(" must light up", actualDigit)
			digitSegments := digitSegmentStrings[actualDigit]
			// Consider all wires lit up for this digit:
			for _, wireLetter := range digitWires {
				// Consider all possible segments that could be lit:
				for _, segmentLetter := range "abcdefg" {
					// This wire can only connect to this segment if it's one of the segments
					// implied by the segment count.  So, if it's not:
					if !strings.ContainsRune(digitSegments, rune(segmentLetter)) {
						// fmt.Printf("  wire %c can't light up segment %c\n", wireLetter, segmentLetter)
						wire, seg := wireLetter-'a', segmentLetter-'a'
						possible[wire][seg] = false
					}
				}
			}
			// Consider all possible input wires:
			for _, wireLetter := range "abcdefg" {
				// If this wire is lit for this observed digit:
				if !strings.ContainsRune(digitWires, rune(wireLetter)) {
					// then this wire can't match up with any segment associated with
					// the only possible digit for this number of segements:
					for _, segmentLetter := range digitSegments {
						// fmt.Printf("  wire %c can't light up segment %c\n", wireLetter, segmentLetter)
						wire, seg := wireLetter-'a', segmentLetter-'a'
						possible[wire][seg] = false
					}
				}
			}
		}

		for _, digit := range append(line.Combos, line.Digits...) {
			switch len(digit) {
			case 2:
				mustLightDigit(digit, 1)
			case 4:
				mustLightDigit(digit, 4)
			case 3:
				mustLightDigit(digit, 7)
			case 7:
				mustLightDigit(digit, 8)
			default:
				continue
			}

		}

		countSet := func(segs [7]bool) (c int) {
			for _, s := range segs {
				if s {
					c++
				}
			}
			return
		}
		fmt.Println("Initial observations require:")
		for i, poss := range possible {
			fmt.Printf(" %c (%d) %v\n", 'a'+i, countSet(poss), poss)
		}

		wire2seg := make(map[int]int)
		for {
			var progress bool

			for wire, possibleSegments := range possible {
				if countSet(possibleSegments) != 1 {
					// no deduction yet
					continue
				}
				if _, ok := wire2seg[wire]; ok {
					// already deduced
					continue
				}

				// Get the segment that is matched
				var seg int
				for s := range possibleSegments {
					if possible[wire][s] {
						seg = s
					}
				}
				wire2seg[wire] = seg
				progress = true

				// Eliminate this as a possible segment for other wires
				for otherWire := range possible {
					if wire != otherWire {
						possible[otherWire][seg] = false
					}
				}
			}

			fmt.Println(wire2seg, progress)
			if !progress {
				break
			}
		}

		if len(wire2seg) != 7 {
			fmt.Println("No mapping for line:", line)
			for i, poss := range possible {
				fmt.Printf(" %c (%d) %v\n", 'a'+i, countSet(poss), poss)
			}
			continue
		}

		var output int
		for _, digitStr := range line.Digits {
			var segments [7]bool
			for _, wireLetter := range digitStr {
				segments[wire2seg[int(wireLetter-'a')]] = true
			}
			digit, ok := segmentsToDigit[segments]
			if !ok {
				fmt.Println("No mapping for segememts", segments)
				continue
			}
			output *= 10
			output += digit
		}
		ret += output
	}

	return
}

func TestPart2(t *testing.T) {
	t.Skip("This doesn't work yet")

	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part2 example 0", `be cfbegad cbdgef fgaecd cgeb fdcge agebfd fecdb fabcd edb | fdgacbe cefdb cefbgd gcbe
		edbfga begcd cbg gc gcadebf fbgde acbgfd abcde gfcbed gfec | fcgedb cgb dgebacf gc
		fgaebd cg bdaec gdafb agbcfd gdcbef bgcad gfac gcb cdgabef | cg cg fdcagb cbg
		fbegcd cbd adcefb dageb afcb bc aefdc ecdab fgdeca fcdbega | efabcd cedba gadfec cb
		aecbfdg fbg gf bafeg dbefa fcge gcbea fcaegb dgceab fcbdga | gecf egdcabf bgf bfgea
		fgeab ca afcebg bdacfeg cfaedg gcfdb baec bfadeg bafgc acf | gebdcfa ecba ca fadegcb
		dbcfg fgd bdegcaf fgec aegbdf ecdfab fbedc dacgb gdcebf gf | cefg dcbef fcge gbcadfe
		bdfegc cbegaf gecbf dfcage bdacg ed bedf ced adcbefg gebcd | ed bcgafe cdgba cbgef
		egadfb cdbfeg cegd fecab cgb gbdefca cg fgcdab egfdb bfceg | gbdfcae bgc cg cgb
		gcafb gcf dcaebfg ecagb gf abcdeg gaef cafbge fdbac fegbdc | fgae cfgab fg bagce`, 26},
		// {"part2 answer", advent.ReadFile(t, "input.txt"), 548},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}
