package csv

import "io/ioutil"
import "lib/csv"
var Patterns  []Pattern

type Pattern struct {
	Cards int    `csv:"cards"`  //
	Value int     `csv:"value"` //
	Fan   int    `csv:"fan"`    //
	Name  int    `csv:"name"`   //
}

func InitPattern()error {
	file:= "./csv/test_patterns.csv"
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	Patterns = []Pattern{}
	err = csv.Unmarshal(data, &Patterns)
	if err != nil {
		return err
	}
	return nil
}
