input = 347991

def part1(val):
    if val == 1:
        return 0
    ring = 0
    while val > (2*ring+1)**2:
        ring += 1
    corner = (2*ring+1)**2
    val = corner - val
    side = 2*ring
    offset = val % side
    mid = side/2
    if offset >= mid:
        return offset-mid + ring
    if offset < mid:
        return mid-offset + ring

tests = [
    (1, 0),
    (12, 3),
    (23, 2),
    (1024, 31),
]
for (i, o) in tests:
    g = part1(i)
    if g == o:
        print('PASS', i, o)
    else:
        print('FAIL', i, g, 'WANT', o)

print ('input', part1(input))
