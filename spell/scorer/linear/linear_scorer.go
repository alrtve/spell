package linear

import (
	"spell"
	"spell/scorer"
)

type Scorer struct {
	Weights *scorer.Vector
	*scorer.Vectoriser
}

func (scoring *Scorer) Compare(a *spell.Suggestion, b *spell.Suggestion) float64 {
	largeValue := 100.0
	if a.Prescription == nil {
		return -largeValue
	}
	if b.Prescription == nil {
		return largeValue
	}

	vectorA := scoring.Vectorize(a.Prescription)
	vectorB := scoring.Vectorize(b.Prescription)
	a.Score = vectorA.ScalarMul(scoring.Weights)
	b.Score = vectorB.ScalarMul(scoring.Weights)
	if vectorA.Sub(vectorB).IsSatisfied(scoring.Weights) {
		return 1
	}
	return -1
}

func (scoring *Scorer) Score(a *spell.Suggestion) float64 {
	largeValue := -100.0
	if a.Prescription == nil {
		return largeValue
	}

	vectorA := scoring.Vectorize(a.Prescription)
	score := vectorA.ScalarMul(scoring.Weights)
	return score
}

func (scoring *Scorer) GetVectorSystem(a *spell.LearningTerm) *VectorSystem {
	baseVector := (*scorer.Vector)(nil)
	for _, suggestion := range a.Suggestions {
		if suggestion.Term == a.Term {
			baseVector = scoring.Vectorize(suggestion.Prescription)
			break
		}
	}
	if baseVector == nil {
		return nil
	}

	vectorSystem := InitVectorSystem()
	for _, suggestion := range a.Suggestions {
		if suggestion.Term != a.Term && a.Misspell != suggestion.Term {
			vector := scoring.Vectorize(suggestion.Prescription)
			vector = vector.Sub(baseVector)
			vectorSystem.Add(vector)
			break
		}
	}
	vectorSystem.Normalize()
	return vectorSystem
}
