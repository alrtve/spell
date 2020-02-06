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


type LearnProgress struct {
	VectorSystemsCount int
	Step               int
	BestScore          int
	RelaxingCount      int
}

type Learner struct {
	*scorer.Vectoriser
	learnProgress LearnProgress
}

func (learner *Learner) Learn(learningData []*spell.LearningTerm) spell.ScoreModel {
	learner.learnProgress = LearnProgress{}
	vectorSystems := make([]*VectorSystem, 0, len(learningData))
	for _, learningTerm := range learningData {
		vectorSystem := learner.vectorSystem(learningTerm)
		if vectorSystem != nil {
			vectorSystems = append(vectorSystems, vectorSystem)
		}
	}
	learner.learnProgress.VectorSystemsCount = len(vectorSystems)

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
		triesLoop:
		for i := 0; i < tries; i++ {
			learner.learnProgress.Step = i + 1
			vector = scorer.RandomVector(vectorSystems[0].Vectors[0].Len())
			currentScore := learner.score(vectorSystems, vector)
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

				if currentScore > bestScore {
					bestScore = currentScore
					bestVector = vector
					learner.learnProgress.BestScore = bestScore
				}
			}
			if bestScore > prevBestScore {
				relaxingCount = 0
			} else {
				relaxingCount++
				learner.learnProgress.RelaxingCount = relaxingCount
				if relaxingCount == maxRelaxCount {
					break triesLoop
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

func (learner *Learner) LearnProgress() string {
	learnProgress := learner.learnProgress;
	return fmt.Sprintf("Step %d. Best score: %d. VS count: %d. Relaxing Count: %d", learnProgress.Step, learnProgress.BestScore, learnProgress.VectorSystemsCount, learnProgress.RelaxingCount)
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


func (learner *Learner) vectorSystem(a *spell.LearningTerm) *VectorSystem {
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