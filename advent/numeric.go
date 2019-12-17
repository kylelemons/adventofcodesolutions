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

package advent

// GCD returns the greatest common divisor of x and y.
func GCD(x, y int) int {
	for y != 0 {
		x, y = y, x%y
	}
	return x
}

// factorize returns the prime factorization of value as the keys of a map whose
// values are the number of times the factor multiplies into n.
func factorize(value int) map[int]int {
	factors := make(map[int]int)
	for divisor := 2; value > 1; divisor++ {
		for value%divisor == 0 {
			factors[divisor]++
			value /= divisor
		}
	}
	return factors
}

// LCM returns the least common multiple of the given values.
func LCM(values ...int) int {
	if len(values) == 0 {
		panic("LCM requires at least one value")
	}

	max := factorize(values[0]) // max[factor] = highest power
	for _, value := range values[1:] {
		for factor, count := range factorize(value) {
			if current := max[factor]; count > current {
				max[factor] = count
			}
		}
	}

	lcm := 1
	for factor, count := range max {
		for i := 0; i < count; i++ {
			lcm *= factor
		}
	}
	return lcm
}
