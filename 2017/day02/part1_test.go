package aocday

import (
	"sort"
	"strconv"
	"strings"
	"testing"
)

func part1(t *testing.T, in string) (sum int) {
	lines := strings.Split(in, "\n")
	for _, line := range lines {
		cells := strings.Fields(line)
		var max, min int
		for i, cell := range cells {
			val, err := strconv.Atoi(cell)
			if err != nil {
				t.Fatalf("cell %q is not a number: %s", cell, err)
			}
			if i == 0 {
				max = val
				min = val
				continue
			}
			if val > max {
				max = val
			}
			if val < min {
				min = val
			}
		}
		sum += max - min
	}
	return
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part1 example", `5 1 9 5
7 5 3
2 4 6 8`, 18},
		{"part1", read(t, "input.txt"), 58975},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part1(t, test.in), test.want; got != want {
				t.Errorf("part1(%#v) = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

func part2(t *testing.T, in string) (sum int) {
	lines := strings.Split(in, "\n")
	for _, line := range lines {
		cells := strings.Fields(line)
		var vals []int
		for _, cell := range cells {
			val, err := strconv.Atoi(cell)
			if err != nil {
				t.Fatalf("cell %q is not a number: %s", cell, err)
			}
			vals = append(vals, val)
		}
		sort.Ints(vals)
		for i, denom := range vals {
			for _, num := range vals[i+1:] {
				if num%denom == 0 {
					sum += num / denom
					t.Log(num, denom, num/denom)
				}
			}
		}
	}
	return
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part2 example", `5 9 2 8
9 4 7 3
3 8 6 5`, 9},
		{"part2", read(t, "input.txt"), 308},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got != want {
				t.Errorf("part2(%#v) = %#v, want %#v", test.in, got, want)
			}
		})
	}
}
