

class TM(object):
    def __init__(self):
        self.tape = {}
        self.cursor = 0
        self.state = 'A'

    def read(self):
        if self.cursor not in self.tape:
            return 0
        return self.tape[self.cursor]

    def write(self, value):
        self.tape[self.cursor] = value
    
    def left(self):
        self.cursor -= 1

    def right(self):
        self.cursor += 1

    def A(self):
        if self.read() == 0:
            self.write(1)
            self.right()
            self.state = 'B'
        else:
            self.write(0)
            self.left()
            self.state = 'D'

    def B(self):
        if self.read() == 0:
            self.write(1)
            self.right()
            self.state = 'C'
        else:
            self.write(0)
            self.right()
            self.state = 'F'

    def C(self):
        if self.read() == 0:
            self.write(1)
            self.left()
            self.state = 'C'
        else:
            self.write(1)
            self.left()
            self.state = 'A'

    def D(self):
        if self.read() == 0:
            self.write(0)
            self.left()
            self.state = 'E'
        else:
            self.write(1)
            self.right()
            self.state = 'A'

    def E(self):
        if self.read() == 0:
            self.write(1)
            self.left()
            self.state = 'A'
        else:
            self.write(0)
            self.right()
            self.state = 'B'

    def F(self):
        if self.read() == 0:
            self.write(0)
            self.right()
            self.state = 'C'
        else:
            self.write(0)
            self.right()
            self.state = 'E'

    def step(self):
        if self.state == 'A': self.A()
        elif self.state == 'B': self.B()
        elif self.state == 'C': self.C()
        elif self.state == 'D': self.D()
        elif self.state == 'E': self.E()
        elif self.state == 'F': self.F()

    def checksum(self):
        return sum(1 if v == 1 else 0 for v in self.tape.itervalues())

m = TM()
for i in range(12302209):
    m.step()

print m.checksum()
