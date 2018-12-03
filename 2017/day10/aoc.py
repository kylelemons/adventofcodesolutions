input = [31,2,85,1,80,109,35,63,98,255,0,13,105,254,128,33]

def part1():
    knot = list(range(256))

    curpos = knot[0]
    skipsize = 0
    lengths = input

    def rev(start, cnt):
        end = start+cnt-1
        while end > start:
            i, j = start%256, end%256
            knot[i], knot[j] = knot[j], knot[i]
            start += 1
            end -= 1
    
    for cnt in lengths:
        rev(curpos, cnt)
        curpos += cnt
        curpos += skipsize
        skipsize += 1

    return knot[0] * knot[1]

print part1()

input = [ord(r) for r in '31,2,85,1,80,109,35,63,98,255,0,13,105,254,128,33']
input += [17, 31, 73, 47, 23]

def part2():
    knot = list(range(256))
    lengths = input

    # hax for nonlocal
    curpos = [knot[0]]
    skipsize = [0]

    def rev(start, cnt):
        end = start+cnt-1
        while end > start:
            i, j = start%256, end%256
            knot[i], knot[j] = knot[j], knot[i]
            start += 1
            end -= 1
    
    def round():
        for cnt in lengths:
            rev(curpos[0], cnt)
            curpos[0] += cnt
            curpos[0] += skipsize[0]
            skipsize[0] += 1

    for i in range(64):
        round()

    dense = [reduce(lambda x,y: x^y, knot[b*16:][:16]) for b in range(16)]

    hx = '0123456789abcdef'
    return ''.join(hx[v/16]+hx[v%16] for v in dense)

print part2()
