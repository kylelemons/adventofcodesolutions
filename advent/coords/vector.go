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

// Vector is a synonym for Coord for when it makes sense for readability.
type Vector = Coord

// Coordinate constants.
var (
	// Cardinal (4-way) directions:
	North = XY(0, -1)
	East  = XY(+1, 0)
	South = XY(0, +1)
	West  = XY(-1, 0)

	Cardinals = []Coord{North, East, South, West}

	// Compass (8-way) directions:
	NorthEast = North.Add(East)
	NorthWest = North.Add(West)
	SouthEast = South.Add(East)
	SouthWest = South.Add(West)

	Compass = []Coord{
		North, NorthEast, East, SouthEast,
		South, SouthWest, West, NorthWest,
	}
)

// Right returns the vector v rotated clockwise 90 degrees around the origin.
func (v Vector) Right() Vector { return Vector{x: -v.y, y: v.x} }

// Left returns the vector v rotated counterclockwise 90 degrees around the origin.
func (v Vector) Left() Vector { return Vector{x: v.y, y: -v.x} }

// Reverse returns the negative of the given vector.
func (v Vector) Reverse() Vector { return v.Scale(-1) }

// Scale returns the vector scaled by the given magnitude.
func (v Vector) Scale(by int) Vector { return Vector{x: v.x * by, y: v.y * by} }
