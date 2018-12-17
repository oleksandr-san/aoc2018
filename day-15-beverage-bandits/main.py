GOBLIN = 1
ELF = 2

END_GAME = 1
ATTACK = 2
MOVE = 3
END_TURN = 4

def adjacent(p, map):
    if p[1] > 0:
        yield (p[0], p[1] - 1)
    if p[0] > 0:
        yield (p[0] - 1, p[1])
    if p[0] < len(map[p[1]]) - 1:
        yield (p[0] + 1, p[1])
    if p[1] < len(map) - 1:
        yield (p[0], p[1] + 1)

def adjacent_empty(p, map):
    return filter(lambda p: map[p[1]][p[0]] == '.', adjacent(p, map))

def select_best_step(b, e, map):
    dist_map = [[-1]*len(line) for line in map]
    stack = [(e, 0)]
    visited = set()
    while len(stack):
        p, d = stack.pop(0)
        visited.add(p)
        if dist_map[p[1]][p[0]] != -1:
            continue
        dist_map[p[1]][p[0]] = d

        for a in adjacent_empty(p, map):
            if a not in visited and dist_map[a[1]][a[0]] == -1:
                stack.append((a, d + 1))

    best_step, min_dist = None, None
    for a in adjacent_empty(b, map):
        dist = dist_map[a[1]][a[0]]
        if dist == -1:
            continue
        if min_dist is None or min_dist > dist:
            best_step, min_dist = a, dist
    #print('\n'.join(''.join(str(i) if i >= 0 else '#' for i in line) for line in dist_map))
    return best_step, min_dist

def select_best_path(pos, targets, map):
    reachable_targets = []
    for target in targets:
        step, dist = select_best_step(pos, target, map)
        if step is not None:
            reachable_targets.append((target, step, dist))

    if reachable_targets:
        reachable_targets.sort(key=lambda p: (p[2], p[0][1], p[0][0]))
        return reachable_targets[0][1]
    return None

def select_heroes(map):
    heroes = []
    for y, line in enumerate(map):
        for x, cell in enumerate(line):
            if type(cell) == list:
                heroes.append((x, y))
    return heroes

def select_enemies(pos, map):
    enemy_type = GOBLIN if map[pos[1]][pos[0]][0] == ELF else ELF

    enemies = []
    for y, line in enumerate(map):
        for x, cell in enumerate(line):
            if type(cell) == list and cell[0] == enemy_type:
                enemies.append((x, y))
    return enemies

def select_enemy_to_attack(pos, map):
    enemy_type = GOBLIN if map[pos[1]][pos[0]][0] == ELF else ELF

    weakest_enemy, enemy_pos = None, None
    for a in adjacent(pos, map):
        cell = map[a[1]][a[0]]
        if type(cell) == list and cell[0] == enemy_type:
            if weakest_enemy is None or weakest_enemy[2] > cell[2]:
                weakest_enemy, enemy_pos = cell, a
    return enemy_pos

def select_action(pos, map):
    enemies = select_enemies(pos, map)
    if not enemies:
        return (END_GAME,)

    enemy_to_attack = select_enemy_to_attack(pos, map)
    if enemy_to_attack is not None:
        return (ATTACK, enemy_to_attack)

    best_step = select_best_path(pos, enemies, map)
    if best_step is not None:
        return (MOVE, best_step)

    return (END_TURN,)

def move_hero(pos, action, map):
    new_pos = action[1]
    map[new_pos[1]][new_pos[0]] = map[pos[1]][pos[0]]
    map[pos[1]][pos[0]] = '.'
    return new_pos

def attack_enemy(pos, action, map):
    hero = map[pos[1]][pos[0]]
    enemy_pos = action[1]
    enemy = map[enemy_pos[1]][enemy_pos[0]]

    enemy[2] -= hero[1]
    if enemy[2] <= 0:
        map[enemy_pos[1]][enemy_pos[0]] = '.'
        return enemy
    return None

def calculate_hitpoints(map):
    sum = 0
    for _, line in enumerate(map):
        for _, cell in enumerate(line):
            if type(cell) == list:
                sum += cell[2]
    return sum

