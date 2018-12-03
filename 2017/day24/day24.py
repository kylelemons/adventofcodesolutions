input = [line.split('/') for line in open('input.txt')]
comps = [tuple(int(v) for v in line) for line in input]
comps = [(v[0], v[1]) if v[0] > v[1] else (v[1], v[0]) for v in comps]

def strength(comp):
    return comp[0] + comp[1]

from heapq import *

def part1(comps):
    comps = set(comps)
    
    # (sort, cur, last, available)
    q = [(0, 0, 0, comps)]
    ret = 0
    while len(q) > 0:
        sort, cur, last, available = heappop(q)
        if cur > ret:
            ret = cur
            print ret

        found = 0
        mx = cur
        for comp in available:
            if comp[0] != last and comp[1] != last:
                continue
            found += 1

            heappush(q, (
                sort - strength(comp),
                cur + strength(comp),
                comp[1] if comp[0] == last else comp[0],
                available - set([comp]),
            ))
    return ret

# print part1(comps)

def part2(comps):
    comps = set(comps)
    
    # (sort, cur, cnt, last, available)
    q = [((0,0), 0, 0, 0, comps)]
    ret = (0, 0)
    while len(q) > 0:
        sort, cur, cnt, last, available = heappop(q)
        if sort > ret:
            ret = sort
            # print ret

        found = 0
        mx = cur
        for comp in available:
            if comp[0] != last and comp[1] != last:
                continue
            found += 1

            heappush(q, (
                (cnt + 1, cur + strength(comp)),
                cur + strength(comp),
                cnt + 1,
                comp[1] if comp[0] == last else comp[0],
                available - set([comp]),
            ))
    return ret

print part2(comps)
