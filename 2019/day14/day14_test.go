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
	"sort"
	"testing"

	"github.com/kylelemons/adventofcodesolutions/advent"
)

type Reactant struct {
	Qty  int
	Chem string
}

type Input struct {
	Reactions map[Reactant][]Reactant // [out] = [in ...]
}

func parseInput(t *testing.T, in string) *Input {
	input := &Input{
		Reactions: make(map[Reactant][]Reactant),
	}
	advent.Lines(in).Extract(t, `(.*) => (.*)`, func(ins, outs string) {
		var out Reactant
		advent.Scanner(outs).Scan(t, advent.Fields(&out)...)
		advent.Split(ins, ',').Each(func(i int, line advent.Scanner) {
			var in Reactant
			line.Scan(t, advent.Fields(&in)...)
			input.Reactions[out] = append(input.Reactions[out], in)
		})
	})
	return input
}

func (input *Input) orePerNFuel(t *testing.T, N int) (ore int) {
	target := map[string]int{
		"FUEL": N,
	}
	leftover := map[string]int{}

	for {
		if len(target) == 1 && target["ORE"] > 0 {
			return target["ORE"]
		}

		chems := make([]string, 0, len(target))
		for chem := range target {
			chems = append(chems, chem)
		}
		sort.Strings(chems)

		for _, chem := range chems {
			if chem == "ORE" {
				continue
			}

			var inputs []Reactant
			var output Reactant
			for out, ins := range input.Reactions {
				if out.Chem != chem {
					continue
				}
				inputs, output = ins, out
			}
			if len(inputs) == 0 {
				t.Fatalf("no raction for %q", chem)
			}

			// Figure out how many we actually need to make:
			needed := target[chem] - leftover[chem]

			// How many reactions are required: (integer divide-with-ceil)
			runs := (needed + output.Qty - 1) / output.Qty

			// How many does that make?
			made := runs * output.Qty

			// How many are leftover?
			left := made - needed
			leftover[chem] += left

			// t.Logf("Running reaction %v <= %v x %d for %d %s (%d of %d left over)", output, inputs, runs, needed, chem, left, made)
			delete(target, chem)

			for _, in := range inputs {
				target[in.Chem] += runs * in.Qty
				for leftover[in.Chem] > 0 && target[in.Chem] > 0 {
					leftover[in.Chem]--
					target[in.Chem]--
				}
				if leftover[in.Chem] == 0 {
					delete(leftover, in.Chem)
				}
				if target[in.Chem] == 0 {
					delete(target, in.Chem)
				}
			}
		}
	}
}

