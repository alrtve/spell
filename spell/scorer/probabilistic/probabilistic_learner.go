package probabilistic

import (
	"spell"
	"spell/scorer"
)

type Learner struct {
	*scorer.Vectoriser
}

func InitLearner(vectoriser *scorer.Vectoriser) *Learner {
	return &Learner{vectoriser}
}

func (learner *Learner) Learn(learningData []*spell.LearningTerm) spell.ScoreModel {
	var weights *scorer.Vector = nil
	for _, learningTerm := range learningData {
		for _, suggestion := range learningTerm.Suggestions{
			if suggestion.Term == learningTerm.Term {
				vector := learner.Vectorize(suggestion.Prescription)
				if weights == nil {
					weights = vector
				} else {
					weights = weights.Add(vector)
				}
			}
		}
	}
	totalW := 0.0
	for _, xs := range weights.Xs {
		totalW += xs
	}
	for i := range weights.Xs {
		weights.Xs[i] /= totalW
	}
	return &ScoreModel{
		Weights:    weights,
		Vectoriser: learner.Vectoriser,
	}
}
