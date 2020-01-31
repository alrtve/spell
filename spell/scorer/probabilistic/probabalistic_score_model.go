package probabilistic

import (
	"math"
	"spell"
	"spell/scorer"
)

var eps = 0.000001

type ScoreModel struct {
	Weights *scorer.Vector
	*scorer.Vectoriser
}

func (scoreModel *ScoreModel) Compare(a *spell.Suggestion, b *spell.Suggestion) float64 {
	largeValue := math.MaxFloat64
	if a.Prescription == nil {
		return -largeValue
	}
	if b.Prescription == nil {
		return largeValue
	}

	vectorA := scoreModel.Vectorize(a.Prescription)
	vectorB := scoreModel.Vectorize(b.Prescription)
	a.Score = a.Count * probabalisticMul(vectorA, scoreModel.Weights)
	b.Score = b.Count * probabalisticMul(vectorB, scoreModel.Weights)
	return b.Score - a.Score
}

func probabalisticMul(a, w *scorer.Vector) float64 {
	result := 1.0
	for i, xs := range a.Xs {
		if math.Abs(xs) > eps {
			result *= math.Pow(w.Xs[i], xs)
		}
	}
	return result
}

