package spell

import (
	"bufio"
	"io"
	"os"
)

type MisspellParser struct {
}

func InitMisspellParser() *MisspellParser {
	return &MisspellParser{}
}

func (parser *MisspellParser) ParseFromFile(fileName string) ([]Misspell, error) {
	fp, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	return parser.Parse(fp)
}

func (parser *MisspellParser) Parse(reader io.Reader) ([]Misspell, error) {
	scanner := bufio.NewScanner(reader)
	misspells := make([]Misspell, 0, 1000)
	for scanner.Scan() {
		line := scanner.Text()
		if line[0] == '$' {
			term := line[1:]
			misspells = append(misspells, Misspell{
				Term:      term,
				Misspells: make([]string, 4),
			})
		} else {
			misspells[len(misspells)-1].Misspells = append(misspells[len(misspells)-1].Misspells, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return misspells, nil

}
