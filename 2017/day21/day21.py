rules = [tuple(line.strip().split(' => ')) for line in open('input.txt') if line != '']

#rules = [('../.#', '##./#../...'), ('.#./..#/###', '#..#/..../..../#..#')]

xfrm = {tuple(i.split('/')): tuple(o.split('/')) for i, o in rules}

pattern = [
    '.#.',
    '..#',
    '###',
]

def flips(chunk):
    # No change
    yield tuple(chunk)

    # Vertical flip
    yield tuple(reversed(chunk))

    # Horizontal flip
    hflip = tuple(''.join(reversed(r)) for r in chunk)
    yield hflip

    # Double flip
    yield tuple(reversed(hflip))

def rotations(chunk):
    n = len(chunk)

    # No change
    yield tuple(chunk)

    # 90-degree rotation increments
    rot90  = tuple(''.join(chunk[n-c-1][r]  for c in range(n)) for r in range(n))
    yield rot90

    rot180 = tuple(''.join(rot90[n-c-1][r]  for c in range(n)) for r in range(n))
    yield rot180

    rot270 = tuple(''.join(rot180[n-c-1][r] for c in range(n)) for r in range(n))
    yield rot270

def variations(chunk):
    for flipped in flips(chunk):
        for rotated in rotations(flipped):
            yield rotated

# Part 1
def expand(pattern, iterations):
    # for line in pattern:
    #     print iterations, line

    if iterations <= 0:
        return pattern

    div = 2 if len(pattern) % 2 == 0 else 3

    out = []
    for r in range(0, len(pattern), div):
        i = len(out)
        out += [''] * (div+1)
        for c in range(0, len(pattern[r]), div):
            chunk = tuple(row[c:][:div] for row in pattern[r:][:div])
            found = False
            for key in variations(chunk):
                if key not in xfrm:
                    continue
                for ii, orow in enumerate(xfrm[key]):
                    found = True
                    out[i+ii] += orow
                break
            if not found:
                print '{} Failed to find any of the following keys:'.format(iterations)
                for key in sorted(['/'.join(key) for key in variations(chunk)]):
                    print '{}   {}'.format(iterations, key)
                for i, o in sorted(xfrm.iteritems()):
                    if len(i) != 3:
                        continue
                    i = '/'.join(i)
                    if i.count('#') != 5:
                        continue
                    if i[5] != '.':
                        continue
                    print 'candidate: {} -> {}'.format(i, '/'.join(o))
                return pattern
    
    return expand(out, iterations-1)

print 'Lights: {}'.format(sum(sum([1 if c == '#' else 0 for c in line])for line in expand(pattern, 18)))

'''
5   ###/#../.#.
5   ###/#../.#.
5   ###/..#/.#.
5   ###/..#/.#.
5   ##./#.#/#..
5   ##./#.#/#..
5   #../#.#/##.
5   #../#.#/##.
5   .##/#.#/..#
5   .##/#.#/..#
5   .#./#../###
5   .#./#../###
5   .#./..#/###
5   .#./..#/###
5   ..#/#.#/.##
5   ..#/#.#/.##
'''
