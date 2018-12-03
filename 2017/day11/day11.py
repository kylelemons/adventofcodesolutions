input = open('input.txt').read().strip()

def do(input):
    input = input.split(',')
    ll, lr, ud = 0, 0, 0
    m = 0
    for dir in input:
        if dir == 'ne':
            ll -= 1
        elif dir == 'sw':
            ll += 1
        elif dir == 'nw':
            lr -= 1
        elif dir == 'se':
            lr += 1
        elif dir == 'n':
            ud -= 1
        elif dir == 's':
            ud += 1
        v = sorted([abs(n) for n in [ll, lr, ud]])
        c = v[-1]+v[-2]
        if c > m:
            m = c

    v = sorted([abs(n) for n in [ll, lr, ud]])
    print v[-1]+v[-2], m

do('ne,ne,ne')
do(input)
