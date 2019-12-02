package aocday

import (
	"fmt"
	"testing"
)

func TestScan(t *testing.T) {
	var reg1, reg2 string
	var val, offset int

	ptrs := func(v ...interface{}) []interface{} { return v }

	t.Run("positive", func(t *testing.T) {
		tests := []struct {
			desc, input, re string
			ptrs            []interface{}
			reg1, reg2      string
			val, offset     int
		}{
			{
				desc:   "basic",
				input:  "jne a 1",
				re:     `jne (\w+) (\d+)`,
				ptrs:   ptrs(&reg1, &offset),
				reg1:   "a",
				offset: 1,
			},
		}

		for _, test := range tests {
			reg1, reg2 = "", ""
			val, offset = 0, 0

			t.Run(fmt.Sprintf("scanner(%q).scan(%#q)", test.input, test.re), func(t *testing.T) {
				if got, want := scanner(test.input).scan(t, test.re, test.ptrs...), true; got != want {
					t.Errorf("scan(...) = %v, want %v", got, want)
				}
				if got, want := reg1, test.reg1; got != want {
					t.Errorf("reg1 = %q, want %q", got, want)
				}
				if got, want := reg2, test.reg2; got != want {
					t.Errorf("reg2 = %q, want %q", got, want)
				}
				if got, want := val, test.val; got != want {
					t.Errorf("val = %v, want %v", got, want)
				}
				if got, want := offset, test.offset; got != want {
					t.Errorf("offset = %v, want %v", got, want)
				}
			})
		}
	})
}
