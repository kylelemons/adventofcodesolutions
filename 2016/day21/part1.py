#!/usr/bin/python2.7

import re

class password:
    def __init__(self, start):
        self.pw = [x for x in start]

    def swapPos(self, x, y):
        x, y = int(x), int(y)
        self.pw[x], self.pw[y] = self.pw[y], self.pw[x]

    def swapLetter(self, x, y):
        mapping = {x:y, y:x}
        self.pw = [mapping[ch] if ch in mapping else ch for ch in self.pw]

    def rotateLeft(self, by):
        by = int(by)
        self.pw = self.pw[by:] + self.pw[:by]

    def rotateRight(self, by):
        by = int(by)
        self.pw = self.pw[-by:] + self.pw[:-by]

    def rotateFromLetter(self, x):
        offset = self.pw.index(x)
        if offset >= 4:
            offset += 1
        offset += 1
        offset %= len(self.pw)
        self.rotateRight(offset)

    def reverseRange(self, x, y):
        x, y = int(x), int(y)
        rev = reversed(self.pw[x:y+1])
        self.pw = self.pw[:x] + list(rev) + self.pw[y+1:]

    def move(self, x, y):
        x, y = int(x), int(y)
        if x < y:
            self.pw = self.pw[:x] + self.pw[x+1:y+1] + [self.pw[x]] + self.pw[y+1:]
        else:
            self.pw = self.pw[:y] + [self.pw[x]] + self.pw[y:x] + self.pw[x+1:]

    def __str__(self):
        return ''.join(self.pw)

    def eval(self, command):
        swapPos = re.match('swap position (\d+) with position (\d+)', command)
        swapLetter = re.match('swap letter (\S+) with letter (\S+)', command)
        rotateLeft = re.match('rotate left (\d+) steps?', command)
        rotateRight = re.match('rotate right (\d+) steps?', command)
        rotateFromLetter = re.match('rotate based on position of letter (\S+)', command)
        reverseRange = re.match('reverse positions? (\d+) through (\d+)', command)
        move = re.match('move position (\d+) to position (\d+)', command)

        if swapPos:
            self.swapPos(swapPos.group(1), swapPos.group(2))
        elif swapLetter:
            self.swapLetter(swapLetter.group(1), swapLetter.group(2))
        elif rotateLeft:
            self.rotateLeft(rotateLeft.group(1))
        elif rotateRight:
            self.rotateRight(rotateRight.group(1))
        elif rotateFromLetter:
            self.rotateFromLetter(rotateFromLetter.group(1))
        elif reverseRange:
            self.reverseRange(reverseRange.group(1), reverseRange.group(2))
        elif move:
            self.move(move.group(1), move.group(2))
        else:
            print "unknown command", command

pw = password('abcdefgh')
for cmd in open('input.txt').readlines():
    pw.eval(cmd)
    print cmd, pw
