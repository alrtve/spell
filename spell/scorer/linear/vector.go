package linear

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

func (a *Vector) EqualTo(b *Vector) bool {
	for i := range a.Xs {
		if math.Abs(a.Xs[i]-b.Xs[i]) > eps {
			return false
		}
	}
	return true
}

func (a *Vector) Vector(initialValue float64) *Vector {
	result := InitVector(len(a.Xs))
	for i := range a.Xs {
		result.Xs[i] = initialValue
	}
	return result
}

func (a *Vector) RandomVector() *Vector {
	result := InitVector(len(a.Xs))
	for i := range a.Xs {
		result.Xs[i] = rand.Float64()
	}
	return result
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

func (a *Vector) IsSimple() bool {
	for _, val := range a.Xs {
		if math.Abs(val) > 1+eps {
			return false
		}
	}
	return true
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

// goodness
func (a *Vector) Gn(vector *Vector) (float64, bool) {
	val := 0.0
	significantCount := 0
	minVal := 1.0
	for i, xs := range a.Xs {
		val += xs * vector.Xs[i]
		if math.Abs(xs) < eps {
			significantCount += 1
		} else {
			minVal = math.Min(minVal, math.Abs(xs))
		}
	}
	if significantCount <= 1 {
		return 1, false
	}
	if math.Abs(val) < eps {
		return 1, true
	}
	if val > 0 {
		return 1 / (1 + val/minVal), true
	}
	return 1 - val/minVal/5, true
}

func (a *Vector) IsSatisfied(wights *Vector) bool {
	val := 0.0
	for i := range a.Xs {
		val += wights.Xs[i] * a.Xs[i]
	}
	return val > 0
}
