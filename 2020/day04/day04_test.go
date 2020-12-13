// Copyright 2020 Kyle Lemons
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
	"strings"
	"testing"

	"github.com/kylelemons/adventofcodesolutions/advent"
)

type Input struct {
	Passports []passport
}

type passport map[string]string

func parseInput(t *testing.T, in string) *Input {
	input := &Input{
		// ...
	}
	for _, v := range strings.Split(in, "\n\n") {
		last := make(passport)
		advent.Words(v).Each(func(i int, token advent.Scanner) {
			var k, v string
			token.Extract(t, `([^:]+):(.*)`, &k, &v)
			last[k] = v
		})
		input.Passports = append(input.Passports, last)
	}
	return input
}

func part1(t *testing.T, in string) (ret int) {
	input := parseInput(t, in)

next:
	for _, p := range input.Passports {
		for _, f := range []string{
			"byr", //(Birth Year)
			"iyr", //(Issue Year)
			"eyr", //(Expiration Year)
			"hgt", //(Height)
			"hcl", //(Hair Color)
			"ecl", //(Eye Color)
			"pid", //(Passport ID)
			// "cid", //(Country ID)
		} {
			if _, ok := p[f]; !ok {
				continue next
			}
		}
		ret++
	}

	return
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part1 example 0", `ecl:gry pid:860033327 eyr:2020 hcl:#fffffd
byr:1937 iyr:2017 cid:147 hgt:183cm

iyr:2013 ecl:amb cid:350 eyr:2023 pid:028048884
hcl:#cfa07d byr:1929

hcl:#ae17e1 iyr:2013
eyr:2024
ecl:brn pid:760753108 byr:1931
hgt:179cm

hcl:#cfa07d eyr:2025 pid:166559648
iyr:2011 ecl:brn hgt:59in`, 2},
		{"part1 anfswer", advent.ReadFile(t, "input.txt"), 239},
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

next:
	for _, p := range input.Passports {
		for _, f := range []string{
			"byr", //(Birth Year)
			"iyr", //(Issue Year)
			"eyr", //(Expiration Year)
			"hgt", //(Height)
			"hcl", //(Hair Color)
			"ecl", //(Eye Color)
			"pid", //(Passport ID)
			// "cid", //(Country ID)
		} {
			if _, ok := p[f]; !ok {
				continue next
			}
		}
		var byr, iyr, eyr int
		advent.Scanner(p["byr"]).Scan(t, &byr)
		advent.Scanner(p["iyr"]).Scan(t, &iyr)
		advent.Scanner(p["eyr"]).Scan(t, &eyr)
		if byr < 1920 || byr > 2002 {
			continue
		}
		if iyr < 2010 || iyr > 2020 {
			continue
		}
		if eyr < 2020 || eyr > 2030 {
			continue
		}

		hgt, hunit := 0, ""
		advent.Scanner(p["hgt"]).Extract(t, `(\d+)(.*)`, &hgt, &hunit)
		switch {
		case hunit == "cm" && hgt >= 150 && hgt <= 193:
		case hunit == "in" && hgt >= 59 && hgt <= 76:
		default:
			continue
		}

		if !advent.Scanner(p["hcl"]).CanExtract(t, `^#[0-9a-f]{6}$`) {
			continue
		}

		switch p["ecl"] {
		case "amb", "blu", "brn", "gry", "grn", "hzl", "oth":
		default:
			continue
		}

		if !advent.Scanner(p["pid"]).CanExtract(t, `^\d{9}$`) {
			continue
		}

		ret++
	}

	return
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part2 example 0", `eyr:1972 cid:100
hcl:#18171d ecl:amb hgt:170 pid:186cm iyr:2018 byr:1926

iyr:2019
hcl:#602927 eyr:1967 hgt:170cm
ecl:grn pid:012533040 byr:1946

hcl:dab227 iyr:2012
ecl:brn hgt:182cm pid:021572410 eyr:2020 byr:1992 cid:277

hgt:59cm ecl:zzz
eyr:2038 hcl:74454a iyr:2023
pid:3556412378 byr:2007

pid:087499704 hgt:74in ecl:grn iyr:2012 eyr:2030 byr:1980
hcl:#623a2f

eyr:2029 ecl:blu cid:129 byr:1989
iyr:2014 pid:896056539 hcl:#a97842 hgt:165cm

hcl:#888785
hgt:164cm byr:2001 iyr:2015 cid:88
pid:545766238 ecl:hzl
eyr:2022

iyr:2010 hgt:158cm hcl:#b6652a ecl:blu byr:1944 eyr:2021 pid:093154719`, 4},
		{"part2 anfswer", advent.ReadFile(t, "input.txt"), 188},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}
