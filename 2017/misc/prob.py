#!/usr/bin/python2.7

from re import *
from itertools import *
from functools import *
from md5 import *
from math import *
from copy import *
from fractions import *

def answer(m):
    def follow(s, path):
        '''
        Computes the probabilities of transitioning to end states from the given starting state.

        input:
            s    - state index
            path - previous states
        returns:
            n    - list of probability numerators
            d    - probability denominator
        '''
        prefix = '..'*len(path)
        print '{}path {}'.format(prefix, path+[s])
        def debug(what, n, d):
            'Print the probability matrix represented by n and d'
            print '{}{}'.format(prefix, what)
            for i, nn in enumerate(n):
                rs = ''.join('{:>7}'.format(str(nnn)+'/'+str(d[i])) for nnn in nn)
                print '{}{}'.format(prefix, rs)

        # Terminal state: 100% probability to be in the matching state
        if not any(m[s]):
            print '{}terminal'.format(prefix)
            return [1 if i == s else 0 for i in xrange(len(m))], 1

        # Add the current state to the path (but don't modify the argument)
        path = path + [s]

        # Compute the weighted probability of each next state recursively
        #   wd - denominator of weights of transitory states
        #   n  - list of numerators for transitory states
        #   d  - list of denominators for transitory states
        #   ad - target denominator for sum of weighted probabilities
        #   ss - next state
        #   ww - next state's weight
        wd = sum(v for i, v in enumerate(m[s]) if i not in path)
        n, d, ad = [], [], wd
        for ss, ww in enumerate(m[s]):
            if ww == 0 or ss in path:
                continue
            nn, dd = follow(ss, path)
            n.append([ww*nnn for nnn in nn])
            d.append(wd*dd)
            ad *= dd

        print '{}at path {}'.format(prefix, path)
        debug('from following', n, d)
        print '{}target x/{}'.format(prefix, ad)

        # Scale the numerators to match the target denominator
        for i, nn in enumerate(n):
            # ... the target denominator is the portion of ad not covered by d
            scale, d[i] = ad / d[i], ad
            n[i] = [v * scale for v in nn]
        debug('after scale', n, d)

        # Sum the columns to compute the total probability
        total = [sum(col) for col in izip(*n)]
        debug('column sum', [total], [ad])

        # Reduce the fractions to keep things tidy and in-bounds
        descale = reduce(gcd, total, ad)
        print '{}gcd = {}'.format(prefix, descale)

        return [t/descale for t in total], ad/descale

    n, d = follow(0, [])
    return [v for i, v in enumerate(n) if not any(m[i])] + [d]

tests = [
    # Non-recursive state transitions
    (([[0,1,1,0],[0,0,1,2],[0]*4,[0]*4],), [2,1,3]),
    # Mutual recursion
    (([[0,1,1,0],[1,0,3,4],[0]*4,[0]*4],), [5,2,7]),
    # Scaling logic (/2 scale to /6 before sum)
    (([[0,1,1,0],[1,0,1,2],[0]*4,[0]*4],), [2,1,3]),
    # Self-recursion: two coin flips to choose between 3 options
    (([[1,1,1,1],[0,0,0,0],[0,0,0,0],[0,0,0,0]],), [1,1,1,3]),
]
for args, ans in tests:
    got = answer(*args)
    print args, got, ans, got == ans
