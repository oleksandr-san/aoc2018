def addr(rs, A, B, C): rs[C] = rs[A] + rs[B]
def addi(rs, A, B, C): rs[C] = rs[A] + B
def mulr(rs, A, B, C): rs[C] = rs[A] * rs[B]
def muli(rs, A, B, C): rs[C] = rs[A] * B
def banr(rs, A, B, C): rs[C] = rs[A] & rs[B]
def bani(rs, A, B, C): rs[C] = rs[A] & B
def borr(rs, A, B, C): rs[C] = rs[A] | rs[B]
def bori(rs, A, B, C): rs[C] = rs[A] | B
def setr(rs, A, B, C): rs[C] = rs[A]
def seti(rs, A, B, C): rs[C] = A
def gtir(rs, A, B, C): rs[C] = 1 if A > rs[B] else 0
def gtri(rs, A, B, C): rs[C] = 1 if rs[A] > B else 0
def gtrr(rs, A, B, C): rs[C] = 1 if rs[A] > rs[B] else 0
def eqir(rs, A, B, C): rs[C] = 1 if A == rs[B] else 0
def eqri(rs, A, B, C): rs[C] = 1 if rs[A] == B else 0
def eqrr(rs, A, B, C): rs[C] = 1 if rs[A] == rs[B] else 0

OPS = {op.__name__: op for op in [ addr, addi, mulr, muli, banr, bani, borr, bori, setr, seti, gtir, gtri, gtrr, eqir, eqri, eqrr, ]}

def read_program(path):
    with open(path) as f:
        ip_reg = int(f.readline()[3:])
        commands = []
        for line in f:
            op, A, B, C = line.split()
            commands.append((op, int(A), int(B), int(C)))
        return ip_reg, commands

def eval_program(rs, ip_reg, commands):
    while rs[ip_reg] >= 0 and rs[ip_reg] < len(commands):
        op, A, B, C = commands[rs[ip_reg]]
        OPS[op](rs, A, B, C)
        rs[ip_reg] += 1
    return rs

# Dissasembled version of a specific program from puzzle input
def eval_fun(rs):
    def factors(n):
        for i in range(1, n + 1):
            if n % i == 0:
                yield i

    n = 900 if rs[0] == 0 else 10_551_300
    rs[0] = sum(factors(n))
    rs[1] = n
    rs[2] = 1
    rs[3] = 257
    rs[4] = n + 1
    rs[5] = n + 1
    return rs

ip_reg, commands = read_program('data.txt')

assert(eval_program([0]*6, ip_reg, commands) == eval_fun([0]*6))

print('Part 1: the value in register 0:', eval_fun([0]*6)[0])
print('Part 2: the value in register 0:', eval_fun([1]+[0]*5)[0])
