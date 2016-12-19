#!/usr/bin/python2.7

import re

class screen:
    def __init__(self, wide, tall):
        self.screen = [['.']*wide for r in xrange(tall)]

    def __str__(self):
        return '\n'.join([''.join(r) for r in self.screen])

    def rect(self, a, b):
        a, b = int(a), int(b)
        for r in xrange(b):
            for c in xrange(a):
                self.screen[r][c] = '#'

    def rotcol(self, x, by):
        x, by = int(x), int(by)
        c = [self.screen[i][x] for i in xrange(len(self.screen))]
        c = c[-by:] + c[:len(c)-by]
        for i in xrange(len(self.screen)):
            self.screen[i][x] = c[i]

    def rotrow(self, y, by):
        y, by = int(y), int(by)
        c = [self.screen[y][i] for i in xrange(len(self.screen[y]))]
        c = c[-by:] + c[:len(c)-by]
        for i in xrange(len(self.screen[y])):
            self.screen[y][i] = c[i]

    def lit(self):
        return sum([r.count('#') for r in self.screen])

def execute(s, command):
    rect = re.match(r'rect (\d+)x(\d+)', command)
    rotcol = re.match(r'rotate column x=(\d+) by (\d+)', command)
    rotrow = re.match(r'rotate row y=(\d+) by (\d+)', command)
    if rect:
        s.rect(*rect.groups())
    elif rotcol:
        s.rotcol(*rotcol.groups())
    elif rotrow:
        s.rotrow(*rotrow.groups())

s = screen(7, 3)
#s.rect(3,2)
execute(s, "rect 3x2")
#s.rotcol(1,1)
execute(s, "rotate column x=1 by 1")
#s.rotrow(0,4)
execute(s, "rotate row y=0 by 4")
#s.rotcol(1,1)
execute(s, "rotate column x=1 by 1")
print s
print s.lit()

real = screen(50, 6)
commands = open('input.txt').readlines()
for c in commands:
    execute(real, c)
print "real:"
print real
print real.lit()
