#!/usr/bin/python2.7

def white_elephant(elves):
    has_present = [i+1 for i in xrange(elves)]
    while len(has_present) > 1:
        if len(has_present) < 20:
            print has_present
        for i in xrange(len(has_present)):
            if has_present[i]:
                has_present[(i+1)%len(has_present)] = 0
        has_present = filter(None, has_present)
    return has_present[0]

print white_elephant(5)
print white_elephant(3004953)
