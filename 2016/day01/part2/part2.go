package main

import (
	"fmt"
	"strings"
)

const (
	North = iota
	East
	South
	West
	MaxDir
)

var dd = [][2]int{
	North: {0, -1},
	East:  {1, 0},
	South: {0, 1},
	West:  {-1, 0},
}

func distance(input string) int {
	directions := strings.Split(input, ", ")
	x, y, curr := 0, 0, North
	visited := map[[2]int]int{{0, 0}: 1}
	for _, d := range directions {
		var turn rune
		var dist int
		if _, err := fmt.Sscanf(d, "%c%d", &turn, &dist); err != nil {
			panic(fmt.Sprintf("failed to parse %q: %s", d, err))
		}
		switch turn {
		case 'L':
			curr--
		case 'R':
			curr++
		default:
			panic("Unknown turn " + string(turn))
		}
		for curr < 0 {
			curr += MaxDir
		}
		for curr >= MaxDir {
			curr -= MaxDir
		}
		for i := 0; i < dist; i++ {
			x += dd[curr][0]
			y += dd[curr][1]
			loc := [2]int{x, y}
			fmt.Println(d, loc)
			if visited[loc] > 0 {
				goto finished
			}
			visited[loc]++
		}
	}

finished:
	if x < 0 {
		x = -x
	}
	if y < 0 {
		y = -y
	}
	return x + y
}

func main() {
	for _, input := range []string{
		"L2, R3",
		"R2, R2, R2",
		"R5, L5, R5, R3",
		"R8, R4, R4, R8",
		"R5, L2, L1, R1, R3, R3, L3, R3, R4, L2, R4, L4, R4, R3, L2, L1, L1, R2, R4, R4, L4, R3, L2, R1, L4, R1, R3, L5, L4, L5, R3, L3, L1, L1, R4, R2, R2, L1, L4, R191, R5, L2, R46, R3, L1, R74, L2, R2, R187, R3, R4, R1, L4, L4, L2, R4, L5, R4, R3, L2, L1, R3, R3, R3, R1, R1, L4, R4, R1, R5, R2, R1, R3, L4, L2, L2, R1, L3, R1, R3, L5, L3, R5, R3, R4, L1, R3, R2, R1, R2, L4, L1, L1, R3, L3, R4, L2, L4, L5, L5, L4, R2, R5, L4, R4, L2, R3, L4, L3, L5, R5, L4, L2, R3, R5, R5, L1, L4, R3, L1, R2, L5, L1, R4, L1, R5, R1, L4, L4, L4, R4, R3, L5, R1, L3, R4, R3, L2, L1, R1, R2, R2, R2, L1, L1, L2, L5, L3, L1",
	} {
		fmt.Printf("f(%q)\n = %d\n", input, distance(input))
	}
}
