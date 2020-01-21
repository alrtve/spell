package linear

type StochasticLearningAlgorithm struct {
}

func (learner *StochasticLearningAlgorithm) Learn(vectorSystems []*VectorSystem) *Vector {
	var bestVector *Vector
	if len(vectorSystems) > 0 && vectorSystems[0] != nil &&
		len(vectorSystems[0].Vectors) > 0 && vectorSystems[0].Vectors[0] != nil {
		var (
			vector *Vector
			tries = 10000
			bestScore = 0
			currentScore = 0
		)
		for i := 0; i < tries; i++ {
			vector = RandomVector(vectorSystems[0].Vectors[0].Len())
			for _, vectorSystem := range vectorSystems {
				if vectorSystem.IsSatisfied(vector) {
					currentScore += 1
				}
			}
			if currentScore > bestScore {
				bestScore = currentScore
				bestVector = vector
			}
		}
	}
	return bestVector
}


