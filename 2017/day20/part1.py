import re
from math import copysign

input = open('input.txt').readlines()

# list of (manhat, pos, vel, accel)
particles = []

linere = re.compile(r'p=<(-?\d+),(-?\d+),(-?\d+)>, v=<(-?\d+),(-?\d+),(-?\d+)>, a=<(-?\d+),(-?\d+),(-?\d+)>')
for line in input:
    m = linere.match(line)
    if not m:
        continue
    g = m.groups()
    p = [int(n) for n in g[0:][:3]]
    v = [int(n) for n in g[3:][:3]]
    a = [int(n) for n in g[6:][:3]]
    particles.append((p,v,a))

# doesn't break ties right
print sorted(enumerate(particles), key=lambda p: sum([abs(n) for n in p[1][2]]))[0]
