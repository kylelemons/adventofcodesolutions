package aocday

import (
	"io/ioutil"
	"strings"
	"testing"
)

func part1(t *testing.T, in string) string {
	lines := strings.Split(in, "\n")

	r, c := 0, strings.Index(lines[0], "|")

	deltas := [][2]int{
		// dr, dc
		{1, 0},
		{0, 1},
		{-1, 0},
		{0, -1},
	}
	dir := 0

	var found string
	for i := 0; i < 1e6; i++ {
		ch := lines[r][c]
		switch ch {
		default:
			found += string(ch)
			fallthrough
		case '|', '-':
			r += deltas[dir][0]
			c += deltas[dir][1]
		case '+':
			ld := (dir + 1) % len(deltas)
			l := deltas[ld]
			r2, c2 := r+l[0], c+l[1]
			if lines[r2][c2] != ' ' {
				dir, r, c = ld, r2, c2
				continue
			}

			rd := (len(deltas) + dir - 1) % len(deltas)
			rt := deltas[rd]
			r2, c2 = r+rt[0], c+rt[1]
			if lines[r2][c2] != ' ' {
				dir, r, c = rd, r2, c2
				continue
			}
		case ' ':
			return found
		}
	}
	panic("timed out")
}

func part2(t *testing.T, in string) int {
	lines := strings.Split(in, "\n")

	r, c := 0, strings.Index(lines[0], "|")

	deltas := [][2]int{
		// dr, dc
		{1, 0},
		{0, 1},
		{-1, 0},
		{0, -1},
	}
	dir := 0

	var found string
	for i := 0; i < 1e6; i++ {
		ch := lines[r][c]
		switch ch {
		default:
			found += string(ch)
			fallthrough
		case '|', '-':
			r += deltas[dir][0]
			c += deltas[dir][1]
		case '+':
			ld := (dir + 1) % len(deltas)
			l := deltas[ld]
			r2, c2 := r+l[0], c+l[1]
			if lines[r2][c2] != ' ' {
				dir, r, c = ld, r2, c2
				continue
			}

			rd := (len(deltas) + dir - 1) % len(deltas)
			rt := deltas[rd]
			r2, c2 = r+rt[0], c+rt[1]
			if lines[r2][c2] != ' ' {
				dir, r, c = rd, r2, c2
				continue
			}
		case ' ':
			return i
		}
	}
	panic("timed out")
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"part1 example", `     |          
     |  +--+    
     A  |  C    
 F---|----E|--+ 
     |  |  |  D 
     +B-+  +--+ `, "ABCDEF"},
		{"part1", read(t, "input.txt"), "LIWQYKMRP"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part1(t, test.in), test.want; got != want {
				t.Errorf("part1(%#v) = %#v, want %#v", test.in, got, want)
			}
		})
	}

	t.Log(part2(t, read(t, "input.txt")))
}

func read(t *testing.T, filename string) string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatalf("failed to read %q: %s", filename, err)
	}
	return string(data)
}
