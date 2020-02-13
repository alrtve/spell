package main

import (
	"compress/gzip"
	"github.com/alrtve/binary"
	"io/ioutil"
	"os"
	"spell"
	"spell/scorer"
	"spell/scorer/linear"
	"strings"
)


type CacheFile struct {
	*os.File
	isCreated bool
	isInTransaction bool
	fileName  string
}

func OpenCacheFile(fileName string) (cf CacheFile, err error) {
	cf.fileName = fileName
	cf.File, err = os.Open(cf.fileName)
	if os.IsNotExist(err) {
		cf.File, err = os.OpenFile(cf.fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err == nil {
			cf.isCreated = true
		}
	}
	return
}

func (cf CacheFile) Close() error {
	err := cf.File.Close()
	if err == nil && cf.isInTransaction && cf.isCreated {
		err = os.Remove(cf.fileName)
	}
	return err
}

func (cf CacheFile) Commit()()  {
	cf.isInTransaction = true
}

func (cf CacheFile) Exist() bool {
	return !cf.isCreated
}


func getLearningTerm(term, misspelledTerm string, model *spell.Model) (learningTerm *spell.LearningTerm) {
	term = strings.ToLower(term)
	if !model.HasTerm(term) {
		return
	}
	misspelledTerm = strings.ToLower(misspelledTerm)
	rawSuggestions := model.GetRawSuggestions(misspelledTerm, true)
	suggestions := make([]spell.Suggestion, 0)
	for _, suggestion := range rawSuggestions {
		if misspelledTerm != suggestion.Term && suggestion.Prescription != nil {
			suggestions = append(suggestions, suggestion)
		}
	}
	if len(suggestions) > 0 {
		learningTerm =  &spell.LearningTerm{
			Term: term,
			Misspell: misspelledTerm,
			Suggestions: suggestions,
		}
	}
	return
}

func GetModelFromCache(modelFileName, trainTextFileName string) (model *spell.Model, err error) {
	modelCacheFile, err := OpenCacheFile(modelFileName)
	if err != nil {
		return
	}

	defer  modelCacheFile.Close()
	if modelCacheFile.Exist() {
		r := spell.Model{}
		r.InitMeasurers()
		err := binary.UnmarshalFrom(modelCacheFile, &r)
		return &r, err

	}
	fp, err := os.Open(trainTextFileName)
	if err != nil {
		return
	}
	defer fp.Close()

	gzipP, err := gzip.NewReader(fp)
	if err != nil {
		return
	}
	defer gzipP.Close()

	data, _ := ioutil.ReadAll(gzipP)

	//model.TrainTerms(terms.Words)
	model = spell.InitModel()
	model.TrainText(data)
	err = binary.MarshalTo(model, modelCacheFile)
	if err == nil {
		modelCacheFile.Commit()
	}
	return
}

func GetLearningDataFromCache(model *spell.Model, learningDataFileName, misspellsFileName string) (learningData []*spell.LearningTerm, err error) {
	learningDataFile, err := OpenCacheFile(learningDataFileName)
	if err != nil {
		return
	}
	defer learningDataFile.Close()
	if learningDataFile.Exist() {
		learningData = []*spell.LearningTerm{}
		binary.UnmarshalFrom(learningDataFile, &learningData)
		return
	}

	learningData = []*spell.LearningTerm{}
	misspellParser := spell.NewMisspellParser()
	fp, err := os.Open(misspellsFileName)
	if err != nil {
		return
	}
	defer fp.Close()

	gzipP, err := gzip.NewReader(fp)
	if err != nil {
		return
	}
	defer gzipP.Close()

	misspells, _ := misspellParser.Parse(gzipP)
	measurer := spell.NewDistanceMeasurer()
	for i, misspell := range misspells {
		ms := make([]string, 0, len(misspell.Misspells))
		for _, m := range misspell.Misspells {
			if d, _ := measurer.Distance(m, misspell.Term, false); d <= 2 {
				ms = append(ms, m)
			}
		}
		misspells[i].Misspells = ms
		term := strings.ToLower(misspell.Term)
		if !model.HasTerm(term) {
			model.AddTerm(term, 1)
		}
	}

	for _, misspell := range misspells {
		for _, misspelledTerm := range misspell.Misspells {
			learningTerm := getLearningTerm(misspell.Term, misspelledTerm, model)
			if learningTerm != nil {
				learningData = append(learningData, learningTerm)
			}
		}
	}
	err = binary.MarshalTo(learningData, learningDataFile)
	if err == nil {
		learningDataFile.Commit()
	}
	return
}

func GetScorerFromCache(learningModelFileName string, learningData []*spell.LearningTerm) (learninigModel spell.ScoreModel, err error) {
	learningDataFile, err := OpenCacheFile(learningModelFileName)
	if err != nil {
		return
	}
	if learningDataFile.Exist() {
		learninigModel = &linear.Scorer{}
		err = binary.UnmarshalFrom(learningDataFile, learninigModel)
		return
	}

	learnAlgorithm := &linear.Learner{Vectoriser: scorer.InitVectoriser()}
	learninigModel = learnAlgorithm.Learn(learningData)
	err = binary.MarshalTo(learninigModel, learningDataFile)
	if err == nil {
		learningDataFile.Commit()
	}
	return

}
