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

A Go wordle engine has been implemented, either completely or at least partially. The application layer hasn't been done. There is no caching.

## Remaining Development Plan

1. Create `pkg/wordlegameengine/` directory. Move `word.go`, `solution.go`, `wordlists.go`, `game.go`, `word_test.go`, `solution_test.go`, `game_test.go`, `wordlists_test.go` to `pkg/wordlegameengine/`. Update all imports in moved files to use relative paths or module path if needed.
4. Run `go mod tidy`.
5. Run `go test ./pkg/wordlegameengine/...` to verify.
6. Create `main.go` in project root.
7. In `main.go`, import \"github.com/sam-bee/wordle-game-engine/pkg/wordlegameengine\" and std libs (net/http, encoding/json, sync).
8. Define structs: `type Turn struct { Guess string \`json:\"guess\"\` Feedback string \`json:\"feedback\"\` }`, `type Request struct { Solution string \`json:\"solution\"\` History []Turn \`json:\"history\"\` Proposed string \`json:\"proposed\"\` }`, `type Response struct { Status string \`json:\"status\"\` Valid bool \`json:\"valid\"\` Before int \`json:\"before\"\` After int \`json:\"after\"\` Reduction float64 \`json:\"reduction\"\` Feedback string \`json:\"feedback\"\` }`.
9. Implement `Feedback.String() string` in pkg/wordlegameengine/feedback.go or extend Feedback: map Grey:'-', Yellow:'Y', Green:'G'.
10. Implement `parseFeedback(s string) (Feedback, error)` vice versa.
11. Implement `FilterSolutions(history []Turn) []Word` function that filters AllowedSolutions where for each turn, candidate.CheckGuess(NewWord(guess)) == parseFeedback(fb).
12. Use parallel workers similar to game.updateSolutionShortlist for efficiency.
13. In main, LoadWordlists(\"./data\") in init() or main.
14. var cache = make(map[string][]Word); var mu sync.RWMutex
15. POST handler `/api/evaluate`: parse req, validate Proposed with NewWord(proposed).Validate()==nil -> valid=true
16. sol := NewSolution(req.Solution)
17. guess := NewWord(req.Proposed)
18. fb := sol.CheckGuess(guess); fbstr := fb.String()
19. var before []Word
20. if len(req.History) == 0 {
21:     before = AllowedSolutions
22: } else {
23:     firstKey := req.History[0].Guess + "|" + req.History[0].Feedback
24:     mu.RLock(); beforeCached, ok := cache[firstKey]; mu.RUnlock()
25:     if ok {
26:         before = append([]Word{}, beforeCached...)
27:     } else {
28:         before = FilterSolutions(req.History[:1])
29:         mu.Lock(); cache[firstKey] = append([]Word{}, before...); mu.Unlock()
30:     }
31: }
32: newHistory := append(req.History, Turn{req.Proposed, fbstr})
33: after := FilterSolutions(newHistory)
34: reduction := 0.0; if len(before) > 0 { reduction = 1.0 - float64(len(after))/float64(len(before)) }
35: turns := len(req.History) + 1
36: status := \"ongoing\"
37: if fb == Feedback{5*Green} { status = \"won\" }
38: else if turns >= 6 { status = \"lost\" }
39: resp := Response{Status: status, Valid: valid, Before: len(before), After: len(after), Reduction: reduction, Feedback: fbstr}
40: json.NewEncoder(w).Encode(resp)
41: If len(req.History)==0 && len(after)>0 { // cache after first
42:     key := req.Proposed + "|" + fbstr
43:     mu.Lock()
44:     // check memory rough: len([]Word)*6*1024 <10e9 etc, skip for now or simple count keys
45:     cache[key] = append([]Word{}, after...)
46:     mu.Unlock()
47: }
48. http.HandleFunc(\"/api/evaluate\", handler)
49. log.Fatal(http.ListenAndServe(\":9111\", nil))
50. Add tests for API: create handler_test.go in pkg or main_test.go
51. Run go test ./...
52. Implement memory limit for cache: track total size, evict LRU if >8GB.
