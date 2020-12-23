#!/usr/bin/python3
from itertools import combinations, product
from functools import reduce

def parse(input: str) -> list[int]:
    return [int(v) for v in input.splitlines()]

def part1(input: str) -> int:
    values = parse(input)
    for pair in combinations(values, 2):
        if sum(pair) == 2020:
            return reduce(lambda x,y: x*y, pair)
    return None

for (input, answer) in [
    ("""1721
979
366
299
675
1456""", 514579),
    (open("input.txt").read(), 355875),
]:
    got = part1(input)
    if got == answer:
        print("part1: PASS: {}".format(got))
    else:
        print("part1: FAIL: {} (want {})".format(got, answer))

def part2(input: str) -> int:
    values = parse(input)
    for trio in combinations(values, 3):
        if sum(trio) == 2020:
            return reduce(lambda x,y: x*y, trio)
    return None

for (input, answer) in [
    ("""1721
979
366
299
675
1456""", 241861950),
    (open("input.txt").read(), 140379120),
]:
    got = part2(input)
    if got == answer:
        print("part2: PASS: {}".format(got))
    else:
        print("part2: FAIL: {} (want {})".format(got, answer))
