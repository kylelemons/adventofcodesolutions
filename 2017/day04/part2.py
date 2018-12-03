phrases = [line.strip().split() for line in open('input.txt')]

valid = 0
for phrase in phrases:
    phrase = [''.join(sorted(list(w))) for w in phrase]
    if len(set(phrase)) == len(phrase):
        valid += 1

print valid
