package aocday

import (
	"testing"
)

func part1(t *testing.T, in string) string {
	return in + "foo"
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"part1 example", "", ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part1(t, test.in), test.want; got != want {
				t.Errorf("part1(%#v) = %#v, want %#v", test.in, got, want)
			}
		})
	}
}
