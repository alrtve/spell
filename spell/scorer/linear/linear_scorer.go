package linear

import (
	"spell"
)

type Scorer struct {
	Weights    *Vector
	Vectorizer SimpleVectorizer
	Learner    LearningAlgorithm
}

func (linear *Scorer) GetDifference(prescription *spell.EditorialPrescription) float64 {
	difference := 0.0
	vector := linear.Vectorizer.Vectorize(prescription)
	for i := range vector.Xs {
		difference += linear.Weights.Xs[i] * vector.Xs[i]
	}
	return difference
}

func (scorer *Scorer) Learn(learningData []spell.LearningData) {
	byTerms := make(map[string][]spell.LearningData)
	vectorSystems := make([]*VectorSystem, 8)
	for _, learning := range learningData {
		if _, ok := byTerms[learning.Term]; !ok {
			byTerms[learning.Term] = make([]spell.LearningData, 4)
		}
		byTerms[learning.Term] = append(byTerms[learning.Term], learning)
	}

	for term, lernings := range byTerms {
		for _, learning := range lernings {
			var (
				baseVector   *Vector
				vectorSystem *VectorSystem
			)
			for _, suggestion := range learning.Suggestions {
				if suggestion.Term == term {
					baseVector = scorer.Vectorizer.Vectorize(suggestion.Prescription)
					break
				}
			}
			if baseVector != nil {
				vectorSystem = InitVectorSystem()
				for _, suggestion := range learning.Suggestions {
					if suggestion.Term != term {
						vector := scorer.Vectorizer.Vectorize(suggestion.Prescription)
						vector = vector.Sub(baseVector)
						vectorSystem.Add(vector)
						break
					}
				}
				vectorSystem.Normalize()
				vectorSystems = append(vectorSystems, vectorSystem)
			}
		}
	}
	scorer.Weights = scorer.Learner.Learn(vectorSystems)
}
