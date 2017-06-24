#!/usr/bin/python2.7

from copy import *
from itertools import *

# Input:
floors = [
# The first floor contains a promethium generator and a promethium-compatible microchip.
    [("PR", "G"), ("PR", "M")],
# The second floor contains a cobalt generator, a curium generator, a ruthenium generator, and a plutonium generator.
    [("CO", "G"), ("CU", "G"), ("RU", "G"), ("PL", "G")],
# The third floor contains a cobalt-compatible microchip, a curium-compatible microchip, a ruthenium-compatible microchip, and a plutonium-compatible microchip.
    [("CO", "M"), ("CU", "M"), ("RU", "M"), ("PL", "M")],
# The fourth floor contains nothing relevant.
    [],
]

# Part 2:
floors[0] += [tuple(t) for t in product(["EL", "DI"], ["G", "M"])]

# Example:
example = [
    [("H", "M"), ("L", "M")],
    [("H", "G")],
    [("L", "G")],
    [],
]

def will_fry(floor):
    """
    input:
      floor - list of (ele, typ) tuples of source floor
    returns:
      True if the floor contents will fry a microchip
    """
    gs = set(obj[0] for obj in floor if obj[1] == "G")
    ms = set(obj[0] for obj in floor if obj[1] == "M")
    ms -= gs
    return len(ms) > 0 and len(gs) > 0

def carries(floor, onto):
    """
    input:
      floor - list of (ele, typ) tuples of source floor
      onto  - list of (ele, typ) tuples of dest floor
    """
    carry2 = combinations(floor, 2)
    carry1 = combinations(floor, 1)
    for objs in chain(carry2, carry1):
        s = set(objs)
        if not will_fry(onto | s) and not will_fry(floor - s):
            yield s

def locate_all(floors):
    """
    input:
      floors - list of lists of (ele, typ) tuples
    outupt:
      yields (element, type, floor) tuples
    """
    for f, floor in enumerate(floors):
        for element, typ in floor:
            yield (element, typ, f)

def key(floors, on_floor):
    """
    input:
      floors   - list of lists of (ele, typ) tuples
      on_floor - the elevator's floor
    output:
      a string representing the state in an element-independent way
    """
    locs = sorted((e, int(t == 'G'), f) for e, t, f in locate_all(floors))
    pairs = sorted(tuple(f for e, t, f in g) for k, g in groupby(locs, key=lambda x: x[0]))
    return str((pairs, on_floor))

def goodness(floors):
    return sum(f*f*len(objs) for f, objs in enumerate(floors))

def floorstr(floor):
    return ' '.join(sorted(''.join(x) for x in floor))

def pretty(floors, on_floor):
    return '\n'.join([
        '     F{} {:1} {}'.format(f+1, 'E' if f == on_floor else ' ', floorstr(floor))
        for f, floor in reversed(list(enumerate(floors)))
    ])

def min_steps(floors, maxdepth = 100):
    # Turn the floors into sets
    floors = [set(f) for f in floors]

    # Keep track of states we've seen
    seen = set()

    # Heap queue
    # (steps, on_floor, floors)
    q = []
    q.append(([], 0, floors))

    # Figure out the elements in use
    elements = sorted(set([obj[0] for floor in floors for obj in floor]))
    print "STARTING STATE:"
    print pretty(floors, 0)
    print ""

    while len(q) > 0:
        steps, on_floor, floors = q.pop(0) # bfs

        # Recursion depth limit
        if len(steps) > maxdepth:
            continue

        floor = floors[on_floor]
        #print "step", len(steps), "on", on_floor, "with", sorted(floor), "key", key(floors, on_floor)

        if len(floors[3]) >= 2*len(elements): # all objects
            print "SOLUTION: {} steps".format(len(steps))
            for s, st in enumerate(steps):
                print "{:3} {}".format(s+1, st)
            return len(steps)

        # EXPERIMENTAL: don't move below the lowest occupied floor
        min_floor = min(f for f, floor in enumerate(floors) if len(floor) > 0)

        # Figure out what floors we can go to
        nexts = []
        if on_floor > 0:
            nexts.append(on_floor-1)
        if on_floor < len(floors)-1:
            nexts.append(on_floor+1)

        for to_floor in nexts:
            if to_floor < min_floor:
                continue
            onto = floors[to_floor]
            for carry in carries(floor, onto):
                nfloors = deepcopy(floors)
                nfloors[on_floor] -= carry
                nfloors[to_floor] |= carry
                nsteps = steps + ["carry [{}] from {} to {}\n{}\n".format(floorstr(carry), on_floor, to_floor, pretty(nfloors, to_floor))]

                # Prune seen states
                k = key(nfloors, to_floor)
                if k in seen:
                    continue
                seen.add(k)

                #print "  carry", sorted(carry), "from", on_floor, "to", to_floor, "which has", sorted(onto)
                q.append((nsteps, to_floor, nfloors))

print min_steps(example, 11)
print min_steps(floors, 1000)

# 29 is incorrect
# 37 is too high
