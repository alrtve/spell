package linear

import "spell"

type LinearScoreModel struct {
	Weights *Vector
	*Vectoriser
}

func (scoreModel *LinearScoreModel) Compare(a *spell.Suggestion, b *spell.Suggestion) float64 {
	largeValue := 100.0
	if a.Prescription == nil {
		return -largeValue
	}
	if b.Prescription == nil {
		return largeValue
	}

	vectorA := scoreModel.Vectorize(a.Prescription)
	vectorB := scoreModel.Vectorize(b.Prescription)
	if vectorA.Sub(vectorB).IsSatisfied(scoreModel.Weights) {
		return 1
	}
	return -1
}


func (scoreModel *LinearScoreModel) GetVectorSystem(a *spell.LearningTerm) *VectorSystem {
	baseVector := (*Vector)(nil)
	for _, suggestion := range a.Suggestions {
		if suggestion.Term == a.Term {
			baseVector = scoreModel.Vectorize(suggestion.Prescription)
			break
		}
	}
	if baseVector == nil {
		return nil
	}

	vectorSystem := InitVectorSystem()
	for _, suggestion := range a.Suggestions {
		if suggestion.Term != a.Term && a.Misspell != suggestion.Term {
			vector := scoreModel.Vectorize(suggestion.Prescription)
			vector = vector.Sub(baseVector)
			vectorSystem.Add(vector)
			break
		}
	}
	vectorSystem.Normalize()
	return vectorSystem
}
