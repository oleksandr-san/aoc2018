import re

def load_nanobots(path):
    regex = re.compile("""pos=<(-?\d+),(-?\d+),(-?\d+)>, r=(-?\d+)""")
    bots = []
    with open(path) as file:
        for line in file:
            match = regex.match(line)
            p = tuple(int(match.group(i)) for i in range(1, 4))
            r = int(match.group(4))
            bots.append((p, r))
    return bots

def manhattan_dist(p1, p2):
    return sum(abs(c1 - c2) for c1, c2 in zip(p1, p2))

bots = load_nanobots('data.txt')
strongest_bot = max(bots, key=lambda bot: bot[1])

print('Part 1:', sum(manhattan_dist(strongest_bot[0], bot[0]) <= strongest_bot[1] for bot in bots))
