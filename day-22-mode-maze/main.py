from heapq import heappush, heappop

def up(p):    return p[0], p[1]-1
def down(p):  return p[0], p[1]+1
def left(p):  return p[0]-1, p[1]
def right(p): return p[0]+1, p[1]

class Map:
    def __init__(self, depth, target):
        self.mouth = (0, 0)
        self.depth = depth
        self.target = target
        self.geo_idxs = {}

    def geo_idx(self, r):
        if r in self.geo_idxs:
            return self.geo_idxs[r]

        stack = [r]
        while stack:
            r, idx = stack.pop(), None
            if r == self.mouth or r == self.target:
                idx = 0
            elif r[0] == 0:
                idx = r[1] * 48271
            elif r[1] == 0:
                idx = r[0] * 16807
            elif left(r) in self.geo_idxs and up(r) in self.geo_idxs:
                idx = self.erosion_lvl(left(r)) * self.erosion_lvl(up(r))

            if idx is None:
                stack.extend((r, left(r), up(r)))
            else:
                self.geo_idxs[r] = idx
        return self.geo_idxs[r]

    def erosion_lvl(self, r):
        return (self.geo_idx(r) + self.depth) % 20183

    def region_type(self, r):
        return self.erosion_lvl(r) % 3

    def total_risk_lvl(self):
        return sum(self.region_type((x, y))
            for x in range(self.target[0]+1)
            for y in range(self.target[1]+1))

ROCKY = 0
WET = 1
NARROW = 2

NOTHING = 0
TORCH = 1
GEAR = 2

EQ = { ROCKY: (TORCH, GEAR), WET: (NOTHING, GEAR), NARROW: (NOTHING, TORCH) }

def adjacent(p):
    if p[0] > 0: yield left(p)
    yield right(p)
    if p[1] > 0: yield up(p)
    yield down(p)

def possible_moves(map, r, eq):
    for nr in adjacent(r):
        if eq in EQ[map.region_type(nr)]:
            yield nr, eq, 1

    for neq in EQ[map.region_type(r)]:
        if neq != eq:
            yield r, neq, 7

def best_route_time(map):
    target, queue, time_map = (map.target, TORCH), [(0, map.mouth, TORCH)], {}

    while queue:
        time, r, eq = heappop(queue)
        if (r, eq) in time_map and time_map[(r, eq)] <= time:
            continue

        time_map[(r, eq)] = time
        if (r, eq) == target:
            return time

        for nr, neq, move_time in possible_moves(map, r, eq):
            heappush(queue, (time+move_time, nr, neq))

    return None

def parse_scan_data(data):
    raw_depth, raw_target = data.splitlines()
    return int(raw_depth[7:]), tuple(map(int, raw_target[8:].split(',')))

scan_data = "depth: 8787\ntarget: 10,725"
m = Map(*parse_scan_data(scan_data))

print('Part 1:', m.total_risk_lvl())
print('Part 2:', best_route_time(m))
