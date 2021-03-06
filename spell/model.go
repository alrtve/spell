package spell

import (
	"regexp"
	"runtime"
	"sort"
	"strings"
	"sync"
	"unicode/utf8"
)

/**
	Model represents misspell corrector structure
 */
type Model struct {
	Terms     []string
	TermsDict map[string]int
	Index     map[string][][]int
	IndexTail map[string]map[string]bool

	Affects      [][][]int // inputlen -> edit len -> term lens
	KnownAffects []bool
	Depth int
	IndexSplitLen int

	TermsCounts []float64
	TotalTerms  float64

	distanceMeasurers  []*DistanceMeasurer `binary:"-"`
	currentMeasurerId  int                 `binary:"-"`
	measurersSemaphore chan bool           `binary:"-"`
	measurersMutex     sync.Mutex          `binary:"-"`
}

func InitModel() *Model {
	model := Model{
		Terms:         []string{},
		TermsDict:     map[string]int{},
		Index:         map[string][][]int{},
		IndexTail:     map[string]map[string]bool{},

		Affects:       [][][]int{},
		KnownAffects:  make([]bool, 15),

		Depth:         2,
		IndexSplitLen: DefaultIndexSplitLen,

		TermsCounts:   make([]float64, 0),
	}

	model.InitMeasurers()
	return &model
}

func (model *Model) InitMeasurers()  {
	procCount := runtime.GOMAXPROCS(-1)
	model.distanceMeasurers = make([]*DistanceMeasurer, procCount)
	for i := 0; i < procCount; i++ {
		model.distanceMeasurers[i] = NewDistanceMeasurer()
	}
	model.currentMeasurerId = 0
	model.measurersSemaphore = make(chan bool, procCount)
	model.measurersMutex = sync.Mutex{}
}

/**
	Has term been added to the model
 */
func (model *Model) HasTerm(term string) bool {
	_, ok := model.TermsDict[strings.ToLower(term)]
	return ok
}

/**
	Parse text, split it to terms and add them to the model
 */
func (model *Model) TrainText(text []byte) {
	var (
		termsRegex = regexp.MustCompile(`(?m)[\p{L}-]+`)
		rawTermsB = termsRegex.FindAll(text, -1)
		terms = map[string]float64{}
	)
	for _, termB := range rawTermsB {
		if utf8.RuneCount(termB) < DefaultMinTermLen {
			continue
		}
		term := string(termB)
		term = strings.ToLower(term)
		terms[term] += 1
	}

	for term, count := range terms {
		if count >= DefaultMinTermCount {
			model.AddTerm(term, count)
		}
	}
}


/**
	Directly add terms to the model
 */
func (model *Model) TrainTerms(terms []string) {
	for _, term := range terms {
		model.AddTerm(term, 1)
	}
}

/**
	Add one term to the model
 */
func (model *Model) AddTerm(term string, count float64) bool {
	var (
		termLen    = utf8.RuneCountInString(term)
		ok         = false
		termId     = 0
		termLo     = strings.ToLower(term)
	)
	if termLen < DefaultMinTermLen {
		return false
	}

	model.TotalTerms += count
	if termId, ok = model.TermsDict[termLo]; ok {
		model.TermsCounts[termId] += count
		return true
	}

	var (
		termsByLen [][]int
		termsIndex []int
		termI      = termLen - 1
	)
	termId = len(model.Terms)
	model.Terms = append(model.Terms, termLo)
	model.TermsCounts = append(model.TermsCounts, count)
	model.TermsDict[termLo] = termId
	edits := GetMultiEdits(termLo, 0.0, float64(model.Depth))

	for edit := range edits {
		editHead, editTail := model.splitEdit(edit)
		if termsByLen, ok = model.Index[editHead]; !ok {
			termsByLen = make([][]int, termLen)
			model.Index[editHead] = termsByLen
		} else {
			if termLen > len(termsByLen) {
				model.Index[editHead] = make([][]int, termLen)
				copy(model.Index[editHead], termsByLen)
				termsByLen = model.Index[editHead]
			}
		}

		termsIndex = termsByLen[termI]
		if termsIndex == nil {
			termsIndex = []int{termId}
			termsByLen[termI] = termsIndex
		} else if termsIndex[len(termsIndex) - 1] != termId {
			termsIndex = append(termsIndex, termId)
		}
		termsByLen[termI] = termsIndex

		// edit tail
		if editTail != "" {
			tailToHeads := model.IndexTail[editTail]
			if tailToHeads == nil {
				tailToHeads = make(map[string]bool)
				model.IndexTail[editTail] = tailToHeads
			}
			tailToHeads[editHead] = true
		}
	}

	// fill known affects
	if termLen > len(model.KnownAffects) {
		knownAffects := make([]bool, termLen)
		copy(knownAffects, model.KnownAffects)
		model.KnownAffects = knownAffects
	}
	if !model.KnownAffects[termI] {
		trackingMultiEdits := GetTrackingMultiEdits(termLo, OperationAffectedChange{0, map[int]bool{}}, float64(model.Depth))
		for edit, trackingEdit := range trackingMultiEdits {
			editLen := len(edit)
			for inputDiff := range trackingEdit.InputLens {
				inputLen := termLen + inputDiff
				if inputLen > len(model.Affects) {
					affects := make([][][]int, inputLen)
					copy(affects, model.Affects)
					model.Affects = affects
				}
				inputAffects := model.Affects[inputLen-1]
				if inputAffects == nil {
					inputAffects = make([][]int, termLen)
					model.Affects[inputLen-1] = inputAffects
				} else if editLen > len(inputAffects) {
					inputAffectsCp := make([][]int, editLen)
					copy(inputAffectsCp, inputAffects)
					inputAffects = inputAffectsCp
					model.Affects[inputLen-1] = inputAffectsCp
				}

				termsLens := model.Affects[inputLen-1][editLen-1]
				if termsLens == nil {
					termsLens = make([]int, 0, 10)
					model.Affects[inputLen-1][editLen-1] = termsLens
				}
				index := sort.SearchInts(termsLens, termI)
				if index == len(termsLens) || termsLens[index] != termI {
					termsLens := append([]int{termI}, termsLens...)
					sort.Ints(termsLens)
					model.Affects[inputLen-1][editLen-1] = termsLens
				}
			}
		}
		model.KnownAffects[termLen-1] = true
	}
	return true
}

