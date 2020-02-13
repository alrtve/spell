package main

import (
	"fmt"
	"github.com/alrtve/binary"
	"log"
	"math/rand"
	"path"
	"reflect"
	"spell/scorer/linear"
	"spell/scorer/probabilistic"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	binary.RegisterType(reflect.TypeOf((*linear.Scorer)(nil)).Elem())
	binary.RegisterType(reflect.TypeOf((*probabilistic.Scorer)(nil)).Elem())

	rootDir := "src/spell/cmd/data"
	trainTextFileName := path.Join(rootDir, "eng-uk_web_2002_1M-sentences.txt.gz")
	modelFileName := path.Join(rootDir, "model.bin")
	misspellsFileName := path.Join(rootDir, "misspells.txt.gz")
	learningDataFileName := path.Join(rootDir, "learning.bin")
	learningModelFileName := path.Join(rootDir, "learning.model.bin")


	model, err := GetModelFromCache(modelFileName, trainTextFileName)
	if err != nil {
		log.Fatal(err)
	}

	learningData, err := GetLearningDataFromCache(model, learningDataFileName, misspellsFileName)
	if err != nil {
		log.Fatal(err)
	}

	scorer, err := GetScorerFromCache(learningModelFileName, learningData)
	if err != nil {
		log.Fatal(err)
	}

	totalCount := 0
	validCount := 0
	for _, learningTerm := range learningData {
		suggestions := model.GetSuggestions(learningTerm.Misspell, scorer, true)
		if model.HasTerm(learningTerm.Term) && len(suggestions) > 0 {
			if suggestions[0].Term != learningTerm.Term && suggestions[0].Prescription != nil {
				fmt.Println(learningTerm.Misspell)
				fmt.Println(learningTerm.Term)
				for i := 0 ; i < 3 && i < len(suggestions); i++ {
					if suggestions[i].Prescription != nil {
						suggestions[i].Prescription.Dump()
					} else {
						fmt.Println(suggestions[i].Term)
					}
				}
				fmt.Printf("\n\n")
				totalCount++
			} else {
				validCount++
			}
		}
	}
	fmt.Println(totalCount)
	fmt.Println(validCount)
}