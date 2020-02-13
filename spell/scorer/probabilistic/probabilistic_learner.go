package probabilistic

import (
	"fmt"
	"spell"
	"spell/scorer"
)

type LearnProgress struct {
	ProcessedTerms int
}


type Learner struct {
	*scorer.Vectoriser
	learnProgress LearnProgress
}

func (learner *Learner) Learn(learningData []*spell.LearningTerm) spell.ScoreModel {
	learner.learnProgress = LearnProgress{}
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
	learner.learnProgress.ProcessedTerms = len(learningData)
	return &Scorer{
		Weights:    weights,
		Vectoriser: learner.Vectoriser,
	}
}

func (learner *Learner) LearnProgress() string {
	learnProgress := learner.learnProgress;
	return fmt.Sprintf("Processed terms: %d", learnProgress.ProcessedTerms)
}
