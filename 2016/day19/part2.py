#!/usr/bin/python2.7

import math

def white_elephant(elves, debug=True):
    has_present = [i+1 for i in xrange(elves)]
    while len(has_present) > 1:
        i = 0
        while i < len(has_present):
            victim = i + int(len(has_present)/2)
            victim %= len(has_present)
            del has_present[victim]
            if victim > i:
                i += 1
    return has_present[0]

# Thanks, Reddit!
#   https://www.reddit.com/r/adventofcode/comments/5j4lp1/2016_day_19_solutions/dbdf4up/
def opt_elephant(elves):
    base = int(math.log(elves, 3))
    end = int(math.pow(3, base))
    if elves == end:
        return end
    offset = elves - end
    if offset <= end:
        return offset
    twos = offset - end
    return end + twos*2

for i in xrange(81):
    print i+1, white_elephant(i+1, i < 6), opt_elephant(i+1)

print opt_elephant(3004953)
