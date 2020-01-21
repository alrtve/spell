package spell

type DistanceMeasurer struct {
	p          [][] int
	e          [][] EditAction
	maxRows    int
	maxColumns int
}

func NewDistanceMeasurer() *DistanceMeasurer {
	return &DistanceMeasurer{}
}

func (measurer *DistanceMeasurer) Distance(a, b string, calcEditorialPrescription bool) (distance int, editorialPrescription *EditorialPrescription) {
	var (
		ar = []rune(a)
		br = []rune(b)
		la = len(ar) + 1
		lb = len(br) + 1
	)
	measurer.ensureSizes(la, lb)
	for i := 0; i < measurer.maxRows; i++ {
		measurer.p[i][0] = i
		measurer.e[i][0] = Insert
	}
	for j := 0; j < measurer.maxColumns; j++ {
		measurer.p[0][j] = j
		measurer.e[0][j] = Delete
	}

	maxDistance := la + lb
	if lb > maxDistance {
		maxDistance = la
	}
	for i := 1; i < lb; i++ {
		for j := 1; j < la; j++ {
			del := measurer.p[i][j-1] + 1
			ins := measurer.p[i-1][j] + 1
			repl := measurer.p[i-1][j-1] + 1
			match := maxDistance
			if br[i-1] == ar[j-1] {
				match = measurer.p[i-1][j-1]
			}

			transpose := maxDistance
			if i >= 2 && j >= 2 {
				if br[i-1] == ar[j-2] && br[i-2] == ar[j-1] {
					transpose = measurer.p[i-2][j-2] + 1
				}
			}

			doupl := maxDistance
			if i >= 2 && ar[j-1] == br[i-1] && ar[j-1] == br[i-2] {
				doupl = measurer.p[i-1][j] + 1
			}
			missDoupl := maxDistance
			if j >= 2 && ar[j-1] == br[i-1] && ar[j-2] == br[i-1] {
				missDoupl = measurer.p[i][j-1] + 1
			}

			triplet := maxDistance
			// for triplets chars must be pairwise different
			if i >= 3 && j >= 3 && a[j-3] != a[j-2] && a[j-3] != a[j-1] && a[j-2] != a[j-1] {
				if ar[j-3] == br[i-1] && ar[j-2] == br[i-3] && ar[j-1] == br[i-2] {
					triplet = measurer.p[i-3][j-3] + 1
				} else if ar[j-3] == br[i-2] && ar[j-2] == br[i-1] && ar[j-1] == br[i-3] {
					triplet = measurer.p[i-3][j-3] + 1
				}
			}

			min := del
			action := Delete
			if ins < min {
				min = ins
				action = Insert
			}
			if match <= min {
				min = match
				action = Match
			}
			if repl < min {
				min = repl
				action = Replace
			}
			if doupl <= min {
				min = doupl
				action = Duplicate
			}
			if missDoupl <= min {
				min = missDoupl
				action = MissDouble
			}
			if transpose <= min {
				min = transpose
				action = Transposition
			}
			// min < del -- special case when no actions are applicable
			if triplet <= min {
				min = triplet
				action = Triplet
			}
			measurer.p[i][j] = min
			measurer.e[i][j] = action
		}
	}
	distance = measurer.p[lb-1][la-1]
	if calcEditorialPrescription {
		editorialPrescription = measurer.getEditorialPrescription(ar, br, la, lb)
	}
	return
}

func (measurer *DistanceMeasurer) getEditorialPrescription(ar, br []rune, la, lb int) *EditorialPrescription {
	var (
		i  = lb - 1
		j  = la - 1
		ia = i + j
	)
	var (
		actions = make([]EditAction, ia)
		froms   = make([]rune, ia)
		tos     = make([]rune, ia)
	)
	ia--
	for i > 0 || j > 0 {
		action := measurer.e[i][j]
		actions[ia] = action
		froms[ia] = rune(0)
		tos[ia] = rune(0)

		switch action {
		case Delete:
			fallthrough
		case MissDouble:
			froms[ia] = ar[j-1]
			j--
		case Insert:
			fallthrough
		case Duplicate:
			tos[ia] = br[i-1]
			i--
		case Match:
			fallthrough
		case Replace:
			froms[ia] = ar[j-1]
			tos[ia] = br[i-1]
			i--
			j--
		case Transposition:
			froms[ia] = ar[j-1]
			froms[ia-1] = ar[j-2]
			tos[ia] = br[i-1]
			tos[ia-1] = br[i-2]

			actions[ia-1] = Transposition
			ia -= 1
			i -= 2
			j -= 2
		case Triplet:
			froms[ia] = ar[j-1]
			froms[ia-1] = ar[j-2]
			froms[ia-2] = ar[j-3]
			tos[ia] = br[i-1]
			tos[ia-1] = br[i-2]
			tos[ia-2] = br[i-3]

			actions[ia-1] = Triplet
			actions[ia-2] = Triplet
			ia -= 2
			i -= 3
			j -= 3
		}
		ia--
	}
	ia++
	froms = froms[ia:]
	tos = tos[ia:]
	actions = actions[ia:]

	return &EditorialPrescription{
		Froms:   froms,
		Tos:     tos,
		Actions: actions,
	}
}

func (measurer *DistanceMeasurer) ensureSizes(la, lb int) {
	if la > measurer.maxColumns || lb > measurer.maxRows {
		measurer.p = make([][]int, lb)
		measurer.e = make([][]EditAction, lb)
		for i := 0; i < lb; i++ {
			measurer.p[i] = make([]int, la)
			measurer.e[i] = make([]EditAction, la)
		}
		measurer.maxRows = lb
		measurer.maxColumns = la
	}
}
