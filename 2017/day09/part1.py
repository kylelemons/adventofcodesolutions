input = open('input.txt').read().strip()

def part12(input):
    # (depth, startpos)
    state = []
    depth = 1
    score = 0
    i = 0
    garbage = False
    removed = 0
    while i < len(input):
        ch = input[i]
        if ch == '!':
            i += 1
        elif garbage and ch == '>':
            garbage = False
        elif garbage:
            removed += 1
            pass
        elif ch == '{':
            state.append((depth, i))
            depth += 1
        elif ch == '}':
            (cdepth, startpos) = state.pop()
            score += cdepth
            depth -= 1
        elif ch == '<':
            garbage = True
        i += 1
    if len(input) > 10:
        input = input[:10]+'...'
    print 'score("{}") = {} (removed: {})'.format(input, score, removed)

part12('{{<!!>},{<!!>},{<!!>},{<!!>}}')
part12(input)
