#!/usr/bin/python3
from itertools import combinations, product
from functools import reduce
from dataclasses import dataclass
from re import findall, search, MULTILINE
from typing import Union

@dataclass
class Record:
    item: dict[str, str]

def parse(input: str) -> list[Record]:
    return [
        Record(
            {
                m[0]: m[1]
                for m in findall(r'(...):(\S+)', line, MULTILINE)
            }
        ) for line in input.split('\n\n')
    ]

def part1(raw: str) -> int:
    input = parse(raw)
    def valid(r: Record) -> bool:
        for req in ['byr', 'iyr', 'eyr', 'hgt', 'hcl', 'ecl', 'pid']:
            if req not in r.item:
                return False
        return True
    return sum([valid(r) for r in input])

for (input, answer) in [
    ("""ecl:gry pid:860033327 eyr:2020 hcl:#fffffd
byr:1937 iyr:2017 cid:147 hgt:183cm

iyr:2013 ecl:amb cid:350 eyr:2023 pid:028048884
hcl:#cfa07d byr:1929

hcl:#ae17e1 iyr:2013
eyr:2024
ecl:brn pid:760753108 byr:1931
hgt:179cm

hcl:#cfa07d eyr:2025 pid:166559648
iyr:2011 ecl:brn hgt:59in""", 2),
    (open("input.txt").read(), 239),
]:
    got = part1(input)
    if got == answer:
        print("part1: PASS: {}".format(got))
    else:
        print("part1: FAIL: {} (want {})".format(got, answer))


def part2(raw: str) -> int:
    input = parse(raw)
    def valid(r: Record) -> bool:
        for req in ['byr', 'iyr', 'eyr', 'hgt', 'hcl', 'ecl', 'pid']:
            if req not in r.item:
                return False
        byr = int(r.item['byr'])
        iyr = int(r.item['iyr'])
        eyr = int(r.item['eyr'])
        if not (1920 <= byr <= 2002):
            return False
        if not (2010 <= iyr <= 2020):
            return False
        if not (2020 <= eyr <= 2030):
            return False

        hm = search(r'(\d+)(in|cm)', r.item['hgt'])
        if not hm:
            return False
        hgt, unit = hm.groups()
        if unit == 'cm' and not (150 <= int(hgt) <= 193):
            return False
        if unit == 'in' and not (59 <= int(hgt) <= 76):
            return False

        if not search(r'^#[a-f0-9]{6}$', r.item['hcl']):
            return False

        if r.item['ecl'] not in ('amb', 'blu', 'brn', 'gry', 'grn', 'hzl', 'oth'):
            return False

        if not search(r'^\d{9}$', r.item['pid']):
            return False
        return True
    return sum([valid(r) for r in input])

for (input, answer) in [
    ("""eyr:1972 cid:100
hcl:#18171d ecl:amb hgt:170 pid:186cm iyr:2018 byr:1926

iyr:2019
hcl:#602927 eyr:1967 hgt:170cm
ecl:grn pid:012533040 byr:1946

hcl:dab227 iyr:2012
ecl:brn hgt:182cm pid:021572410 eyr:2020 byr:1992 cid:277

hgt:59cm ecl:zzz
eyr:2038 hcl:74454a iyr:2023
pid:3556412378 byr:2007""", 0),
    ("""pid:087499704 hgt:74in ecl:grn iyr:2012 eyr:2030 byr:1980
hcl:#623a2f

eyr:2029 ecl:blu cid:129 byr:1989
iyr:2014 pid:896056539 hcl:#a97842 hgt:165cm

hcl:#888785
hgt:164cm byr:2001 iyr:2015 cid:88
pid:545766238 ecl:hzl
eyr:2022

iyr:2010 hgt:158cm hcl:#b6652a ecl:blu byr:1944 eyr:2021 pid:093154719""", 4),
    (open("input.txt").read(), 188),
]:
    got = part2(input)
    if got == answer:
        print("part2: PASS: {}".format(got))
    else:
        print("part2: FAIL: {} (want {})".format(got, answer))
