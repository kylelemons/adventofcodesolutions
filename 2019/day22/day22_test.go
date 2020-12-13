// Copyright 2019 Kyle Lemons
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
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kylelemons/adventofcodesolutions/advent"
)

type Helper struct {
	Deck  []int
	Lines *advent.Delimited
}

func parseInput(t *testing.T, in string, N int) *Helper {
	h := &Helper{
		Deck:  make([]int, N),
		Lines: advent.Lines(in),
	}
	for i := range h.Deck {
		h.Deck[i] = i
	}

	return h
}

func part1(t *testing.T, in string, N, target int) (ret int, deck []int) {
	h := parseInput(t, in, N)

	var (
		reverse = func() {
			for i, j := 0, len(h.Deck)-1; i < j; i, j = i+1, j-1 {
				h.Deck[i], h.Deck[j] = h.Deck[j], h.Deck[i]
			}
		}

		cutN = func(n int) {
			if n < 0 {
				n += len(h.Deck)
			}
			h.Deck = append(h.Deck[n:], h.Deck[:n]...)

		}

		incN = func(n int) {
			deck := make([]int, len(h.Deck))
			for i, v := range h.Deck {
				pos := (i * n) % len(deck)
				deck[pos] = v
			}
			h.Deck = deck

		}
	)

	h.Lines.Each(func(i int, line advent.Scanner) {
		var x int
		switch {
		case line.CanExtract(t, `deal into new stack`):
			reverse()
		case line.CanExtract(t, `cut (-?\d+)`, &x):
			cutN(x)
		case line.CanExtract(t, `deal with increment (\d+)`, &x):
			incN(x)
		default:
			t.Fatalf("Unrecognized line %q", line)
		}
	})

	if len(h.Deck) < 100 {
		t.Logf("Final: %v", h.Deck)
	}

	for i, v := range h.Deck {
		if v == target {
			return i, h.Deck
		}
	}
	return -1, nil
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name   string
		in     string
		N      int
		target int
		want   int
	}{
		{"part1 example 0", `deal with increment 7
deal into new stack
deal into new stack`, 10, 9, 3},
		{"part1 example 1", `cut 6
deal with increment 7
deal into new stack`, 10, 7, 2},
		{"part1 example 2", `deal with increment 7
deal with increment 9
cut -2`, 10, 7, 3},
		{"part1 example 2 x2", `deal with increment 7
deal with increment 9
cut -2
deal with increment 7
deal with increment 9
cut -2`, 10, 7, 1},
		{"part1 example 3", `deal into new stack
cut -2
deal with increment 7
cut 8
cut -4
deal with increment 7
cut 3
deal with increment 9
deal with increment 3
cut -1`, 10, 7, 6},
		{"part1 answer", advent.ReadFile(t, "input.txt"), 10007, 2019, 6417},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, _ := part1(t, test.in, test.N, test.target)
			if got, want := res, test.want; got != want {
				t.Errorf("part1(...)\n = %#v, want %#v", got, want)
			}
		})
	}
}

func part1linearized(t *testing.T, in string, N, times int64) func(int64) int64 {
	a, b := linearize(t, in, N)
	aExp, bExp := expLinear(a, b, times, N)
	return linearFunc(aExp, bExp, N)
}

func TestPart1Linearized(t *testing.T) {
	tests := []struct {
		name string
		in   string
		N    int
	}{
		// Dumb tests
		{"cut", `cut 6`, 11},
		{"incr", `deal with increment 7`, 11},
		{"deal", `deal into new stack`, 11},

		// Example tests
		{"part1 example 0", `deal with increment 7
deal into new stack
deal into new stack`, 11},
		{"part1 example 1", `cut 6
deal with increment 7
deal into new stack`, 11},
		{"part1 example 2", `deal with increment 7
deal with increment 9
cut -2`, 11},
		{"part1 example 2 x2", `deal with increment 7
deal with increment 9
cut -2
deal with increment 7
deal with increment 9
cut -2`, 11},
		{"part1 example 3", `deal into new stack
cut -2
deal with increment 7
cut 8
cut -4
deal with increment 7
cut 3
deal with increment 9
deal with increment 3
cut -1
`,
			11},
		{"part1 answer", advent.ReadFile(t, "input.txt"), 10007},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for _, times := range []int64{1, 2, 7, 13} {
				t.Run(fmt.Sprintf("x%d", times), func(t *testing.T) {
					input := strings.Repeat(test.in+"\n", int(times))
					_, baseline := part1(t, input, test.N, 0)

					part1linearized(t, test.in, int64(test.N), 1)
					f := part1linearized(t, test.in, int64(test.N), times)
					lin := make([]int, len(baseline))
					for i := range lin {
						lin[int(f(int64(i)))] = i
					}

					if diff := cmp.Diff(lin, baseline); diff != "" {
						t.Errorf("-linearized vs +iterative:\n%s", diff)
					}
				})
			}
		})
	}
}

