package probabilistic

import (
	"math"
	"spell"
	"spell/scorer"
)

var eps = 0.000001

type Scorer struct {
	Weights *scorer.Vector
	*scorer.Vectoriser
}

func (scoring *Scorer) Compare(a *spell.Suggestion, b *spell.Suggestion) float64 {
	largeValue := math.MaxFloat64
	if a.Prescription == nil {
		return -largeValue
	}
	if b.Prescription == nil {
		return largeValue
	}

	vectorA := scoring.Vectorize(a.Prescription)
	vectorB := scoring.Vectorize(b.Prescription)
	a.Score = a.Count * probabilisticMul(vectorA, scoring.Weights)
	b.Score = b.Count * probabilisticMul(vectorB, scoring.Weights)
	return b.Score - a.Score
}

func (scoring *Scorer) Score(a *spell.Suggestion) float64 {
	vectorA := scoring.Vectorize(a.Prescription)
	score := a.Count * probabilisticMul(vectorA, scoring.Weights)
	return 1 - score
}

func probabilisticMul(a, w *scorer.Vector) float64 {
	result := 1.0
	for i, xs := range a.Xs {
		if math.Abs(xs) > eps {
			result *= math.Pow(w.Xs[i], xs)
		}
	}
	return result
}

