package spell


const (
	Insert        EditAction = 0
	Delete        EditAction = 1
	Replace       EditAction = 2
	Match         EditAction = 3
	Transposition EditAction = 4
	Duplicate     EditAction = 5
	MissDouble    EditAction = 6
	Triplet       EditAction = 7
)

type OperationWeight struct {
	AffectedLen int
	Weight float64
	MisspellLens []int
}

type EditVariance struct {
	MinTermLen int
	PossibleEdits [][]EditAction
}

type OperationAffectedChange struct {
	Weight float64
	InputLens map[int]bool
}

type Suggestion struct {
	Term         string
	Distance     int
	Score        float64
	Prescription *EditorialPrescription
}

var OperationsAffects = map[EditAction]int{
	Insert: 0,
	Delete: 1,
	Replace: 1,
	Transposition: 2,
	Duplicate: 0,
	MissDouble: 1,
	Triplet: 3,
}