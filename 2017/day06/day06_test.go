package main

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func cycles(s string) (index int, size int) {
	var buckets []int
	for _, f := range strings.Fields(s) {
		v, err := strconv.Atoi(f)
		if err != nil {
			panic(err)
		}
		buckets = append(buckets, v)
	}

	hash := func(b []int) string {
		return fmt.Sprint(b)
	}
	seen := map[string]int{}

	for iter := 0; true; iter++ {
		h := hash(buckets)
		if at, ok := seen[h]; ok {
			return iter, iter - at
		}
		seen[h] = iter

		at, max := 0, buckets[0]
		for i, cur := range buckets {
			if cur > max {
				at, max = i, cur
			}
		}

		buckets[at] = 0
		at++

		for max > 0 {
			buckets[at%len(buckets)]++
			max--
			at++
		}
	}
	panic("unreachable")
}

func TestCycles(t *testing.T) {
	for _, input := range []struct {
		str          string
		part1, part2 int
	}{
		{"0 2 7 0", 5, 4},
		{"4 1 15 12 0 9 9 5 5 8 7 3 14 5 12 3", 6681, 2392},
		{"10 3 15 10 5 15 5 15 9 2 5 8 5 2 3 6", 14029, 2765},
	} {
		part1, part2 := cycles(input.str)
		if got, want := part1, input.part1; want != -1 && got != want {
			t.Errorf("part1(%q) = %v, want %v\n", input.str, got, want)
		} else {
			t.Logf("part1(%q) = %v [OK]", input.str, part1)
		}
		if got, want := part2, input.part2; want != -1 && got != want {
			t.Errorf("part2(%q) = %v, want %v\n", input.str, got, want)
			continue
		} else {
			t.Logf("part2(%q) = %v [OK]", input.str, part2)
		}
	}
}
