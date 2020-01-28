package linear

import (
	"spell"
)

type StochasticLearner struct {
	*Vectoriser
}

func (learner *StochasticLearner) Learn(learningData []*spell.LearningTerm) spell.ScoreModel {
	vectorSystems := make([]*VectorSystem, 0, len(learningData))
	for _, learningTerm := range learningData {
		vectorSystem := learner.GetVectorSystem(learningTerm)
		if vectorSystem != nil {
			vectorSystems = append(vectorSystems, vectorSystem)
		}
	}

	var bestVector *Vector
	if len(vectorSystems) > 0 && vectorSystems[0] != nil &&
		len(vectorSystems[0].Vectors) > 0 && vectorSystems[0].Vectors[0] != nil {
		var (
			vector       *Vector
			tries        = 10000
			bestScore    = 0
		)
		for i := 0; i < tries; i++ {
			currentScore := 0
			vector = RandomVector(vectorSystems[0].Vectors[0].Len())
			for _, vectorSystem := range vectorSystems {
				if vectorSystem.IsSatisfied(vector) {
					currentScore += 1
				}
			}
			if currentScore > bestScore {
				bestScore = currentScore
				bestVector = vector
				//fmt.Println(currentScore)
				//fmt.Println(vector.Xs)
			}
		}
	}
	return &LinearScoreModel{
		Weights:bestVector,
		Vectoriser: InitVectoriser(),
	}
}


func (learner *StochasticLearner) GetVectorSystem(a *spell.LearningTerm) *VectorSystem {
	baseVector := (*Vector)(nil)
	for _, suggestion := range a.Suggestions {
		if suggestion.Term == a.Term {
			baseVector = learner.Vectorize(suggestion.Prescription)
			break
		}
	}
	if baseVector == nil {
		return nil
	}

	vectorSystem := InitVectorSystem()
	for _, suggestion := range a.Suggestions {
		if suggestion.Term != a.Term && a.Misspell != suggestion.Term {
			vector := learner.Vectorize(suggestion.Prescription)
			vector = vector.Sub(baseVector)
			vectorSystem.Add(vector)
		}
	}
	vectorSystem.Normalize()
	if len(vectorSystem.Vectors) > 0 {
		return vectorSystem
	}
	return nil
}