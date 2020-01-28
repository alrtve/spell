package spell

type EditAction uint32

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
	AffectedLen  int
	Weight       float64
	MisspellLens []int
}

type EditVariance struct {
	MinTermLen    int
	PossibleEdits [][]EditAction
}

type OperationAffectedChange struct {
	Weight    float64
	InputLens map[int]bool
}

type EditorialPrescription struct {
	Froms   []rune
	Tos     []rune
	Actions []EditAction
}

type Suggestion struct {
	Term         string
	Distance     int
	Score        float64
	Prescription *EditorialPrescription
}

type Misspell struct {
	Term      string
	Misspells []string
}

type LearningTerm struct {
	Term        string
	Misspell    string
	Suggestions []Suggestion
}

type Scorer interface {
	GetDifference(prescription *EditorialPrescription) float64
	Learn(learningData []LearningTerm)
}

type ScoreModel interface {
	Compare(a *Suggestion, b *Suggestion) float64
}

type Learner interface {
	Learn(learningData []*LearningTerm) ScoreModel
}