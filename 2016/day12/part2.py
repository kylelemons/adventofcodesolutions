#!/usr/bin/python2.7

from re import *
from itertools import *
from functools import *
from math import *

registers = {k:0 for k in "abcd"}

instructions_text = """cpy 1 a
cpy 1 b
cpy 26 d
jnz c 2
jnz 1 5
cpy 7 c
inc d
dec c
jnz c -2
cpy a c
inc a
dec b
jnz b -2
cpy c b
dec d
jnz d -6
cpy 16 c
cpy 17 d
inc a
dec d
jnz d -2
dec c
jnz c -5"""

instrs = [line.split() for line in instructions_text.splitlines()]
registers['c'] = 1

pc = 0
while pc < len(instrs) and pc >= 0:
    instr = instrs[pc]
    print registers, instr
    if instr[0] == ("cpy"):
        x, y = instr[1], instr[2]
        if y in registers:
            if x >= 'a' and x <= 'z':
                registers[y] = registers[x]
            else:
                registers[y] = int(x)
    elif instr[0] == ("inc"):
        x = instr[1]
        if x in registers:
            registers[x] += 1
    elif instr[0] == ("dec"):
        x = instr[1]
        if x in registers:
            registers[x] -= 1
    elif instr[0] == ("jnz"):
        x, y = instr[1], instr[2]
        try:
            y = int(y)
        except:
            if y in registers:
                y = registers[y]
        try:
            x = int(x)
        except:
            if x in registers:
                x = registers[x]
        if x:
            pc += y
            continue
    #print "after " , instr, registers, instrs
    pc += 1

print registers
