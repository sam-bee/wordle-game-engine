# Wordle Game Engine

_Game Logic for Wordle_

A small, dependency-light Go package that implements **Wordle-style game logic**: given a hidden solution and a proposed guess, it returns per-letter feedback (**grey / yellow / green**) with correct handling of **duplicate letters**.

The repo contains wordlists needed for playing Wordle.

---

## Features

- **Pure game logic**
  - Represent a game-in-progress (solution + guess history + feedback history)
  - Evaluate a guess against the hidden solution
- **Correct duplicate-letter scoring**
  - Matches Wordle behaviour for repeated letters (greens allocated first, then yellows up to remaining counts)
- **Bundled wordlists**
  - Includes a **solutions** list and a **valid guesses** list (Wordle-style: solutions ⊂ valid guesses)
  - Helpers to load/access these lists without extra dependencies
- **Testable API**
  - Wordlists loaded from `./data/` directory

---

## What's in the box (and what isn't)

### Included
- Scoring logic (grey/yellow/green)
- A `Game` type to track attempts and feedback
- **Wordlists** for solutions and allowed guesses

### Not included
- Solver / strategy code
- UI rendering
- Anything to do with reinforcement learning, directly.

---

## Core concepts

### Tile feedback

A guess is scored into a slice of tiles, one per letter:

Grey: letter not present (or present but all instances already 'consumed' by greens/yellows)

Yellow: letter present in solution, but in a different position

Green: correct letter in the correct position

This package implement`s the standard two-pass approach used by Wordle:

Allocate greens for exact matches.

For remaining letters, allocate yellows only while the solution still has unmatched counts of that letter; otherwise grey.

### Wordlists

This repo includes two lists:

- `allowed-solutions.txt`: the set of possible solutions. 2,309 words.
- `allowed-guesses.txt`: the set of valid guesses. 14,855 words.

## Testing

```
go test ./...
```

The test suite should include:

Known Wordle edge cases (repeated letters in solution and/or guess)

Property-ish checks (e.g., number of greens+yellows for a letter never exceeds its count in the solution)

Regression fixtures (table-driven tests)

## End with 'Why'

For a reinforcement learning project, in which a machine learning model will be produced that can play Wordle well.
