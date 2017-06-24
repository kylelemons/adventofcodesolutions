#!/usr/bin/python2.7

input = open('input.txt')

fwd = {}
rev = {}
subs = []

def xfrm(mol):
    short = [
        ("Y", ","),
        ("Rn", "{"),
        ("Ar", "}"),
        ("Al", "Z"),
        ("Ca", "X"),
        ("Mg", "W"),
        ("Si", "Q"),
        ("Th", "K"),
        ("Ti", "L"),
    ]
    for before, after in short:
        mol = mol.replace(before, after)
    return mol

while True:
    line = input.readline().strip()
    if line == "":
        break
    line = xfrm(line)
    pieces = line.split(" => ")
    fwd.setdefault(pieces[0], []).append(pieces[1])
    rev.setdefault(pieces[1], []).append(pieces[0])
    subs.append((pieces[0], pieces[1]))

molecule = xfrm(input.readline().strip())

def findall(haystack, needle):
    start = 0
    while True:
        start = haystack.find(needle, start)
        if start == -1:
            return
        yield start, start+len(needle)
        start += 1

def step1(molecule):
    next = set()
    for base, repls in fwd.iteritems():
        for start, end in findall(molecule, base):
            for repl in repls:
                replaced = molecule[:start] + repl + molecule[end:]
                next.add(replaced)
    return len(next)

def analyze():
    global constr

    rhs = rev.keys()
    lhs = fwd.keys()

    # Find unique constructions on the RHS
    pats = {}
    patv = set()
    for r in rhs:
        rem = r
        for l in lhs:
            rem = rem.replace(l, "")
        if rem:
            pats.setdefault(rem, []).append(r)
            patv.add(r)
    for pat, rs in pats.iteritems():
        print "pattern {}: {}".format(pat, rs)

    # Print out the replacements again, in two parts:
    print "Expansions:"
    for k, v in subs:
        if v in patv:
            continue
        print "  {} => {}".format(k, v)
    print "Constructions:"
    for pat, rs in pats.iteritems():
        print "  {}:".format(pat)
        for r in rs:
            for l in rev[r]:
                print "    {} => {}".format(l, r)

analyze()

from heapq import *
from random import *

def step2(molecule):
    print molecule

    seen = set()
    q = [(len(molecule), 0, 0, molecule)]
    min_len = len(molecule)
    timeout = 10000
    while len(q) > 0:
        timeout -= 1
        if timeout <= 0:
            return step2(molecule)
        _, steps, _, mol = heappop(q)
        if len(mol) == 1:
            return steps, mol
        if len(mol) < min_len:
            min_len = len(mol)
            print "{} after {} steps: {}".format(min_len, steps, mol)
        # We know what to do about runs of the recursive construction
        #  ... but don't shorten it beyond 3 yet
        while "XXX" in mol:
            steps += 1
            mol = mol.replace("XXX", "XX", 1)
        # Shorten with "brute force"
        for base, repls in rev.iteritems():
            for start, end in findall(mol, base):
                for repl in repls:
                    replaced = mol[:start] + repl + mol[end:]
                    if replaced in seen:
                        continue
                    seen.add(replaced)
                    heappush(q, (len(replaced), steps+1, randint(0, 1000), replaced))

print step2(xfrm("HPMg"))
print step2(xfrm("NThCaPTiMg"))
print step2(molecule)
