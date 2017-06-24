#!/usr/bin/python2.7

from re import *
from itertools import *
from functools import *
from math import *

registers = {k:0 for k in "abcd"}

instructions_text = """cpy 41 a
inc a
inc a
dec a
jnz a 2
dec a"""
instructions_text = """cpy 2 a
tgl a
tgl a
tgl a
cpy 1 a
dec a
dec a"""
instructions_text = """cpy a b
dec b
cpy a d
cpy 0 a
cpy b c
inc a
dec c
jnz c -2
dec d
jnz d -5
dec b
cpy b c
cpy c d
dec d
inc c
jnz d -2
tgl c
cpy -16 c
jnz 1 c
cpy 86 c
jnz 78 d
inc a
inc d
jnz d -2
inc c
jnz c -5"""

instrs = [line.split() for line in instructions_text.splitlines()]
registers['a'] = 12

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
    elif instr[0] == ("tgl"):
        x = instr[1]
        if x in registers:
            offset = int(registers[x])
            if pc+offset >= 0 and pc+offset < len(instrs):
                victim = instrs[pc+offset]
                if victim[0] == "cpy":
                    instrs[pc+offset][0] = "jnz"
                elif victim[0] == "inc":
                    instrs[pc+offset][0] = "dec"
                elif victim[0] == "dec":
                    instrs[pc+offset][0] = "inc"
                elif victim[0] == "jnz":
                    instrs[pc+offset][0] = "cpy"
                elif victim[0] == "tgl":
                    instrs[pc+offset][0] = "inc"
    #print "after " , instr, registers, instrs
    pc += 1

print registers
