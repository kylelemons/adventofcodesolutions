#!/usr/bin/python3
from itertools import combinations, product
from functools import reduce
from dataclasses import dataclass
from re import findall, MULTILINE

@dataclass
class Row:
    lo: int
    hi: int
    letter: str
    password: str

def parse(input: str) -> list[Row]:
    return [Row(int(m[0]), int(m[1]), m[2], m[3]) for m in findall(r'^(\d+)-(\d+) (.): (.*)$', input, MULTILINE)]

def part1(input: str) -> int:
    rows = parse(input)
    print(rows[0])
    return sum([
        row.lo <= len(list(filter(lambda x: x == row.letter, row.password))) <= row.hi
        for row in rows
    ])

for (input, answer) in [
    ("""1-3 a: abcde
1-3 b: cdefg
2-9 c: ccccccccc""", 2),
    (open("input.txt").read(), 460),
]:
    got = part1(input)
    if got == answer:
        print("part1: PASS: {}".format(got))
    else:
        print("part1: FAIL: {} (want {})".format(got, answer))


def part2(input: str) -> int:
    rows = parse(input)
    print(rows[0])
    return sum([
        (row.password[row.lo-1] == row.letter) != (row.password[row.hi-1] == row.letter)
        for row in rows
    ])

for (input, answer) in [
    ("""1-3 a: abcde
1-3 b: cdefg
2-9 c: ccccccccc""", 1),
    (open("input.txt").read(), 251),
]:
    got = part2(input)
    if got == answer:
        print("part2: PASS: {}".format(got))
    else:
        print("part2: FAIL: {} (want {})".format(got, answer))
