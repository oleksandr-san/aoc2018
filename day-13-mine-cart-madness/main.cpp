#include <vector>
#include <map>
#include <iostream>
#include <fstream>
#include <string>
#include <optional>

enum class Direction {
    Unknown,
    Up,
    Down,
    Left,
    Right
};

auto left(Direction d) {
    switch (d) {
    case Direction::Up:    return Direction::Left;
    case Direction::Down:  return Direction::Right;
    case Direction::Left:  return Direction::Down;
    case Direction::Right: return Direction::Up;
    default:               return Direction::Unknown;
    }
}

auto right(Direction d) {
    switch (d) {
    case Direction::Up:    return Direction::Right;
    case Direction::Down:  return Direction::Left;
    case Direction::Left:  return Direction::Up;
    case Direction::Right: return Direction::Down;
    default:               return Direction::Unknown;
    }
}

auto curveLeft(Direction d) {
    switch (d) {
    case Direction::Up:
    case Direction::Down:
        return left(d);
    case Direction::Left:
    case Direction::Right:
        return right(d);
    default:
        return Direction::Unknown;
    }
}

auto curveRight(Direction d) {
    switch (d) {
    case Direction::Up:
    case Direction::Down:
        return right(d);
    case Direction::Left:
    case Direction::Right:
        return left(d);
    default:
        return Direction::Unknown;
    }
}

auto toDirection(char c) {
    switch (c) {
    case '^': return Direction::Up;
    case 'v': return Direction::Down;
    case '<': return Direction::Left;
    case '>': return Direction::Right;
    default:  return Direction::Unknown;
    }
}

auto toChar(Direction d) {
    switch (d) {
    case Direction::Up:    return '^';
    case Direction::Down:  return 'v';
    case Direction::Left:  return '<';
    case Direction::Right: return '>';
    default:               return '?';
    }
}

enum class Path {
    Unknown,
    StraightH,
    StraightV,
    CurveR,
    CurveL,
    Intersection
};

auto toPath(char c) {
    switch (c) {
    case '-':  return Path::StraightH;
    case '|':  return Path::StraightV;
    case '/':  return Path::CurveR;
    case '\\': return Path::CurveL;
    case '+':  return Path::Intersection;
    default:   return Path::Unknown;
    }
}

auto toPath(Direction d) {
    switch (d) {
    case Direction::Up:
    case Direction::Down:
        return Path::StraightV;
    case Direction::Left:
    case Direction::Right:
        return Path::StraightH;
    default:
        return Path::Unknown;
    }
}

auto toChar(Path p) {
    switch (p) {
    case Path::StraightH:    return '-';
    case Path::StraightV:    return '|';
    case Path::CurveR:       return '/';
    case Path::CurveL:       return '\\';
    case Path::Intersection: return '+';
    default:                 return '?';
    }
}

struct Location {
    size_t x, y;
};

struct Cart {
    Direction dir;
    uint8_t option;
};

struct Less {
    auto operator()(const Location& l, const Location& r) const {
        return l.y == r.y ? l.x < r.x : l.y < r.y;
    }
};

using Carts = std::map<Location, Cart, Less>;
using Tracks = std::vector<std::map<size_t, Path>>;

auto parseTracksMap(std::istream& stream) -> std::pair<Tracks, Carts> {
    auto [ tracks, carts, y, line ] = std::tuple<Tracks, Carts, size_t, std::string>{};

    for (; std::getline(stream, line); ++y) {
        auto& tracksLine = tracks.emplace_back();
        for (size_t x = 0; x < line.size(); ++x) {
            if (line[x] == ' ')
                continue;

            if (auto path = toPath(line[x]); path != Path::Unknown) {
                tracksLine[x] = path;
            } else if (auto dir = toDirection(line[x]); dir != Direction::Unknown) {
                tracksLine[x] = toPath(dir);
                carts[{ x, y }] = { dir, 0 };
            }
        }
    }

    return { tracks, carts };
}

auto printTracksMap(const Tracks& tracks, const Carts& carts, std::ostream& stream) {
    for (size_t y = 0; y < tracks.size(); ++y) {
        size_t prevX = 0;
        for (auto [x, path] : tracks[y]) {
            for (auto i = prevX; i < x; ++i) stream << ' ';

            if (auto it = carts.find(Location{x, y}); it != carts.end())
                stream << toChar(it->second.dir);
            else
                stream << toChar(path);

            prevX = x + 1;
        }
        stream << '\n';
    }
}

auto moveLocation(Location loc, Direction d) -> Location {
    switch(d) {
    case Direction::Up:    --loc.y; break;
    case Direction::Down:  ++loc.y; break;
    case Direction::Left:  --loc.x; break;
    case Direction::Right: ++loc.x; break;
    }
    return loc;
}

auto moveCart(Cart cart, Path path) -> Cart {
    switch (path) {
    case Path::CurveR:
        cart.dir = curveRight(cart.dir);
        break;
    case Path::CurveL:
        cart.dir = curveLeft(cart.dir);
        break;
    case Path::Intersection:
        switch (cart.option) {
        case 0: cart.dir = left(cart.dir); break;
        case 2: cart.dir = right(cart.dir); break;
        }
        cart.option = (cart.option + 1) % 3;
        break;
    }

    return cart;
}

auto moveCart(Location loc, Cart cart, const Tracks& tracks) -> std::pair<Location, Cart> {
    loc = moveLocation(loc, cart.dir);
    return { loc, moveCart(cart, tracks.at(loc.y).at(loc.x)) };
}

auto locateFirstCrash(const Tracks& tracks, Carts carts) -> Location {
    while (true) {
        Carts movedCarts;

        for (auto [location, cart] : carts) {
            auto [movedLocation, movedCart] = moveCart(location, cart, tracks);
            if (!carts.count(movedLocation) && !movedCarts.erase(movedLocation))
                movedCarts[movedLocation] = movedCart;
            else
                return movedLocation;
        }

        carts = std::move(movedCarts);
    }
}

auto locateLastCart(const Tracks& tracks, Carts carts) -> std::optional<Location> {
    while (carts.size() > 1) {
        Carts movedCarts;

        for (auto it = carts.begin(); it != carts.end();) {
            auto [location, cart] = moveCart(it->first, it->second, tracks);
            if (auto found = carts.find(location); found != carts.end()) {
                it = carts.erase(found);
            } else if (movedCarts.erase(location)) {
                ++it;
            } else {
                ++it;
                movedCarts[location] = cart;
            }
        }

        carts = std::move(movedCarts);
    }

    if (!carts.empty())
        return carts.begin()->first;
    return {};
}

int main(int argc, char** argv) {
    if (argc < 2) {
        std::cerr << "Not enough arguments\n";
        return -1;
    }

    std::fstream fs{argv[1]};
    if (!fs.is_open()) {
        std::cerr << "Can't open file " << argv[1] << '\n';
        return -1;
    }

    auto [ tracks, carts ] = parseTracksMap(fs);

    auto l1 = locateFirstCrash(tracks, carts);
    std::cout << "First crash at " << l1.x  << "," << l1.y << "\n";

    if (auto l2 = locateLastCart(tracks, carts))
        std::cout << "Last cart at " << l2->x << "," << l2->y << "\n";
    else
        std::cout << "No carts left\n";

    return 0;
}