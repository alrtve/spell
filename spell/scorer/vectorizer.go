package scorer

import (
	"math"
	"spell"
)

const (
	RFirst         = 0
	RConsonant     = 1
	RDistance1     = 2
	RDistance2     = 3
	RDistance3     = 4
	RDistance4     = 5
	RDistanceOther = 6
	DFirst         = 7
	DLast          = 8
	DMiddle        = 9
	IFirst         = 10
	IMiddle        = 11
	T              = 12
	U              = 13
	P              = 14
	J              = 15
)

type Vectoriser struct {
}

func InitVectoriser() *Vectoriser  {
	return new (Vectoriser)
}

func (maker *Vectoriser) Vectorize(prescription *spell.EditorialPrescription) *Vector {
	inequality := InitVector(16)
	for i := 0; i < len(prescription.Actions); {
		action := prescription.Actions[i]
		from := prescription.Froms[i]
		to := prescription.Tos[i]
		switch action {
		case spell.Replace:
			if i == 0 {
				inequality.Xs[RFirst] += 1
			} else {
				consonantDistance := maker.consonantDistance(from, to)
				if consonantDistance >= 0 {
					inequality.Xs[RConsonant] += 1
				} else {
					qwertyDistance := maker.qwertyDistance(from, to)
					if qwertyDistance >= 0 {
						switch true {
						case qwertyDistance <= 1+eps:
							inequality.Xs[RDistance1] += 1
						case qwertyDistance <= 2+eps:
							inequality.Xs[RDistance2] += 1
						case qwertyDistance <= 3+eps:
							inequality.Xs[RDistance3] += 1
						case qwertyDistance <= 4+eps:
							inequality.Xs[RDistance4] += 1
						default:
							inequality.Xs[RDistanceOther] += 1
						}
					}
				}
			}
		case spell.Delete:
			switch true {
			case i == 0:
				inequality.Xs[DFirst] += 1
			case i == len(prescription.Actions)-1:
				inequality.Xs[DLast] += 1
			default:
				inequality.Xs[DMiddle] += 1
			}
		case spell.Insert:
			if i == 0 {
				inequality.Xs[IFirst] += 1
			} else {
				inequality.Xs[IMiddle] += 1
			}
		case spell.Transposition:
			inequality.Xs[T] += 1
			i++
		case spell.Duplicate:
			inequality.Xs[P] += 1
			i++
		case spell.MissDouble:
			inequality.Xs[U] += 1
			i++
		case spell.Triplet:
			inequality.Xs[J] += 1
			i += 2
		}
		i++
	}
	return inequality
}

func (maker *Vectoriser) consonantDistance(from, to rune) float64 {
	var sounds = [][]rune{
		{'d', 't'},
		{'u', 'y'},
		{'u', 'a'},
		{'v', 'w'},
		{'j', 'g'},
		{'c', 'k'},
		{'s', 'z'},
	}
	for _, sound := range sounds {
		if from == sound[0] {
			for i := 1; i < len(sound); i++ {
				if sound[i] == to {
					return 1
				}
			}
		}
	}
	return -1
}

func (maker *Vectoriser) qwertyDistance(from, to rune) float64 {
	locationFrom, locationFromExists := distances[from]
	locationTo, locationToExists := distances[from]

	if locationFromExists && locationToExists {
		return math.Abs(locationTo.X-locationFrom.X) + math.Abs(locationTo.Y-locationFrom.Y)
	}
	return -1
}

type KeyLocation struct {
	X float64
	Y float64
}

var (
	offset1 = 0.5
	offset2 = 0.9
	offset3 = 1.3
)

var distances = map[rune]KeyLocation{
	// number row
	'0': KeyLocation{0, 0},
	'1': KeyLocation{1, 0},
	'2': KeyLocation{2, 0},
	'3': KeyLocation{3, 0},
	'4': KeyLocation{4, 0},
	'5': KeyLocation{5, 0},
	'6': KeyLocation{6, 0},
	'7': KeyLocation{7, 0},
	'8': KeyLocation{8, 0},
	'9': KeyLocation{9, 0},

	// first row
	'q': KeyLocation{0 + offset1, 1},
	'w': KeyLocation{1 + offset1, 1},
	'e': KeyLocation{2 + offset1, 1},
	'r': KeyLocation{3 + offset1, 1},
	't': KeyLocation{4 + offset1, 1},
	'y': KeyLocation{5 + offset1, 1},
	'u': KeyLocation{6 + offset1, 1},
	'i': KeyLocation{7 + offset1, 1},
	'o': KeyLocation{8 + offset1, 1},
	'p': KeyLocation{9 + offset1, 1},

	// second row
	'a': KeyLocation{0 + offset2, 2},
	's': KeyLocation{1 + offset2, 2},
	'd': KeyLocation{2 + offset2, 2},
	'f': KeyLocation{3 + offset2, 2},
	'g': KeyLocation{4 + offset2, 2},
	'h': KeyLocation{5 + offset2, 2},
	'j': KeyLocation{6 + offset2, 2},
	'k': KeyLocation{7 + offset2, 2},
	'l': KeyLocation{8 + offset2, 2},

	// third row
	'z': KeyLocation{0 + offset3, 3},
	'x': KeyLocation{1 + offset3, 3},
	'c': KeyLocation{2 + offset3, 3},
	'v': KeyLocation{3 + offset3, 3},
	'b': KeyLocation{4 + offset3, 3},
	'n': KeyLocation{5 + offset3, 3},
	'm': KeyLocation{6 + offset3, 3},
}