func part1(t *testing.T, in string) (ret int) {
	input := parseInput(t, in)
	return input.orePerNFuel(t, 1)
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part1 example 0", `9 ORE => 2 A
8 ORE => 3 B
7 ORE => 5 C
3 A, 4 B => 1 AB
5 B, 7 C => 1 BC
4 C, 1 A => 1 CA
2 AB, 3 BC, 4 CA => 1 FUEL`, 165},
		{"part1 example 0", `157 ORE => 5 NZVS
165 ORE => 6 DCFZ
44 XJWVT, 5 KHKGT, 1 QDVJ, 29 NZVS, 9 GPVTF, 48 HKGWZ => 1 FUEL
12 HKGWZ, 1 GPVTF, 8 PSHF => 9 QDVJ
179 ORE => 7 PSHF
177 ORE => 5 HKGWZ
7 DCFZ, 7 PSHF => 2 XJWVT
165 ORE => 2 GPVTF
3 DCFZ, 7 NZVS, 5 HKGWZ, 10 PSHF => 8 KHKGT`, 13312},
		{"part1 example 0", `2 VPVL, 7 FWMGM, 2 CXFTF, 11 MNCFX => 1 STKFG
17 NVRVD, 3 JNWZP => 8 VPVL
53 STKFG, 6 MNCFX, 46 VJHF, 81 HVMC, 68 CXFTF, 25 GNMV => 1 FUEL
22 VJHF, 37 MNCFX => 5 FWMGM
139 ORE => 4 NVRVD
144 ORE => 7 JNWZP
5 MNCFX, 7 RFSQX, 2 FWMGM, 2 VPVL, 19 CXFTF => 3 HVMC
5 VJHF, 7 MNCFX, 9 VPVL, 37 CXFTF => 6 GNMV
145 ORE => 6 MNCFX
1 NVRVD => 8 CXFTF
1 VJHF, 6 MNCFX => 4 RFSQX
176 ORE => 6 VJHF`, 180697},
		{"part1 example 0", `171 ORE => 8 CNZTR
7 ZLQW, 3 BMBT, 9 XCVML, 26 XMNCP, 1 WPTQ, 2 MZWV, 1 RJRHP => 4 PLWSL
114 ORE => 4 BHXH
14 VRPVC => 6 BMBT
6 BHXH, 18 KTJDG, 12 WPTQ, 7 PLWSL, 31 FHTLT, 37 ZDVW => 1 FUEL
6 WPTQ, 2 BMBT, 8 ZLQW, 18 KTJDG, 1 XMNCP, 6 MZWV, 1 RJRHP => 6 FHTLT
15 XDBXC, 2 LTCX, 1 VRPVC => 6 ZLQW
13 WPTQ, 10 LTCX, 3 RJRHP, 14 XMNCP, 2 MZWV, 1 ZLQW => 1 ZDVW
5 BMBT => 4 WPTQ
189 ORE => 9 KTJDG
1 MZWV, 17 XDBXC, 3 XCVML => 2 XMNCP
12 VRPVC, 27 CNZTR => 2 XDBXC
15 KTJDG, 12 BHXH => 5 XCVML
3 BHXH, 2 VRPVC => 7 MZWV
121 ORE => 7 VRPVC
7 XCVML => 6 RJRHP
5 BHXH, 4 VRPVC => 5 LTCX`, 2210736},
		{"part1 answer", advent.ReadFile(t, "input.txt"), 870051},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part1(t, test.in), test.want; got != want {
				t.Errorf("part1(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

func part2(t *testing.T, in string) (fuel int) {
	input := parseInput(t, in)

	target := 1000000000000
	return sort.Search(target, func(fuel int) bool {
		return input.orePerNFuel(t, fuel) > target
	}) - 1 // Search returns the first that's true, so subtract 1
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{

		{"part2 example 0", `157 ORE => 5 NZVS
165 ORE => 6 DCFZ
44 XJWVT, 5 KHKGT, 1 QDVJ, 29 NZVS, 9 GPVTF, 48 HKGWZ => 1 FUEL
12 HKGWZ, 1 GPVTF, 8 PSHF => 9 QDVJ
179 ORE => 7 PSHF
177 ORE => 5 HKGWZ
7 DCFZ, 7 PSHF => 2 XJWVT
165 ORE => 2 GPVTF
3 DCFZ, 7 NZVS, 5 HKGWZ, 10 PSHF => 8 KHKGT`, 82892753},
		{"part2 example 1", `2 VPVL, 7 FWMGM, 2 CXFTF, 11 MNCFX => 1 STKFG
17 NVRVD, 3 JNWZP => 8 VPVL
53 STKFG, 6 MNCFX, 46 VJHF, 81 HVMC, 68 CXFTF, 25 GNMV => 1 FUEL
22 VJHF, 37 MNCFX => 5 FWMGM
139 ORE => 4 NVRVD
144 ORE => 7 JNWZP
5 MNCFX, 7 RFSQX, 2 FWMGM, 2 VPVL, 19 CXFTF => 3 HVMC
5 VJHF, 7 MNCFX, 9 VPVL, 37 CXFTF => 6 GNMV
145 ORE => 6 MNCFX
1 NVRVD => 8 CXFTF
1 VJHF, 6 MNCFX => 4 RFSQX
176 ORE => 6 VJHF`, 5586022},
		{"part2 example 2", `171 ORE => 8 CNZTR
7 ZLQW, 3 BMBT, 9 XCVML, 26 XMNCP, 1 WPTQ, 2 MZWV, 1 RJRHP => 4 PLWSL
114 ORE => 4 BHXH
14 VRPVC => 6 BMBT
6 BHXH, 18 KTJDG, 12 WPTQ, 7 PLWSL, 31 FHTLT, 37 ZDVW => 1 FUEL
6 WPTQ, 2 BMBT, 8 ZLQW, 18 KTJDG, 1 XMNCP, 6 MZWV, 1 RJRHP => 6 FHTLT
15 XDBXC, 2 LTCX, 1 VRPVC => 6 ZLQW
13 WPTQ, 10 LTCX, 3 RJRHP, 14 XMNCP, 2 MZWV, 1 ZLQW => 1 ZDVW
5 BMBT => 4 WPTQ
189 ORE => 9 KTJDG
1 MZWV, 17 XDBXC, 3 XCVML => 2 XMNCP
12 VRPVC, 27 CNZTR => 2 XDBXC
15 KTJDG, 12 BHXH => 5 XCVML
3 BHXH, 2 VRPVC => 7 MZWV
121 ORE => 7 VRPVC
7 XCVML => 6 RJRHP
5 BHXH, 4 VRPVC => 5 LTCX`, 460664},
		{"part2 answer", advent.ReadFile(t, "input.txt"), 1863741},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}
