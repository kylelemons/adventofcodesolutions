commands = [line.strip().split(' ') for line in open('input.txt')]
optimized = [line.strip().split(' ') for line in open('optimized.txt')]
primes = set([int(i) for line in open('primes.txt') for i in line.split('\t')])
print primes

def get(regs, v):
    try:
        return int(v)
    except ValueError:
        if v not in regs:
            return 0
        return regs[v]

def part1(commands):
    pc = 0
    regs = {}

    muls = 0
    while 0 <= pc < len(commands):
        #print regs
        c, a, b = commands[pc]
        if c == 'set':
            regs[a] = get(regs, b)
        elif c == 'sub':
            regs[a] = get(regs, a) - get(regs, b)
        elif c == 'mul':
            regs[a] = get(regs, a) * get(regs, b)
            muls += 1
        elif c == 'jnz':
            if get(regs, a) != 0:
                pc += get(regs, b)
                continue
        pc += 1

    print 'MUL count: {}'.format(muls)

def part2(commands):
    pc = 0
    regs = {'a': 1}
    watchdog = 10000000

    muls = 0
    while 0 <= pc < len(commands) and watchdog >= 0:
        watchdog -= 1
        c, a, b = commands[pc]
        if c == 'isprime':
            print pc+1, c, regs
        if c == 'set':
            regs[a] = get(regs, b)
        elif c == 'sub':
            regs[a] = get(regs, a) - get(regs, b)
        elif c == 'mul':
            regs[a] = get(regs, a) * get(regs, b)
            muls += 1
        elif c == 'jnz':
            if get(regs, a) != 0:
                pc += get(regs, b)
                continue
        elif c == 'isprime':
            regs[a] = 1 if get(regs, b) in primes else 0
        pc += 1

    print get(regs, 'h')

part2(optimized)

"""
# b := 108100
# c := 125100
1   set b 81
2   set c b
3   jnz a 2
4   jnz 1 5
5   mul b 100
6   sub b -100000
7   set c b
8   sub c -17000

# d := 2
# f := 1
9   set f 1
10  set d 2

# e := 2
11  set e 2

# if d*e != b: goto 18
12  set g d
13  mul g e
14  sub g b
15  jnz g 2

# if d*e == b: f := 0
16  set f 0

# if e != b: e++, goto 12
17  sub e -1
18  set g e
19  sub g b
20  jnz g -8

# Past here: e == b

# d++
21  sub d -1

# if d != b: goto 11
22  set g d
23  sub g b
24  jnz g -13

# Past here: d == b == e
# If d*e was ever b, f is 0

# If a d,e pair was found, h++
25  jnz f 2
26  sub h -1

# if b == c: exit
27  set g b
28  sub g c
29  jnz g 2
30  jnz 1 3

# b += 17, goto 9
31  sub b -17
32  jnz 1 -23

---
b := 108100
c := 125100
d := 2
f := 1

for ; d != b; d++ {
    e := 2
    while e != b:
        if d*e == b:
            f := 0
            e++
}
"""
