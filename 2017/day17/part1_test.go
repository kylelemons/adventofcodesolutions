package main_test

import (
	"testing"
)

func part1(offset int) int {
	const N = 2017
	list1 := make([]int, 0, N+1)
	list2 := make([]int, 0, N+1)

	list1 = append(list1, 0)

	var pos int
	for i := 1; i <= N; i++ {
		pos += offset
		pos %= len(list1)
		pos++
		list2 = list2[:0]
		list2 = append(list2, list1[:pos]...)
		list2 = append(list2, i)
		list2 = append(list2, list1[pos:]...)
		list1, list2 = list2, list1
	}
	return list1[(pos+1)%len(list1)]
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   int
		want int
	}{
		{"part1 example", 3, 638},
		{"part1", 355, 1912},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part1(test.in), test.want; got != want {
				t.Errorf("part1(%#v) = %#v, want %#v", test.in, got, want)
			}
		})
	}
}
