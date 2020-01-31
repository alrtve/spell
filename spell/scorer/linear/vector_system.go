package linear

import (
	"fmt"
	"math"
	"spell/scorer"
	"strings"
)

var eps = 0.000001

type VectorSystem struct {
	Vectors []*scorer.Vector
}

func InitVectorSystem() *VectorSystem {
	return &VectorSystem{
		Vectors: make([]*scorer.Vector, 0),
	}
}

func (system *VectorSystem) Add(inequality *scorer.Vector) {
	system.Vectors = append(system.Vectors, inequality)
}

func (system *VectorSystem) Normalize() {
	inequalities := make([]*scorer.Vector, 0, len(system.Vectors))
	for _, inequality := range system.Vectors {
		if inequality.IsZero() {
			continue
		}
		isUniq := true
		for _, uniqInequality := range inequalities {
			if inequality.EqualTo(uniqInequality) {
				isUniq = false
				break
			}
		}
		if isUniq {
			inequalities = append(inequalities, inequality)
		}
	}
	system.Vectors = inequalities
}

func (system *VectorSystem) IsSimple() bool {
	for _, inequality := range system.Vectors {
		if !inequality.IsSimple() {
			return false
		}
	}
	return true
}

func (system *VectorSystem) Dump() {
	for _, inequality := range system.Vectors {
		displayValues := make([]string, 0, 20)
		for i, val := range inequality.Xs {
			if math.Abs(val) > eps {
				valStr := fmt.Sprintf("%0.1f*x%d", math.Abs(val), i)
				if len(displayValues) > 0 || val < 0 {
					sign := "+"
					if val < 0 {
						sign = "-"
					}
					displayValues = append(displayValues, sign)
				}
				displayValues = append(displayValues, valStr)
			}
		}
		if len(displayValues) > 0 {
			displayValues = append(displayValues, "> 0")
			displayStr := strings.Join(displayValues, " ")
			fmt.Printf("%s\n", displayStr)
		}
	}
}

func (system *VectorSystem) Gn(vector *scorer.Vector) float64 {
	min, resultApplicable := -1000.0, false
	for i := 0; i < len(system.Vectors); i++ {
		val, applicable := system.Vectors[i].Gn(vector)
		min = math.Max(min, val)
		resultApplicable = resultApplicable || applicable
	}
	if resultApplicable {
		return min
	}
	return 2
}

func (system *VectorSystem) IsSatisfied(vector *scorer.Vector) bool {
	for _, inequality := range system.Vectors {
		if !inequality.IsSatisfied(vector) {
			return false
		}
	}
	return true
}

func (system *VectorSystem) Compare(a *scorer.Vector, b *scorer.Vector) int {
	result := 0.0
	for i := 0; i < len(system.Vectors); i++ {
		aVal, aApplicable := system.Vectors[i].Gn(a)
		bVal, bApplicable := system.Vectors[i].Gn(b)
		if aApplicable && bApplicable {
			result += aVal - bVal
		}
	}
	if math.Abs(result) < eps {
		return 0
	}
	if result < 0 {
		return -1
	}
	return 1
}
