from collections import defaultdict

input = open('data.txt').read().strip('^$\n')

DELTA = {'N': 1, 'E': 1j, 'S': -1, 'W': -1j}

dist_map = defaultdict(int)
curr, prev = 0, 0
stack = []

for c in input:
    if c in DELTA:
        curr += DELTA[c]
        if curr in dist_map:
            dist_map[curr] = min(dist_map[curr], dist_map[prev]+1)
        else:
            dist_map[curr] = dist_map[prev]+1
    elif c == '(':
        stack.append(curr)
    elif c == ')':
        curr = stack.pop()
    elif c == '|':
        curr = stack[-1]
    prev = curr

print('Part 1:', max(dist_map.values()))
print('Part 2:', sum(dist >= 1000 for dist in dist_map.values()))
