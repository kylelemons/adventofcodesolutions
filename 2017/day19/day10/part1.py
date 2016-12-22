#!/usr/bin/python2.7

import re

class bot:
    def __init__(self, name, lo_to, hi_to, bots, watch):
        self.name = name
        self.lo_to = lo_to
        self.hi_to = hi_to
        self.bots = bots
        self.watch = watch
        self.has = []

    def give(self, i):
        i = int(i)
        self.has = sorted(self.has + [i])

        if tuple(self.has) == self.watch:
            print self.name, self.has, "ANSWER"
        if len(self.has) == 2:
            self.bots[self.lo_to].give(self.has[0])
            self.bots[self.hi_to].give(self.has[1])
        print self.name, self.has

class output:
    def __init__(self, i):
        self.index = i

    def give(self, i):
        print 'output', self.index, i

inputs = '''value 5 goes to bot 2
bot 2 gives low to bot 1 and high to bot 0
value 3 goes to bot 1
bot 1 gives low to output 1 and high to bot 0
bot 0 gives low to output 2 and high to output 0
value 2 goes to bot 2'''
commands = inputs.splitlines()

commands = open('input.txt').readlines()
watch = (17,61)
bots = {'output '+str(i): output(i) for i in xrange(0,21)}
start = []
for cmd in commands:
    rule = re.match(r'(bot \d+) gives low to ((?:bot|output) \d+) and high to ((?:bot|output) \d+)', cmd)
    give = re.match(r'value (\d+) goes to (bot \d+)', cmd)
    if rule:
        name, a, b = rule.group(1), rule.group(2), rule.group(3)
        bots[name] = bot(name, a, b, bots, watch)
    elif give:
        val, name = give.group(1), give.group(2)
        start.append((val, name))
for val, name in start:
    bots[name].give(val)
