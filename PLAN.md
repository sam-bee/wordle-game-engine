# Project Plan

## Overall Goal

We need a Go service listening on port 9111. It will be receiving POST requests from a local Python programme, which is a machine learning model. The ML model is using reinforcement learning to learn to play Wordle.

This Go application has knowledge of the Wordle wordlists for valid guesses, and valid solutions, in the `./data/` folder.

The API will be very simple, with no extraneous features or fields. There will not be versioning.

### Input

When a request comes in, the details in the request should be the correct solution, and between 0-5 turns. A turn will have the word played, and the feedback. Feedback looks like 'GY--Y', where G = green tile in a Wordle game, Y = yellow tile, '-' = grey

### Output

When sent data about the state of a game to date, and a proposed new move, the Go service responds with a JSON body containing the following information:
- Whether the game is won, lost or ongoing after the latest turn. You get 6 turns.
- Whether the turn is valid or invalid. It should have been a 5-letter word in lower case on the allowed guess list
- The 'solution shortlist reduction'. This is the number of possible solutions before the latest turn, the number of possible solutions after the latest turn, and the reduction on a scale of 0-1.
- The 'feedback'.

### Caching

There will be a large, in-memory cache. It will be used for identifying the solution shortlist after the first turn in a game only. It will not be possible to cache solution shortlists for subsequent terms. The cache keys are the first turn the player took, and the feedback string. The value is the entire remaining shortlist, which will often contain hundreds or thousands of items. No more than 10GB should be used for the cache, roughly

## Progress so far

A Go wordle engine has been implemented in pkg/wordlegameengine. A lightweight HTTP server has been added in main.go, listening on port 9111 with /api/evaluate POST endpoint accepting JSON ({&quot;solution&quot;: &quot;...&quot;, &quot;turns&quot;: [...], &quot;proposed_guess&quot;: &quot;...&quot;}) and returning dummy JSON matching spec: {&quot;game_status&quot;: &quot;...&quot;, &quot;turn_valid&quot;: bool, &quot;shortlist_reduction&quot;: {&quot;before&quot;: int, &quot;after&quot;: int, &quot;ratio&quot;: float}, &quot;feedback&quot;: &quot;...&quot;}. Dummy validation (length/lowercase). Tests pass. No caching yet.

## Completed Iterations

### ✅ Iteration 1: Integrate Engine for Request Validation

**Status: COMPLETE**

All acceptance criteria met:
1. ✅ On startup, main.go calls wordlegameengine.LoadWordlists('./data'); log.Fatal if error.
2. ✅ Import "./pkg/wordlegameengine" and "log".
3. ✅ In evaluateHandler validates req.Solution, req.ProposedGuess, and req.Turns using engine
4. ✅ Tests updated with 10 test cases covering real word/non-word validation
5. ✅ `go test ./...` passes (25 tests total)
6. ✅ Server responds 400 with descriptive error for invalid words

## Current Iteration: Compute Real Game Status

### Acceptance Criteria

1. Use `wordlegameengine.NewGame(solution)` to create a game instance in evaluateHandler.

2. For each turn in `req.Turns`:
   - Call `game.MakeGuess(turn.Guess)` to apply the turn
   - Verify the returned `Feedback.String()` matches `turn.Feedback` from the request
   - If feedback doesn't match, return 400 with descriptive error

3. Determine `game_status` based on game state after applying all turns:
   - `"won"` if `game.IsWon()` is true
   - `"lost"` if `game.IsLost()` is true (6 turns played and not won)
   - `"ongoing"` otherwise

4. Set `TurnValid = true` only if all guesses are valid AND feedbacks match.

5. `go test ./...` passes.

6. Update/add tests for:
   - Game won in 1-6 turns
   - Game lost after 6 turns
   - Game ongoing with 0-5 turns
   - Feedback mismatch returns 400

### Tasks

- [ ] **go-coder**: Implement real game status computation using Game/MakeGuess/IsWon/IsLost. Run tests.

- [ ] **qa-requirements**: Verify AC met, tests pass, game status computed correctly, feedback verification works.

## Future Iterations

- Compute real feedback for proposed_guess against solution.

- Compute real shortlist reduction using Game logic.

- Implement caching for first turn.