/**
	Calculate raw unsorted suggestions
 */
func (model *Model) GetRawSuggestions(input string, calcEditorialPrescription bool) map[string]Suggestion {
	result := make(map[string]Suggestion)
	input = strings.ToLower(input)
	var (
		termsByLen   [][]int
		termsIndex   []int
		ok           = false
		inputLen     = utf8.RuneCountInString(input)
		editLen      int
		inputAffects [][]int
		editAffects  []int
		term         string
	)

	// todo add min input len check

	// exact match
	if termIndex, ok := model.TermsDict[input]; ok {
		result[input] = Suggestion{
			Term:     input,
			Distance: 0,
			Score:    0,
			Count:    model.TermsCounts[termIndex],
		}
	}

	// Index doesn't have any term that can be potentially mathed to input
	if inputLen > len(model.Affects) {
		return result
	}
	inputAffects = model.Affects[inputLen-1]
	if inputAffects == nil {
		return result
	}

	edits := GetMultiEdits(input, 0.0, float64(model.Depth))
	for edit := range edits {
		editHead, editTail := model.splitEdit(edit)

		if termsByLen, ok = model.Index[editHead]; !ok {
			continue
		}

		if editTail != ""  && (model.IndexTail[editTail] == nil || !model.IndexTail[editTail][editHead]) {
			continue
		}

		editLen = utf8.RuneCountInString(edit)
		if editLen > len(inputAffects) {
			continue
		}
		editAffects = inputAffects[editLen-1]
		if editAffects == nil {
			continue
		}

		for _, termI := range editAffects {
			if termI >= len(termsByLen) {
				continue
			}
			termsIndex = termsByLen[termI]
			if termsIndex == nil {
				continue
			}

			for _, termIndex := range termsIndex {
				term = model.Terms[termIndex]

				model.measurersSemaphore <- true
				model.measurersMutex.Lock()
				model.currentMeasurerId++
				if model.currentMeasurerId == len(model.distanceMeasurers) {
					model.currentMeasurerId = 0
				}
				measurer := model.distanceMeasurers[model.currentMeasurerId]
				model.measurersMutex.Unlock()
				distance, editorialPrescription := measurer.Distance(term, input, calcEditorialPrescription)
				<- model.measurersSemaphore

				if distance > model.Depth {
					continue
				}
				if _, ok = result[term]; !ok {
					result[term] = Suggestion{
						Term:         term,
						Distance:     distance,
						Prescription: editorialPrescription,
						Score:        0,
						Count:    model.TermsCounts[termIndex],
					}
				}
			}
		}
	}

	return result
}

/**
	Return suggestions sorted by given scorer
*/
func (model *Model) GetSuggestions(input string, scoreModel ScoreModel, calcEditorialPrescription bool) []Suggestion  {
	var rawSuggestions = model.GetRawSuggestions(input, calcEditorialPrescription)
	suggestions := make([]Suggestion, 0, len(rawSuggestions))
	for _, suggestion := range rawSuggestions {
		suggestion.Score = scoreModel.Score(&suggestion)
		suggestions = append(suggestions, suggestion)
	}
	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].Score < suggestions[j].Score
	})
	return suggestions
}

