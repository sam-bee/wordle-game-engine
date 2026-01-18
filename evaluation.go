package wordlegameengine

func matchesAllFeedback(g *Game, candidate Word) bool {
	candidateSolution := Solution(candidate)

	for i, guess := range g.Guesses {
		expectedFeedback := g.Feedbacks[i]
		actualFeedback := candidateSolution.CheckGuess(guess)

		if actualFeedback != expectedFeedback {
			return false
		}
	}

	return true
}
