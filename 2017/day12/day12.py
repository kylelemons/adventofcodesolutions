input = open('input.txt').readlines()
input = [line.strip() for line in input]

class c(object):
    def __init__(self):
        self.roots = {}
        pass
    
    def root(self, n):
        if n in self.roots:
            return self.roots[n]
        self.roots[n] = n
        return n

    def combine(self, n1, n2):
        for k, v in self.roots.iteritems():
            if v == n2:
                self.roots[k] = n1

    def do(self, input):
        for line in input:
            prog, neigh = line.split(' <-> ')
            neigh = neigh.split(', ')
            r = self.root(prog)
            for n in neigh:
                self.combine(r, self.root(n))
        
        r = self.root('0')
        return (
            len([n for n in self.roots if self.root(n) == r]),
            len(set(self.roots.itervalues())),
        )

print c().do(input)
