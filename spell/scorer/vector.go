package scorer

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
)

var eps = 0.000001

type Vector struct {
	Xs []float64
}

func InitVector(length int) *Vector {
	return &Vector{
		Xs: make([]float64, length),
	}
}

func RandomVector(length int) *Vector {
	result := InitVector(length)
	for i := range result.Xs {
		result.Xs[i] = rand.Float64()
	}
	return result
}


func (a *Vector) Len() int {
	return len(a.Xs)
}

func (a *Vector) Sub(b *Vector) *Vector {
	result := InitVector(len(a.Xs))
	for i, xa := range a.Xs {
		result.Xs[i] = xa - b.Xs[i]
	}
	return result
}

func (a *Vector) Add(b *Vector) *Vector {
	result := InitVector(len(a.Xs))
	for i, xa := range a.Xs {
		result.Xs[i] = xa + b.Xs[i]
	}
	return result
}

func (a *Vector) EqualTo(b *Vector) bool {
	for i := range a.Xs {
		if math.Abs(a.Xs[i]-b.Xs[i]) > eps {
			return false
		}
	}
	return true
}

func (a *Vector) Clone() *Vector {
	result := InitVector(len(a.Xs))
	for i, xs := range a.Xs {
		result.Xs[i] = xs
	}
	return result
}

func (a *Vector) IsZero() bool {
	for _, val := range a.Xs {
		if math.Abs(val) > eps {
			return false
		}
	}
	return true
}

func (a *Vector) Variate(min, max, d float64) []*Vector {
	result := make([]*Vector, 0, 2 * len(a.Xs))
	for i := range a.Xs {
		v := a.Clone()
		v.Xs[i] += d
		if v.Xs[i] <= max {
			result = append(result, v)
		}
		v = a.Clone()
		v.Xs[i] -= d
		if v.Xs[i] >= min {
			result = append(result, v)
		}
	}
	return result
}

func (a *Vector) MoveToward(vector *Vector, l float64) *Vector{
	result := a.Clone()
	for i := range a.Xs {
		result.Xs[i] += l * (vector.Xs[i] - a.Xs[i])
	}
	return result
}

func (inequality *Vector) Dump() {
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

func (a *Vector) IsSatisfied(wights *Vector) bool {
	val := 0.0
	for i := range a.Xs {
		val += wights.Xs[i] * a.Xs[i]
	}
	return val > 0
}

func (a *Vector) ScalarMul(wights *Vector) float64 {
	val := 0.0
	for i := range a.Xs {
		val += wights.Xs[i] * a.Xs[i]
	}
	return val
}
