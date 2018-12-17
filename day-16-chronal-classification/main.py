from collections import defaultdict

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

OPS = [ addr, addi, mulr, muli, banr, bani, borr, bori, setr, seti, gtir, gtri, gtrr, eqir, eqri, eqrr, ]

def match_ops(brs, ars, A, B, C):
    matched = []
    for op in OPS:
        rs = brs[:]
        op(rs, A, B, C)
        if rs == ars:
            matched.append(op)
    return matched

def read_samples(path):
    samples = []
    with open(path) as f:
        for raw_sample in f.read().split('\n\n'):
            raw_sample = raw_sample.split('\n')

            brs = [int(c) for c in raw_sample[0][9:-1].split(', ')]
            ars = [int(c) for c in raw_sample[2][9:-1].split(', ')]
            code, A, B, C = (int(c) for c in raw_sample[1].split(' '))

            samples.append((code, brs, (A, B, C), ars))
    return samples

def eval_op_codes(matched_samples):
    ambiguous_ops = defaultdict(set)
    for code, ops in matched_samples:
        for op in ops:
            ambiguous_ops[op].add(code)

    op_code = dict()
    while len(op_code) != len(OPS):
        for op, codes in ambiguous_ops.items():
            codes.difference_update(op_code.keys())
            if len(codes) == 1:
                op_code[codes.pop()] = op

        for op in op_code.values():
            if op in ambiguous_ops:
                ambiguous_ops.pop(op)
    return op_code

def eval_program(path, ops):
    rs = [0] * 4
    with open(path) as f:
        for line in f:
            code, A, B, C = [int(c) for c in line.split()]
            ops[code](rs, A, B, C)
    return rs

matched_samples = [(code, match_ops(brs, ars, *args)) for code, brs, args, ars in read_samples('samples.txt')]
print('Number of samples that match three or more opcodes:', sum(1 if len(ops) >= 3 else 0 for _, ops in matched_samples))
print('Value of register 0 after program execution:', eval_program('program.txt', eval_op_codes(matched_samples))[0])
