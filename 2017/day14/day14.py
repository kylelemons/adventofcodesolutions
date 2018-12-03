
input = 'hwlqcszp'

def rev(knot, start, cnt):
    end = start+cnt-1
    while end > start:
        i, j = start%256, end%256
        knot[i], knot[j] = knot[j], knot[i]
        start += 1
        end -= 1

def knothash_round(knot, lengths, curpos, skipsize):
    for cnt in lengths:
        rev(knot, curpos, cnt)
        curpos += cnt
        curpos += skipsize
        skipsize += 1
    return (curpos, skipsize)

HEXCHARS = '0123456789abcdef'

def knothash(bytes):
    knot = list(range(256))
    lengths = list(ord(b) for b in bytes)
    lengths += [17, 31, 73, 47, 23]

    curpos, skipsize = knot[0], 0
    for i in xrange(64):
        curpos, skipsize = knothash_round(knot, lengths, curpos, skipsize)

    dense = [reduce(lambda x, y: x^y, knot[b*16:][:16]) for b in range(16)]

    kh_hex = ''.join(HEXCHARS[v/16]+HEXCHARS[v%16] for v in dense)
    kh_bin = ''.join(bin(v)[2:].zfill(8) for v in dense)
    return kh_hex, kh_bin

def get_disk(input):
    return [[['.', '#'][int(v)] for v in knothash('{}-{}'.format(input, i))[1]] for i in xrange(128)]

# print 'Filled: {}'.format(sum(sum(1 if v == '#' else 0 for v in r) for r in get_disk(input)))

def flood_fill(disk, r, c, ch):
    if r < 0 or r >= len(disk):
        return
    if c < 0 or c >= len(disk[r]):
        return
    if disk[r][c] != '#':
        return

    dirs = [(-1, 0), (0, 1), (1, 0), (0, -1)] # urdl
    disk[r][c] = ch
    for dr, dc in dirs:
        flood_fill(disk, r+dr, c+dc, ch)

def count_regions(input):
    disk = get_disk(input)
    regions = 0
    for r in range(len(disk)):
        for c in range(len(disk[r])):
            if disk[r][c] == '#':
                regions += 1
                flood_fill(disk, r, c, str(regions%10))
    for row in disk:
        print input, ''.join(row)
    return regions

print count_regions('flqrgnkx')
print count_regions(input)
