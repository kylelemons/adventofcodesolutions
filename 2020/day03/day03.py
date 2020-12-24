#!/usr/bin/python3
from itertools import combinations, product
from functools import reduce
from dataclasses import dataclass
from re import findall, MULTILINE

@dataclass
class Row:
    map: list[str]

@dataclass
class Point:
    r: int
    c: int

    def add(self, p):
        return Point(self.r + p.r, self.c + p.c)

def parse(input: str) -> list[Row]:
    return Row(input.splitlines())

def part1(input: str) -> int:
    inp = parse(input)
    loc = Point(0,0)
    delta = Point(1,3)
    count = 0
    while loc.r < len(inp.map):
        if inp.map[loc.r][loc.c % len(inp.map[loc.r])] == '#':
            count += 1
        loc = loc.add(delta)
    return count

for (input, answer) in [
    ("""..##.......
#...#...#..
.#....#..#.
..#.#...#.#
.#...##..#.
..#.##.....
.#.#.#....#
.#........#
#.##...#...
#...##....#
.#..#...#.#""", 7),
    (open("input.txt").read(), 211),
]:
    got = part1(input)
    if got == answer:
        print("part1: PASS: {}".format(got))
    else:
        print("part1: FAIL: {} (want {})".format(got, answer))

def part2(input: str) -> int:
    def do(delta: Point) -> int:
        inp = parse(input)
        loc = Point(0,0)
        count = 0
        while loc.r < len(inp.map):
            if inp.map[loc.r][loc.c % len(inp.map[loc.r])] == '#':
                count += 1
            loc = loc.add(delta)
        return count
    return reduce(lambda x,y: x*y, [
        do(delta)
        for delta in [
            Point(1,1),
            Point(1,3),
            Point(1,5),
            Point(1,7),
            Point(2,1),
        ]
    ])

for (input, answer) in [
    ("""..##.......
#...#...#..
.#....#..#.
..#.#...#.#
.#...##..#.
..#.##.....
.#.#.#....#
.#........#
#.##...#...
#...##....#
.#..#...#.#""", 336),
    (open("input.txt").read(), 3584591857),
]:
    got = part2(input)
    if got == answer:
        print("part2: PASS: {}".format(got))
    else:
        print("part2: FAIL: {} (want {})".format(got, answer))

