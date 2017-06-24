'''adventofcode 2015 day16'''

import re

CLUES = {k: int(v) for k, v in re.findall(r'(\w+): (\d+)', open("clues.txt").read())}

SUES = {s: {k: int(v) for k, v in re.findall(r'(\w+): (\d+)', c)}
        for s, c in re.findall(r'Sue (\d+): (.*)', open("input.txt").read())}

for sue, info in SUES.iteritems():
    matches = 0
    for k, v in info.iteritems():
        if k == "cats" or k == "trees":
            if CLUES[k] < v:
                matches += 1
        elif k == "goldfish" or k == "pomeranians":
            if CLUES[k] > v:
                matches += 1
        elif CLUES[k] == v:
            matches += 1
    if matches == len(info):
        print sue, info
