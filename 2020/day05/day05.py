#!/usr/bin/python3
from itertools import combinations, product
from functools import reduce
from dataclasses import dataclass
from re import findall, search, MULTILINE
from typing import Union

@dataclass
class Record:
    item: dict[str, str]

def parse(input: str) -> list[str]:
    return input.splitlines()

def part1(raw: str) -> int:
    input = parse(raw)
    def id(r: str) -> int:
        bin = r
        for rep in [('F', '0'), ('B', '1'), ('L', '0'), ('R', '1')]:
            bin = bin.replace(*rep)
        return int(bin, 2)
    return max([id(r) for r in input])

for (input, answer) in [
    ("BFFFBBFRRR", 567),
    ("FFFBBBFRRR", 119),
    ("BBFFBBFRLL", 820),
    (open("input.txt").read(), 930),
]:
    got = part1(input)
    if got == answer:
        print("part1: PASS: {}".format(got))
    else:
        print("part1: FAIL: {} (want {})".format(got, answer))


def part2(raw: str) -> int:
    input = parse(raw)
    def id(r: str) -> int:
        bin = r
        for rep in [('F', '0'), ('B', '1'), ('L', '0'), ('R', '1')]:
            bin = bin.replace(*rep)
        return int(bin, 2)
    tickets = sorted([id(r) for r in input])
    for i, id in enumerate(tickets):
        if tickets[i+1] != id+1:
            return id+1
    return None

for (input, answer) in [
    (open("input.txt").read(), 515),
]:
    got = part2(input)
    if got == answer:
        print("part2: PASS: {}".format(got))
    else:
        print("part2: FAIL: {} (want {})".format(got, answer))
