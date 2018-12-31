import sys

def coordinate_range(r):
    if type(r) == tuple:
        yield from range(r[0], r[1]+1)
    else:
        yield r

class Map:
    def __init__(self):
        self.tiles = {}
        self.minY = sys.maxsize
        self.maxY = 0

    def add(self, coordinates, value):
        for x in coordinate_range(coordinates[0]):
            for y in coordinate_range(coordinates[1]):
                self.tiles[(x, y)] = value
                self.minY = min(self.minY, y)
                self.maxY = max(self.maxY, y)

    def get(self, coordinate, default=None):
        return self.tiles.get(coordinate, default)

def down(c):  return (c[0], c[1]+1)
def left(c):  return (c[0]-1, c[1])
def right(c): return (c[0]+1, c[1])

def find_bound(shift, map, enter):
    p = enter
    while True:
        if map.get(shift(p)) is not None:
            return (p, True)

        p = shift(p)
        if map.get(down(p)) not in ('#', '~'):
            return (p, False)

def fill_map(map, startX=500):
    levels = [[(startX, map.minY), None, None]]

    while levels:
        level = levels[-1]
        curr = level[0]

        if curr[1] >= map.maxY or map.get(down(curr)) == '|':
            map.add(curr, '|')
            levels.pop()

        elif map.get(down(curr)) is None:
            levels.append([down(curr), None, None])

        elif level[1] is None:
            level[1] = bound, closed = find_bound(left, map, curr)
            if not closed:
                levels.append([down(bound), None, None])

        elif level[2] is None:
            level[2] = bound, closed = find_bound(right, map, curr)
            if not closed:
                levels.append([down(bound), None, None])

        elif level[1] is not None and level[2] is not None:
            lbound, lclosed = level[1]
            if not lclosed:
                nlbound, nlclosed = find_bound(left, map, curr)
                if not nlclosed and lbound != nlbound:
                    level[1] = None
                    continue
                lbound, lclosed = nlbound, nlclosed

            rbound, rclosed = level[2]
            if not rclosed:
                nrbound, nrclosed = find_bound(right, map, curr)
                if not nrclosed and rbound != nrbound:
                    level[2] = None
                    continue
                rbound, rclosed = nrbound, nrclosed

            map.add(tuple(zip(lbound, rbound)), '~' if lclosed and rclosed else '|')
            levels.pop()

def parse_ranges(input):
    ranges = []
    for line in input.splitlines():
        if len(line) == 0:
            continue

        x, y = None, None
        for coordinate in line.split(', '):
            bounds = tuple(map(int, coordinate[2:].split('..')))
            bound = bounds[0] if len(bounds) == 1 else tuple(bounds)
            if coordinate[0] == 'x':
                x = bound
            else:
                y = bound
        ranges.append((x, y))
    return ranges

def count_tiles(map, tiles):
    return sum(1 if t in tiles else 0 for t in map.tiles.values())

m = Map()
with open('data.txt') as f:
    input = f.read()
    for xs, ys in parse_ranges(input):
        m.add((xs, ys), '#')

fill_map(m)
print("The number of tiles that water can reach:", count_tiles(m, ('|', '~')))
print("The number of rest water tiles:", count_tiles(m, ('~',)))