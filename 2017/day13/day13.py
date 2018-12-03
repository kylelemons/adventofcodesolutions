input = open('input.txt').readlines()

scanners = [tuple(int(v) for v in line.strip().split(': ')) for line in input]
#scanners = [(0,3), (1,2), (4,4), (6,4)]

def try_at(t, scanners):
    states = [(depth, rng, (depth+t)%(rng*2-2)) for depth, rng in scanners]
    #print t, states
    severity = sum(depth*rng for depth, rng, pos in states if pos == 0)
    catches = sum(1 for depth, rng, pos in states if pos == 0)
    return (severity, catches)

for t in xrange(10000000):
    if try_at(t, scanners)[1] == 0:
        print t
        break
