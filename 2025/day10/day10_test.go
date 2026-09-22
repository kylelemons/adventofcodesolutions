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
	"bufio"
	_ "embed"
	"fmt"
	"math"
	"math/bits"
	"strconv"
	"strings"
	"testing"
)

var _ = fmt.Println // keep imported

type Machine struct {
	Pattern Lights
	Buttons []Lights
}

type Input struct {
	Machines []Machine
}

type Lights struct {
	Bits   uint32
	Length int
}

// String returns the pattern representation (e.g., "[..#.##]").
func (l Lights) String() string {
	var sb strings.Builder
	sb.WriteByte('[')
	for i := 0; i < l.Length; i++ {
		if (l.Bits & (1 << i)) != 0 {
			sb.WriteByte('#')
		} else {
			sb.WriteByte('.')
		}
	}
	sb.WriteByte(']')
	return sb.String()
}

// Print outputs the lights representation to stdout.
func (l Lights) Print() {
	fmt.Println(l.String())
}

// input2pattern parses light/indicator patterns like "[..#.##]" into a Lights struct.
// The first light is stored in the least significant bit (bit 0).
func input2pattern(input string) Lights {
	s := strings.Trim(input, "[] \t\r\n")
	var bits uint32
	for i, ch := range s {
		if ch == '#' {
			bits |= 1 << i
		}
	}
	return Lights{
		Bits:   bits,
		Length: len(s),
	}
}

// input2button parses comma-separated index lists like "(1,2,5)" or "[1,2,5]"
// and sets the corresponding bits in a Lights struct of the given length.
func input2button(input string, length int) Lights {
	s := strings.Trim(input, "[]() \t\r\n")
	if s == "" {
		return Lights{Length: length}
	}
	var bits uint32
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if idx, err := strconv.Atoi(part); err == nil && idx >= 0 && idx < length {
			bits |= 1 << idx
		}
	}
	return Lights{
		Bits:   bits,
		Length: length,
	}
}

// setFromXor returns the XOR combination of lights a and b.
func setFromXor(a, b Lights) Lights {
	return Lights{
		Bits:   a.Bits ^ b.Bits,
		Length: a.Length,
	}
}

// Xor is a convenience method on Lights.
func (l Lights) Xor(other Lights) Lights {
	return Lights{
		Bits:   l.Bits ^ other.Bits,
		Length: l.Length,
	}
}

// pushButtons XORs together the Lights corresponding to the set bits in toPush
// (where button 0 corresponds to the least significant bit, button 1 to bit 1, etc.).
func pushButtons(buttons []Lights, toPush uint32) Lights {
	var result Lights
	if len(buttons) > 0 {
		result.Length = buttons[0].Length
	}

	for i, btn := range buttons {
		if (toPush & (1 << i)) != 0 {
			result.Bits ^= btn.Bits
		}
	}
	return result
}

func parseInput(t *testing.T, in string) *Input {
	input := &Input{
		// ...
	}
	lines := bufio.NewScanner(strings.NewReader(strings.TrimSpace(in)))
	for lines.Scan() {
		words := strings.Split(lines.Text(), " ")
		patRaw, words := words[0], words[1:]
		words = words[:len(words)-1] // take off joltage

		pat := input2pattern(patRaw)
		buttons := make([]Lights, len(words))
		for i := range buttons {
			buttons[i] = input2button(words[i], pat.Length)
		}

		// fmt.Println("pat", patRaw, "words", words)
		// fmt.Println("pat", pat, "buttons", buttons)

		input.Machines = append(input.Machines, Machine{
			Pattern: pat,
			Buttons: buttons,
		})
	}
	return input
}

func part1(t *testing.T, in string) (ret int) {
	input := parseInput(t, in)
	_ = input

	for _, input := range input.Machines {

		minN := math.MaxInt
		for i := range uint32(1) << len(input.Buttons) {
			output := pushButtons(input.Buttons, i)
			if output.Bits == input.Pattern.Bits {
				n := bits.OnesCount32(i)
				if n < minN {
					minN = n
				}

				// fmt.Printf("Pushing %d %b : %v\n", n, i, output)
			}
		}
		ret += minN
	}

	return
}

//go:embed input.txt
var inputFile string

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part1 example 0", "[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}\n[...#.] (0,2,3,4) (2,3) (0,4) (0,1,2) (1,2,3,4) {7,5,12,7,2}\n[.###.#] (0,1,2,3,4) (0,3,4) (0,1,2,4,5) (1,2) {10,11,11,5,10,5}", 7},
		{"part1 answer", inputFile, 578},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part1(t, test.in), test.want; got != want {
				t.Errorf("part1(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}
