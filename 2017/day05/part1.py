input = [int(i) for i in open('input.txt')]

steps = 0
pc = 0
while pc >= 0 and pc < len(input):
    next = pc + input[pc]
    if input[pc] >= 3:
        input[pc] -= 1
    else:
        input[pc] += 1
    pc = next
    steps += 1

print steps
