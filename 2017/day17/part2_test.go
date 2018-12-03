package main_test

import (
	"testing"
)

func part2(offset int) (last int) {
	var pos int
	length := 1
	for i := 1; i <= 50000000; i++ {
		pos += offset
		pos %= length
		pos++
		length++

		if pos == 1 {
			last = i
		}
	}
	return
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name string
		in   int
		want int
	}{
		{"part2", 355, 21066990},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(test.in), test.want; got != want {
				t.Errorf("part2(%#v) = %#v, want %#v", test.in, got, want)
			}
		})
	}
}
