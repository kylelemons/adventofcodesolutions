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

package coords

// In2D uses the corod to index into a 2D array of bytes.
func (c Coord) In2D(matrix [][]byte) byte {
	return matrix[c.R()][c.C()]
}

// InBounds2D uses the corod to index into a 2D array of bytes.
//
// If the coordinate is not in bounds, false will be returned.
func (c Coord) InBounds2D(matrix [][]byte) (byte, bool) {
	row, col := c.R(), c.C()
	if row < 0 || col < 0 || row >= len(matrix) || col >= len(matrix[row]) {
		return 0, false
	}
	return matrix[row][col], true
}
