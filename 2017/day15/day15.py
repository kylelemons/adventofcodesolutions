# Generator A starts with 277
# Generator B starts with 349
input = (277, 349)

# Part 1
from itertools import starmap, islice

def gen(val, factor, mod, test):
    while True:
        val = (val*factor) % mod
        if val % test != 0:
            continue
        yield val

def pairs(a, b, atest, btest):
    mod = 2147483647
    gena = gen(a, 16807, mod, atest)
    genb = gen(b, 48271, mod, btest)
    while True:
        yield next(gena)&0xFFFF == next(genb)&0xFFFF

def judge_count(a, b, count, atest=1, btest=1):
    return sum(1 if eq else 0 for eq in islice(pairs(a,b,atest,btest), count))

#print judge_count(65, 8921, 5)
#print judge_count(input[0], input[1], 40000000)

# Part 2

count = 5e6
print judge_count(65, 8921, count, 4, 8)
print judge_count(input[0], input[1], count, 4, 8)
