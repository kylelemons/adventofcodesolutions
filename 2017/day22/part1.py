class foo(object):
    def __init__(self):
        self.stored = {}
        self.infected = 0
        self.grid = [list(l.strip()) for l in open('input.txt').readlines()]

    def get(self, pos):
        if pos in self.stored:
            return self.stored[pos]
        if pos[0] < 0 or pos[0] >= len(self.grid):
            return '.'
        if pos[1] < 0 or pos[1] >= len(self.grid[0]):
            return '.'
        return self.grid[pos[0]][pos[1]]

    def store(self, pos, val):
        if val == '#':
            self.infected += 1
        if pos in self.stored:
            self.stored[pos] = val
            return
        if pos[0] < 0 or pos[0] >= len(self.grid):
            self.stored[pos] = val
            return
        if pos[1] < 0 or pos[1] >= len(self.grid[0]):
            self.stored[pos] = val
            return
        self.grid[pos[0]][pos[1]] = val

    def dump(self, what):
        return
        print '------ {} ------'.format(what)
        for line in self.grid:
            print ' '.join(line)

f = foo()

pos = (len(f.grid)/2, len(f.grid[0])/2)

dir = 0
dirs = [
    # urdl
    (-1,0),
    (0,1),
    (1,0),
    (0,-1),
]

f.dump('before')
for burst in range(10000):
    f.dump(burst)
    # print pos, f.get(pos), dir
    if f.get(pos) == '#':
        # turn right
        dir = (dir+1) % 4
        f.store(pos, '.')
    else:
        # turn left
        dir = (4+dir-1) % 4
        f.store(pos, '#')
    pos = (pos[0] + dirs[dir][0], pos[1] + dirs[dir][1])

f.dump('after')
print f.infected