def run_game(map, exit_on_dead_elf=False):
    turn = 0

    while True:
        heroes = select_heroes(map)
        if not heroes:
            return turn

        turn += 1
        for hero in heroes:
            hero_obj = map[hero[1]][hero[0]]
            if type(hero_obj) != list or hero_obj[2] <= 0:
                continue

            action = select_action(hero, map)
            if action[0] == END_GAME:
                return turn, True
            elif action[0] == MOVE:
                hero = move_hero(hero, action, map)
                action = select_action(hero, map)
                if action[0] == ATTACK:
                    dead = attack_enemy(hero, action, map)
                    if exit_on_dead_elf and dead is not None and dead[0] == ELF:
                        return turn, False
            elif action[0] == ATTACK:
                dead = attack_enemy(hero, action, map)
                if exit_on_dead_elf and dead is not None and dead[0] == ELF:
                    return turn, False
        print("After round {} there are {} hitpoints left".format(turn, calculate_hitpoints(map)))

def print_map(map):
    def stringify_cell(cell):
        if type(cell) == str:
            return cell
        elif cell[0] == GOBLIN:
            return 'G'
        elif cell[0] == ELF:
            return 'E'
        else:
            return '?'

    for line in map:
        print(
            ''.join(stringify_cell(c) for c in line),
            '\t',
            ', '.join('{}({})'.format(stringify_cell(c), c[2]) for c in line if type(c) == list)
        )

def parse_map(input, elf_power=3):
    def create_cell(c):
        if c == 'G':
            return [GOBLIN, 3, 200]
        elif c == 'E':
            return [ELF, elf_power, 200]
        return c

    return [[create_cell(c) for c in line] for line in input.splitlines() if len(line) > 0]

input = """
#########
#G..G..G#
#.......#
#.......#
#G..E..G#
#.......#
#.......#
#G..G..G#
#########"""

input = """
#######
#.G...#
#...EG#
#.#.#G#
#..G#E#
#.....#
#######
"""

# input = """
# #######
# #G..#E#
# #E#E.E#
# #G.##.#
# #...#E#
# #...E.#
# #######
# """

# input = """
# #########
# #G......#
# #.E.#...#
# #..##..G#
# #...##..#
# #...#...#
# #.G...G.#
# #.....G.#
# #########
# """

input = """
################################
###########...G...#.##..########
###########...#..G#..G...#######
#########.G.#....##.#GG..#######
#########.#.........G....#######
#########.#..............#######
#########.#...............######
#########.GG#.G...........######
########.##...............##..##
#####.G..##G.......E....G......#
#####.#..##......E............##
#####.#..##..........EG....#.###
########......#####...E.##.#.#.#
########.#...#######......E....#
########..G.#########..E...###.#
####.###..#.#########.....E.####
####....G.#.#########.....E.####
#.........#G#########......#####
####....###G#########......##..#
###.....###..#######....##..#..#
####....#.....#####.....###....#
######..#.G...........##########
######...............###########
####.....G.......#.#############
####..#...##.##..#.#############
####......#####E...#############
#.....###...####E..#############
##.....####....#...#############
####.########..#...#############
####...######.###..#############
####..##########################
################################
"""

# map = parse_map(input)

# print(select_enemies((4, 4), map))
# print(select_enemy_to_attack((4, 4), map))
# print(select_best_path((4, 4), [(1, 1), (4, 1), (7, 1), (1, 4), (7, 4), (1, 7), (4, 7), (7, 7)], map))
# print(select_best_step((4, 4), (4, 1), map))

# print(select_action((4, 4), map))
# print(select_action((4, 1), map))

# apply_action((4, 4), select_action((4, 4), map), map)
# print_map(map)

map = parse_map(input, 20)
last_round, success = run_game(map, True)
if success:
    hitpoints = calculate_hitpoints(map)
    print("Final state")
    print_map(map)
    print('\nFull rounds: {}, hitpoints: {}, outcome: {}'.format(last_round-1, hitpoints, (last_round-1)*hitpoints))
else:
    print('\nSome elf died at {} turn'.format(last_round))