func linearize(t *testing.T, in string, N int64) (stride, offset *big.Int) {
	// It turns out that it's possible to compute the output index for an input
	// index using a linear equation:
	//
	//   o = (a*i + b) % N
	//
	// Each of the operations is a version of this:
	//
	//   Operation | Equation       | Coefficients
	//   --------- | -------------- | ---------------
	//   deal      | o = -1*i + -1  | a = -1, b = -1
	//   cut_n     | o =  1*i + N-n | a =  1, b = N-n
	//   deal_n    | o =  n*i + 0   | a =  n, b = 0
	//
	// Luckily, it seems like linear operations like this can be combined before
	// the modulo is applied.
	//
	// The default coefficients are (a,b) = (1,0) of course.

	var (
		a = big.NewInt(1)
		b = big.NewInt(0)

		bigN     = big.NewInt(N)
		bigNsub1 = big.NewInt(N - 1)
		neg1     = big.NewInt(-1)
	)
	_, _ = bigNsub1, neg1

	advent.Lines(in).Each(func(i int, line advent.Scanner) {
		var n int64
		switch {
		case line.CanExtract(t, `deal into new stack`):
			// o = -1*(a*i + b) + -1
			//   = -a*i + (-b-1)
			a.Mul(a, neg1)
			b.Mul(b, neg1).Add(b, neg1)
		case line.CanExtract(t, `cut (-?\d+)`, &n):
			// o = 1*(a*i + b) + N-n
			//   = a*i + b+N-n
			b.Add(b, big.NewInt(N-n))
		case line.CanExtract(t, `deal with increment (\d+)`, &n):
			// o = n*(a*i + b) + 0
			//   = n*a*i + n*b
			bign := big.NewInt(n)
			a.Mul(a, bign)
			b.Mul(b, bign)
		default:
			t.Fatalf("Unrecognized line %q", line)
		}

		// Keep the numbers small
		a.Mod(a, bigN)
		b.Mod(b, bigN)

		// t.Logf("%s:", line)
		// t.Logf("  o = %d * i + %d", a, b)
	})
	return a, b
}

func invertLinear(a, b *big.Int, N int64) (aInv, bInv *big.Int) {
	// We want to derive a function, f_inv, such that:
	//
	//   f_inv(a*x + b) = x | mod m
	//
	// The resulting function should be another linear congruence:
	//
	//   f_inv(i) = a_inv*i + b_inv | mod m
	//
	// We can base this on two known fixed points, x = 1 and 0:
	//
	//   x = 1  ->  f_inv(a*1+b) = 1  ->  a_inv*a + a_inv*b + b_inv = 1 | mod m
	//   x = 0  ->  f_inv(a*0+b) = 0  ->            a_inv*b + b_inv = 0 | mod m
	//
	// Subtracting the two equations results in this:
	//
	//   a_inv*a = 1 | mod m
	//
	// This is the form required for using ModInverse:
	//
	//   a_inv = modinv(a, m)
	//
	// This can then be substituted back into one of our fixed points:
	//
	//   b_inv = -a_inv * b | mod m
	//
	// Once these are computed, the input index (and thus the numeric value) of
	// an output index is:
	//
	//   f_inv(o) = i = a_inv*o + b | mod m

	bigN := big.NewInt(N)
	aInv = new(big.Int).ModInverse(a, bigN)
	bInv = new(big.Int).Set(aInv)
	bInv.Neg(bInv).Mul(bInv, b).Mod(bInv, bigN)
	return aInv, bInv
}

