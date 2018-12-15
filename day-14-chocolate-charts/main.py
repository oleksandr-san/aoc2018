def digits(n):
    return divmod(n, 10) if n >= 10 else (n,)

def eval_recipes(shift, cnt):
    recipes, elves = [3, 7], [0, 1]

    while len(recipes) < (shift + cnt):
        recipes.extend(digits(sum(recipes[elf] for elf in elves)))
        elves = [(elf + recipes[elf] + 1) % len(recipes) for elf in elves]

    return recipes[shift:shift+cnt]

def eval_position(seq):
    recipes, elves = [3, 7], [0, 1]

    while True:
        for digit in digits(sum(recipes[elf] for elf in elves)):
            recipes.append(digit)
            if recipes[-len(seq):] == seq:
                return len(recipes) - len(seq)

        elves = [(elf + recipes[elf] + 1) % len(recipes) for elf in elves]

input = "327901"

print('Part 1:', ''.join(map(str, eval_recipes(int(input), 10))))
print('Part 2:', eval_position(list(map(int, input))))
