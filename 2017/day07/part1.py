import re
from collections import defaultdict

below = {}
above = {}
weights = {}

for line in open('input.txt'):
    match = re.match(r'(\w+) \((\d+)\)(?: -> (.*))?', line.strip())
    if not match:
        continue
    (bot, weight, carried) = match.groups()
    weights[bot] = int(weight)
    if not carried:
        continue
    for cbot in carried.split(', '):
        below[cbot] = bot
        above.setdefault(bot, []).append(cbot)

base = None
for bot in weights:
    if bot not in below:
        base = bot

print 'base:', base

# returns (weight_of_base_and_carried, weight_of_carried, balanced)
def trace(below, above, weights, base):
    weight = weights[base]
    carried = 0
    balanced = True
    subs = defaultdict(int)
    subcarries = {}
    if base in above:
        for bot in above[base]:
            (sub, subcarried, subbal) = trace(below, above, weights, bot)
            if not subbal:
                balanced = False
            weight += sub
            carried += sub
            subs[sub] += 1
            subcarries[sub] = subcarried
    if not balanced:
        return (weight, carried, balanced)
    if len(subs) <= 1:
        return (weight, carried, True)
    print subs, subcarries
    return (weight, carried, False)

trace(below, above, weights, base)