func (model *Model) splitEdit(edit string) (string, string) {
	var (
		editR     = []rune(edit)
		editRHead = editR

		editHead = edit
		editTail string
	)
	if len(editR) > model.IndexSplitLen {
		editRHead = editR[:model.IndexSplitLen]
		editHead = string(editRHead)
		editTail = string(editR[model.IndexSplitLen:])
	}
	return editHead, editTail
}

func GetMultiEdits(term string, usedWeight float64, maxWeight float64) map[string]float64 {
	edits := GetEdits(term, usedWeight, maxWeight)
	if usedWeight < maxWeight {
		traversalEdits := make(map[string]float64)
		for k, v := range edits {
			traversalEdits[k] = v
		}
		for term, weight := range traversalEdits {
			subedits := GetMultiEdits(term, usedWeight+weight, maxWeight)
			for subterm, v := range subedits {
				edits[subterm] = v
			}
		}
	}
	return edits
}

func GetEdits(term string, usedWeight float64, maxWeight float64) map[string]float64 {
	result := make(map[string]float64)
	termR := []rune(term)
	lenF := float64(len(termR))
	for _, operationWeight := range OperationWeights {
		if lenF-operationWeight.Weight < DefaultMinSpanningLen {
			break
		}
		if usedWeight+operationWeight.Weight > maxWeight {
			break
		}
		edit := string(termR[operationWeight.AffectedLen:])
		if existingWeight, ok := result[edit]; !ok || existingWeight > operationWeight.Weight {
			result[edit] = usedWeight + operationWeight.Weight
		}
		lenR := len(termR) - operationWeight.AffectedLen + 1
		for i := 1; i < lenR; i++ {
			edit := string(string(termR[0:i]) + string(termR[i+operationWeight.AffectedLen:]))
			if existingWeight, ok := result[edit]; !ok || existingWeight > operationWeight.Weight {
				result[edit] = usedWeight + operationWeight.Weight
			}
		}
	}
	return result;
}

func GetTrackingMultiEdits(term string, usedAffectedChange OperationAffectedChange, maxWeight float64) map[string]*OperationAffectedChange {
	edits := GetTrackingEdits(term, usedAffectedChange, maxWeight)
	if usedAffectedChange.Weight < maxWeight {
		traversalEdits := make(map[string]*OperationAffectedChange)
		for k, v := range edits {
			traversalEdits[k] = v
		}
		for term, affectedChange := range traversalEdits {
			u := OperationAffectedChange{
				Weight:    usedAffectedChange.Weight + affectedChange.Weight,
				InputLens: map[int]bool{},
			}
			for l := range affectedChange.InputLens {
				u.InputLens[l] = true
			}
			subedits := GetTrackingEdits(term, u, maxWeight)
			for subterm, v := range subedits {
				edits[subterm] = v
			}
		}
	}
	return edits
}

func GetTrackingEdits(term string, usedAffectedChange OperationAffectedChange, maxWeight float64) map[string]*OperationAffectedChange {
	result := make(map[string]*OperationAffectedChange)
	termR := []rune(term)
	lenF := float64(len(termR))
	for _, operationWeight := range CheckOperationWeight {
		if lenF-operationWeight.Weight < DefaultMinSpanningLen {
			break
		}
		if usedAffectedChange.Weight+operationWeight.Weight > maxWeight {
			break
		}
		edit := string(termR[operationWeight.AffectedLen:])
		if existingWeight, ok := result[edit]; !ok {
			existingWeight = &OperationAffectedChange{
				Weight:    usedAffectedChange.Weight + operationWeight.Weight,
				InputLens: map[int]bool{},
			}
			for _, l1 := range operationWeight.MisspellLens {
				if len(usedAffectedChange.InputLens) > 0 {
					for l2 := range usedAffectedChange.InputLens {
						existingWeight.InputLens[l1+l2] = true
					}
				} else {
					existingWeight.InputLens[l1] = true
				}
			}
			result[edit] = existingWeight
		}
		lenR := len(termR) - operationWeight.AffectedLen + 1
		for i := 1; i < lenR; i++ {
			edit := string(termR[0:i]) + string(termR[i+operationWeight.AffectedLen:])
			if existingWeight, ok := result[edit]; !ok {
				existingWeight = &OperationAffectedChange{
					Weight:    usedAffectedChange.Weight + operationWeight.Weight,
					InputLens: map[int]bool{},
				}
				for _, l1 := range operationWeight.MisspellLens {
					if len(usedAffectedChange.InputLens) > 0 {
						for l2 := range usedAffectedChange.InputLens {
							existingWeight.InputLens[l1+l2] = true
						}
					} else {
						existingWeight.InputLens[l1] = true
					}
				}
				result[edit] = existingWeight
			}
		}
	}
	return result;
}
