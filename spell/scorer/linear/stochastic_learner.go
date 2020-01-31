package linear

import (
	"fmt"
	"spell"
	"spell/scorer"
)

type vectorScore struct {
	*scorer.Vector
	score float64

}

type Learner struct {
	*scorer.Vectoriser
}

func (learner *Learner) Learn(learningData []*spell.LearningTerm) spell.ScoreModel {
	vectorSystems := make([]*VectorSystem, 0, len(learningData))
	for _, learningTerm := range learningData {
		vectorSystem := learner.GetVectorSystem(learningTerm)
		if vectorSystem != nil {
			vectorSystems = append(vectorSystems, vectorSystem)
		}
	}

	fmt.Println(len(vectorSystems))
	var bestVector *scorer.Vector
	if len(vectorSystems) > 0 && vectorSystems[0] != nil &&
		len(vectorSystems[0].Vectors) > 0 && vectorSystems[0].Vectors[0] != nil {
		var (
			vector               *scorer.Vector
			tries                = 100
			maxRelaxCount        = 10
			relaxingCount        = 0
			maxDirectSearchCount = 1000
			bestScore            = 0
			prevBestScore        = 0
		)
		tries_loop:
		for i := 0; i < tries; i++ {
			vector = scorer.RandomVector(vectorSystems[0].Vectors[0].Len())
			currentScore := learner.score(vectorSystems, vector)
			fmt.Println(i)
			for k := 0; k < maxDirectSearchCount; k++ {
				vVectors := vector.Variate(0, 1, 0.1)
				valuableVectors := make([]vectorScore, 0, len(vVectors))
				for _, v := range vVectors {
					vScore := learner.score(vectorSystems, v)
					if vScore > currentScore {
						valuableVectors = append(valuableVectors, vectorScore{
							Vector: v,
							score:  float64(vScore),
						})
					}
				}
				if len(valuableVectors) > 0 {
					maxScore := 0.0
					for _, vs := range valuableVectors {
						if vs.score > maxScore {
							maxScore = vs.score
						}
					}
					nextVector := vector.Clone()
					for _, vs := range valuableVectors {
						nextVector = nextVector.MoveToward(vs.Vector, vs.score / maxScore)
					}
					currentScore = learner.score(vectorSystems, nextVector)
					vector = nextVector
				} else {
					break
				}

				fmt.Println(currentScore)
				fmt.Println(bestScore)
				fmt.Println("")
				if currentScore > bestScore {
					fmt.Println(bestScore)
					bestScore = currentScore
					bestVector = vector
				}
			}
			if bestScore > prevBestScore {
				relaxingCount = 0
			} else {
				relaxingCount++
				if relaxingCount == maxRelaxCount {
					break tries_loop
				}
			}
			prevBestScore = bestScore
		}
	}
	return &ScoreModel{
		Weights:    bestVector,
		Vectoriser: scorer.InitVectoriser(),
	}
}

func (learner *Learner) score(vectorSystems []*VectorSystem, vector *scorer.Vector) int {
	currentScore := 0
	for _, vectorSystem := range vectorSystems {
		if vectorSystem.IsSatisfied(vector) {
			currentScore += 1
		}
	}
	return currentScore
}


func (learner *Learner) GetVectorSystem(a *spell.LearningTerm) *VectorSystem {
	baseVector := (*scorer.Vector)(nil)
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