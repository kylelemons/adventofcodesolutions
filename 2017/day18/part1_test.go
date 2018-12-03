package main_test

import (
	"strconv"
	"strings"
	"testing"
)

func part1(input string) (ret int) {
	registers := map[string]int{}

	get := func(s string) int {
		if strings.IndexAny(s, "0123456789") != -1 {
			v, err := strconv.Atoi(s)
			if err != nil {
				panic(err)
			}
			return v
		}
		return registers[s]
	}

	commands := strings.Split(input, "\n")
	var instr, lastPlayed int
	for count := 0; instr >= 0 && instr < len(commands); count++ {
		if count > 100000 {
			panic(count)
		}

		command := strings.TrimSpace(commands[instr])
		if command == "" {
			instr++
			continue
		}
		//fmt.Println(command, registers, lastPlayed)

		args := strings.Fields(commands[instr])
		switch cmd, argv := args[0], args[1:]; cmd {
		case "snd":
			lastPlayed = get(argv[0])
		case "set":
			registers[argv[0]] = get(argv[1])
		case "add":
			registers[argv[0]] += get(argv[1])
		case "mul":
			registers[argv[0]] *= get(argv[1])
		case "mod":
			registers[argv[0]] %= get(argv[1])
		case "rcv":
			if get(argv[0]) == 0 {
				break
			}
			return lastPlayed
			// registers[argv[0]] = lastPlayed
		case "jgz":
			if get(argv[0]) <= 0 {
				break
			}
			instr += get(argv[1])
			continue
		default:
			panic(command)
		}
		instr++
	}
	panic("exit")
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part1 example", `
			set a 1
			add a 2
			mul a a
			mod a 5
			snd a
			set a 0
			rcv a
			jgz a -1
			set a 1
			jgz a -2`, 4},
		{"part1", `set i 31
				set a 1
				mul p 17
				jgz p p
				mul a 2
				add i -1
				jgz i -2
				add a -1
				set i 127
				set p 680
				mul p 8505
				mod p a
				mul p 129749
				add p 12345
				mod p a
				set b p
				mod b 10000
				snd b
				add i -1
				jgz i -9
				jgz a 3
				rcv b
				jgz b -1
				set f 0
				set i 126
				rcv a
				rcv b
				set p a
				mul p -1
				add p b
				jgz p 4
				snd a
				set a b
				jgz 1 3
				snd b
				set f 1
				add i -1
				jgz i -11
				snd a
				jgz f -16
				jgz a -19`, 3188},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part1(test.in), test.want; got != want {
				t.Errorf("part1(%#v) = %#v, want %#v", test.in, got, want)
			}
		})
	}
}
