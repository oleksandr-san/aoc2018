import re
import copy
from collections import defaultdict

GROUP_REGEXP = re.compile("""(\d+) units each with (\d+) hit points (\(.*\)\s)*with an attack that does (\d+) (\w+) damage at initiative (\d+)""")

class Group:
    def __init__(self, id, type, units, hp, dp, dt, initiative, weaknesses=None, immunities=None):
        self.id = id
        self.type = type
        self.units = units
        self.hit_points = hp
        self.damage_points = dp
        self.damage_type = dt
        self.initiative = initiative
        self.weaknesses = weaknesses
        self.immunities = immunities

    @property
    def power(self):
        return self.units * self.damage_points

    @property
    def name(self):
        return '{} group {}'.format(self.type, self.id)

def eval_damage(attacker, defender):
    if attacker.damage_type in defender.weaknesses:
        return attacker.power * 2
    elif attacker.damage_type in defender.immunities:
        return 0
    else:
        return attacker.power

def select_targets(groups):
    ordered_groups = sorted(groups, key=lambda g: (g.power, g.initiative), reverse=True)
    selected_groups = set()
    selections = {}

    for attacker in ordered_groups:
        targets = list(filter(
            lambda g: (g.type != attacker.type and g not in selected_groups),
            ordered_groups))

        if targets:
            defender = max(targets, key=lambda g: (eval_damage(attacker, g), g.power, g.initiative))
            if eval_damage(attacker, defender) != 0:
                selected_groups.add(defender)
                selections[attacker] = defender

    return selections

def attack_targets(groups, selections):
    ordered_groups = sorted(groups, key=lambda g: g.initiative, reverse=True)
    killed_targets = []

    for attacker in ordered_groups:
        if attacker.units <= 0 or attacker not in selections:
            continue

        defender = selections[attacker]
        damage = eval_damage(attacker, defender)
        killed_units = min(defender.units, damage // defender.hit_points)
        defender.units -= killed_units
        if defender.units <= 0:
            killed_targets.append(defender)

    return killed_targets

def parse_group(id, type, input):
    match = GROUP_REGEXP.search(input)
    if match is None:
        return None

    units = int(match.group(1))
    hp = int(match.group(2))
    weaknesses, immunities = [], []
    if match.group(3) is not None:
        for desc in match.group(3).strip('() ').split('; '):
            if desc.startswith('weak to'):
                weaknesses.extend(desc[8:].split(', '))
            elif desc.startswith('immune to'):
                immunities.extend(desc[10:].split(', '))
    dp = int(match.group(4))
    dt = match.group(5)
    initiative = int(match.group(6))

    return Group(id, type, units, hp, dp, dt, initiative, weaknesses, immunities)

def parse_groups(input):
    ids = defaultdict(lambda: 1)
    type = None
    groups = []

    for line in input.splitlines():
        if 'Immune System' in line:
            type = 'immune system'
        elif 'Infection' in line:
            type = 'infection'
        elif line:
            group = parse_group(ids[type], type, line)
            if group:
                ids[type] += 1
                groups.append(group)
    return groups

def run_fight(groups, damage_boost=None):
    fighting_groups = copy.deepcopy(groups)

    if damage_boost:
        for group in fighting_groups:
            group.damage_points += damage_boost.get(group.type, 0)

    while True:
        selections = select_targets(fighting_groups)
        if not selections:
            break
  
        for killed_group in attack_targets(fighting_groups, selections):
            fighting_groups.remove(killed_group)

    return fighting_groups

with open('data.txt') as file:
    groups = parse_groups(file.read())

print('Part 1:', sum(group.units for group in run_fight(groups)))
print('Part 2:', sum(group.units for group in run_fight(groups, {'immune system': 35})))