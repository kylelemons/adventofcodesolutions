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

	DesiredJoltages []int
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

// input2joltage parses brace-enclosed comma-separated integers (e.g., "{1,2,3,4}") into an []int.
func input2joltage(input string) []int {
	s := strings.Trim(input, "{} \t\r\n")
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	res := make([]int, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if v, err := strconv.Atoi(part); err == nil {
			res = append(res, v)
		}
	}
	return res
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
		desJoltRaw := words[len(words)-1]
		words = words[:len(words)-1] // take off joltage

		pat := input2pattern(patRaw)
		buttons := make([]Lights, len(words))
		for i := range buttons {
			buttons[i] = input2button(words[i], pat.Length)
		}
		desired := input2joltage(desJoltRaw)

		// fmt.Println("pat", patRaw, "words", words)
		// fmt.Println("pat", pat, "buttons", buttons, "joltage", desired)

		input.Machines = append(input.Machines, Machine{
			Pattern:         pat,
			Buttons:         buttons,
			DesiredJoltages: desired,
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

func part2(t *testing.T, in string) (ret int) {
	input := parseInput(t, in)
	_ = input

	for _, input := range input.Machines {
		ret += part2machine(input)
	}

	return
}

func sum(joltages []int) int {
	total := 0
	for _, v := range joltages {
		total += v
	}
	return total
}

// addJoltageFromButton clones joltages and increments each position toggled
// by the buttons specified in mask.
func addJoltageFromButton(joltages []int, buttons []Lights, mask uint32) []int {
	next := make([]int, len(joltages))
	copy(next, joltages)

	for i, btn := range buttons {
		if (mask & (1 << i)) != 0 { // if the button is pressed based on the mask
			for b := 0; b < btn.Length; b++ {
				if (btn.Bits & (1 << b)) != 0 { // if the mask says to increase this joltage
					next[b]++
				}
			}
		}
	}
	return next
}

type trackedState struct {
	CurrentJoltages  []int
	ButtonsPressed   int
	RemainingJoltage int
}

type trackedStateHeap []trackedState

func (h trackedStateHeap) Len() int { return len(h) }
func (h trackedStateHeap) Less(i, j int) bool {
	// Primary: fewest buttons pressed (guarantees shortest path / min presses)
	if h[i].ButtonsPressed != h[j].ButtonsPressed {
		return h[i].ButtonsPressed < h[j].ButtonsPressed
	}
	// Tie-breaker: lowest remaining joltage (closest to goal)
	return h[i].RemainingJoltage < h[j].RemainingJoltage
}
func (h trackedStateHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *trackedStateHeap) Push(x any) {
	*h = append(*h, x.(trackedState))
}

func (h *trackedStateHeap) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[0 : n-1]
	return item
}

func part2machine(input Machine) int {
	numEquations := input.Pattern.Length
	numVars := len(input.Buttons)

	// Build augmented matrix: [A | b]
	// Matrix has numEquations rows and (numVars + 1) columns.
	mat := make([][]float64, numEquations)
	for i := 0; i < numEquations; i++ {
		mat[i] = make([]float64, numVars+1)
		for j, btn := range input.Buttons {
			if (btn.Bits & (1 << i)) != 0 {
				mat[i][j] = 1.0
			}
		}
		mat[i][numVars] = float64(input.DesiredJoltages[i])
	}

	// Gaussian elimination to Reduced Row Echelon Form (RREF)
	pivotRow := 0
	pivotCols := make([]int, 0)
	varColIsPivot := make([]bool, numVars)

	const eps = 1e-9

	for col := 0; col < numVars && pivotRow < numEquations; col++ {
		// Find pivot
		sel := -1
		for r := pivotRow; r < numEquations; r++ {
			if math.Abs(mat[r][col]) > eps {
				sel = r
				break
			}
		}
		if sel == -1 {
			continue
		}

		// Swap rows
		mat[pivotRow], mat[sel] = mat[sel], mat[pivotRow]

		// Normalize pivot row
		div := mat[pivotRow][col]
		for c := col; c <= numVars; c++ {
			mat[pivotRow][c] /= div
		}

		// Eliminate all other rows in this column
		for r := 0; r < numEquations; r++ {
			if r != pivotRow && math.Abs(mat[r][col]) > eps {
				factor := mat[r][col]
				for c := col; c <= numVars; c++ {
					mat[r][c] -= factor * mat[pivotRow][c]
				}
			}
		}

		pivotCols = append(pivotCols, col)
		varColIsPivot[col] = true
		pivotRow++
	}

	// Check for contradictory rows (0 = non-zero)
	for r := pivotRow; r < numEquations; r++ {
		if math.Abs(mat[r][numVars]) > eps {
			return math.MaxInt // No solution
		}
	}

	// Identify free variables
	var freeCols []int
	for col := 0; col < numVars; col++ {
		if !varColIsPivot[col] {
			freeCols = append(freeCols, col)
		}
	}

	// Upper bounds for free variables: cannot exceed min desired joltage for the lights it toggles
	freeBounds := make([]int, len(freeCols))
	for idx, fcol := range freeCols {
		bound := math.MaxInt
		btn := input.Buttons[fcol]
		for b := 0; b < btn.Length; b++ {
			if (btn.Bits & (1 << b)) != 0 {
				if input.DesiredJoltages[b] < bound {
					bound = input.DesiredJoltages[b]
				}
			}
		}
		if bound == math.MaxInt {
			bound = 0
		}
		freeBounds[idx] = bound
	}

	minTotalPresses := math.MaxInt

	// Search all combinations of values for the free variables
	var searchFree func(freeIdx int, currentFreeVals []int)
	searchFree = func(freeIdx int, currentFreeVals []int) {
		if freeIdx == len(freeCols) {
			// Compute pivot variable values
			x := make([]int, numVars)
			for i, fcol := range freeCols {
				x[fcol] = currentFreeVals[i]
			}

			valid := true
			for r, pcol := range pivotCols {
				val := mat[r][numVars]
				for i, fcol := range freeCols {
					val -= mat[r][fcol] * float64(currentFreeVals[i])
				}
				roundVal := math.Round(val)
				if math.Abs(val-roundVal) > 1e-5 || roundVal < -eps {
					valid = false
					break
				}
				x[pcol] = int(roundVal)
			}

			if valid {
				total := 0
				for _, presses := range x {
					total += presses
				}
				if total < minTotalPresses {
					minTotalPresses = total
				}
			}
			return
		}

		for val := 0; val <= freeBounds[freeIdx]; val++ {
			// Pruning: if current free variable presses already exceed the minimum found
			currentSum := val
			for _, v := range currentFreeVals {
				currentSum += v
			}
			if currentSum >= minTotalPresses {
				break
			}

			searchFree(freeIdx+1, append(currentFreeVals, val))
		}
	}

	searchFree(0, nil)

	return minTotalPresses
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part2 example 0", "[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}\n[...#.] (0,2,3,4) (2,3) (0,4) (0,1,2) (1,2,3,4) {7,5,12,7,2}\n[.###.#] (0,1,2,3,4) (0,3,4) (0,1,2,4,5) (1,2) {10,11,11,5,10,5}", 33},
		{"part2 answer", inputFile, 20709},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got != want {
				t.Errorf("part2(...)\n = %#v, want %#v", got, want)
			}
		})
	}
}
