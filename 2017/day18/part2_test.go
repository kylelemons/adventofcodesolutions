package main_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"
)

func part2(input string, id int, out, in chan int) (ret int) {
	registers := map[string]int{
		"p": id,
	}

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
	var instr int
	for count := 0; instr >= 0 && instr < len(commands); count++ {
		if count > 100000 {
			panic(count)
		}

		command := strings.TrimSpace(commands[instr])
		if command == "" {
			instr++
			continue
		}

		args := strings.Fields(commands[instr])
		switch cmd, argv := args[0], args[1:]; cmd {
		case "set":
			registers[argv[0]] = get(argv[1])
		case "add":
			registers[argv[0]] += get(argv[1])
		case "mul":
			registers[argv[0]] *= get(argv[1])
		case "mod":
			registers[argv[0]] %= get(argv[1])
		case "snd":
			select {
			case out <- get(argv[0]):
				fmt.Println("send", get(argv[0]), "from", id)
				ret++
			case <-time.After(1 * time.Second):
				return
			}
		case "rcv":
			select {
			case registers[argv[0]] = <-in:
				fmt.Println("got", get(argv[0]), "at", id)

			case <-time.After(1 * time.Second):
				return
			}
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

func TestPart2(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part2 example", `
			snd 1
			snd 2
			snd p
			rcv a
			rcv b
			rcv c
			rcv d`, 3},
		{"part2", `set i 31
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
			jgz a -19`, 7112},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c01 := make(chan int, 10000)
			c10 := make(chan int, 10000)

			go part2(test.in, 0, c01, c10)

			if got, want := part2(test.in, 1, c10, c01), test.want; got != want {
				t.Errorf("part2(%#v) = %#v, want %#v", test.in, got, want)
			}
		})
	}
}
