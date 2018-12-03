import re

input = [line.strip() for line in open('input.txt').readlines()]

parse = re.compile('(.+) (inc|dec) (-?\d+) if (.+) ([><=!]+) (-?\d+)')

regs = {}
during = 0
for line in input:
    match = parse.match(line)
    if not match:
        print 'bad line: {}'.format(line)
        continue
    (victim, iod, by, target, op, lim) = match.groups()
    by, lim = int(by), int(lim)
    
    prev, val = 0, 0
    if victim in regs:
        prev = regs[victim]
    if target in regs:
        val = regs[target]

    if prev > during:
        during = prev

    do = False
    if op == "==":
        do = val == lim
    elif op == "!=":
        do = val != lim
    elif op == ">":
        do = val > lim
    elif op == "<":
        do = val < lim
    elif op == ">=":
        do = val >= lim
    elif op == "<=":
        do = val <= lim
    else:
        print 'bad op: {}'.format(op)
    if not do:
        continue

    if iod == 'inc':
        regs[victim] = prev + by
    elif iod == 'dec':
        regs[victim] = prev - by
    else:
        print 'bad iod: {}'.format(iod)

print sorted([(v, k) for k, v in regs.iteritems()])[-1]
print during
