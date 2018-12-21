# Reverse engineered version of the program from the puzzle input
def f(n=0):
    r5, r1 = 65536, 8595037

    while True:
        r1 = (((r1 + (r5 & 255)) & 16777215) * 65899) & 16777215

        if r5 < 256:
            yield r1
            if r1 == n:
                break

            r5 = r1 | 65536
            r1 = 8595037
        else:
            r5 = int(r5 / 256)

def min_f(n=0, iters=1_000_000):
    min_val, cnt = None, 0
    for val in f(n):
        if not min_val or min_val > val:
            min_val = val
            cnt = 0
        else:
            cnt += 1
            if cnt >= iters:
                break
    return min_val

first, last, seen_vals = None, None, set()

for val in f():
    if first is None:
        first = val
    if val in seen_vals:
        break
    seen_vals.add(val)
    last = val

print('Part 1:', first) # 15883666
print('Part 2:', last) # 3253305
print('Part 2?:', min_f()) # 762
