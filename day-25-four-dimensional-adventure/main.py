def manhattan_dist(p1, p2):
    return sum(abs(c1 - c2) for c1, c2 in zip(p1, p2))

def count_constellations(points):
    consts = []
    for p in points:
        matched_consts = []
        for const in consts:
            for pc in const:
                if manhattan_dist(p, pc) <= 3:
                    matched_consts.append(const)
                    break

        if not matched_consts:
            consts.append([p])
        elif len(matched_consts) == 1:
            matched_consts[0].append(p)
        else:
            merged_const = matched_consts.pop()
            merged_const.append(p)
            for const in matched_consts:
                merged_const.extend(const)
                consts.remove(const)
    return len(consts)

points = []
with open('data.txt') as file:
    for line in file:
        points.append(tuple(map(int, line.strip().split(','))))

print('Part 1:', count_constellations(points))