func expLinear(a, b *big.Int, times, N int64) (aExp, bExp *big.Int) {
	if a.Int64() == 1 {
		// Special case: a == 1
		//
		//   y1 = a*x + b
		//      = x + b
		//
		//   y2 = (x + b) + b
		//      = x + 2*b
		//
		//   y3 = (x + 2*b) + b
		//      = x + 3*b
		//
		//   yn = x + n*b
		bExp = big.NewInt(times)
		bExp.Mul(bExp, b)
		return a, bExp
	}

	// Now, we need to consider the number of times that we repeat the operation
	// in sequence:
	//
	//   y1 = a*x + b
	//
	//   y2 = a*(a*x + b) + b
	//      = a^2*x + a*b + b
	//
	//   y3 = a*(a^2*x + a*b + b) + b
	//      = a^3*x + a^2*b + a*b + b
	//
	// The general pattern, then, looks like this:
	//
	//   yn = a^n*x + a^(n-1)*b + a^(n-2)*b + ... + b
	//      = a^n*x + b*(a^(n-1) + a^(n-2) + ... + a^0)
	//
	// Astute observers (i.e. not me) will recognize the infinite sum there:
	//
	//    Sum[a^k, {k, 0, n-1}] = (a^n - 1) / (a - 1)
	//
	// This allows us to write that as a finite expression:
	//
	//    yn = a^n*x + b*( (a^n - 1) / ( a - 1 ) )
	//
	// Unfortunately a^n can't be computed directly.
	// So, if we break it into separate terms:
	//
	//    yn = a^n*x + b * (a^n - 1) * (a-1)^-1 | mod m
	//
	// We can compute:
	//   a^n - 1  as exp(a,n,m)-1
	//   (a-1)^-1 as exp(a-1,m-2,m)
	//
	// This is because Fermat's little theorem implies: (for the second term)
	//   a^-1 = a^(n-2) | mod n

	if advent.GCD(int(N), int(times)) != 1 {
		panic(fmt.Sprintf("cannot use fast invmod for non-coprime n=%d m=%d", times, N))
	}

	var (
		bigN     = big.NewInt(N)
		bigTimes = big.NewInt(times)
		neg1     = big.NewInt(-1)
	)

	aExp = new(big.Int).Exp(a, big.NewInt(times), bigN)

	bExp = new(big.Int).Set(b)

	t1 := new(big.Int).Exp(a, bigTimes, bigN)
	t1.Add(t1, neg1)
	bExp.Mul(bExp, t1)

	t2 := new(big.Int).Set(a)
	t2.Add(t2, neg1)
	t2.Exp(t2, big.NewInt(N-2), bigN)
	bExp.Mul(bExp, t2)

	aExp.Mod(aExp, bigN)
	bExp.Mod(bExp, bigN)

	return
}

func linearFunc(a, b *big.Int, N int64) func(int64) int64 {
	v := new(big.Int)
	bigN := big.NewInt(N)
	return func(x int64) int64 {
		return v.SetInt64(x).Mul(a, v).Add(v, b).Mod(v, bigN).Int64()
	}
}

func part2(t *testing.T, in string, N, times int64) func(int64) int64 {
	a, b := linearize(t, in, N)
	aInv, bInv := invertLinear(a, b, N)
	aExp, bExp := expLinear(aInv, bInv, times, N)
	return linearFunc(aExp, bExp, N)
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name  string
		in    string
		N     int64
		at    int64
		times int64
		want  int64
	}{
		{"part1 answer", advent.ReadFile(t, "input.txt"), 10007, 6417, 1, 2019},
		{"part2 answer", advent.ReadFile(t, "input.txt"), 119315717514047, 2020, 101741582076661, 98461321956136},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := part2(t, test.in, test.N, test.times)
			if got, want := f(test.at), test.want; got != want {
				t.Errorf("part2(...): index %v = %#v, want %#v", test.at, got, want)
			}
		})
	}
}
