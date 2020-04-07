package service

import (
	"fmt"
	"testing"

	"def2sql/helpers"
	"def2sql/regions"

	"github.com/stretchr/testify/assert"
)

func TestRecord_String(t *testing.T) {
	r := record{
		First:        100,
		Last:         300,
		Range:        200,
		Provider:     "A",
		SourceRegion: "B",
		Region:       10,
	}

	expected := "100-300; A, B"

	if !assert.Equal(t, expected, fmt.Sprint(r)) {
		return
	}
}

func TestRecords_ReadFromCSV(t *testing.T) {
	data := []byte(`
		100;20;30;11;A;B
		100;300;200;100;C;D
		1000;300;400;101;E;F;EXTRA FIELD 1; EXTRA FIELD 2
		AAA;2;3;4;5;6
		1;BBB;3;4;5;6
		1;2;CCC;4;5;6
		1;2;3;DDD;5;6`)

	expected := records{
		Correct: []record{
			{
				First:        10020,
				Last:         10030,
				Range:        11,
				Provider:     "A",
				SourceRegion: "B",
				Region:       0,
			},
		},
		UnknownRegion: nil,
		Wrong: map[string][]string{
			"record on line 4: wrong number of fields": {"1000;300;400;101;E;F;EXTRA FIELD 1;EXTRA FIELD 2"},
			"wrong range": {"100;300;200;100;C;D"},
			"strconv.Atoi: parsing \"1BBB\": invalid syntax": {"1;BBB;3;4;5;6"},
			"strconv.Atoi: parsing \"AAA2\": invalid syntax": {"AAA;2;3;4;5;6"},
			"strconv.Atoi: parsing \"1CCC\": invalid syntax": {"1;2;CCC;4;5;6"},
			"strconv.Atoi: parsing \"DDD\": invalid syntax":  {"1;2;3;DDD;5;6"},
		},
		Warnings: nil,
	}

	res := records{
		Wrong: map[string][]string{},
	}

	res.readFromCSV(data)

	if !assert.Equal(t, expected, res) {

	}
}

func TestRecord_checkRange(t *testing.T) {
	r := records{
		Correct: []record{
			{
				First:        10030,
				Last:         10050,
				Range:        21,
				Provider:     "C",
				SourceRegion: "D",
				Region:       20,
			},
			{
				First:        10020,
				Last:         10030,
				Range:        11,
				Provider:     "A",
				SourceRegion: "B",
				Region:       10,
			},
		},
		UnknownRegion: nil,
		Wrong:         nil,
		Warnings:      nil,
	}

	expected := records{
		Correct: []record{
			{
				First:        10020,
				Last:         10030,
				Range:        11,
				Provider:     "A",
				SourceRegion: "B",
				Region:       10,
			},
			{
				First:        10030,
				Last:         10050,
				Range:        21,
				Provider:     "C",
				SourceRegion: "D",
				Region:       20,
			},
		},
		UnknownRegion: nil,
		Wrong:         nil,
		Warnings: []string{
			"last value 10030 of previous record exceeds or equals initial value 10030",
		},
	}

	r.checkRange()

	if !assert.Equal(t, expected, r) {
		return
	}
}

func TestRecords_fillRegions(t *testing.T) {
	regs := regions.Regions{
		1: {"V", []string{"V"}, nil},
		2: {"W", []string{"W"}, []string{"V"}},
		3: {"X", []string{"X", "Y"}, nil},
		4: {"Y", []string{"Y"}, nil},
	}

	recs := records{
		Correct: []record{
			{
				First:        10021,
				Last:         10030,
				Range:        10,
				Provider:     "A",
				SourceRegion: "VW",
				Region:       1,
			},
			{
				First:        10031,
				Last:         10050,
				Range:        20,
				Provider:     "B",
				SourceRegion: "X",
				Region:       3,
			},
			{
				First:        10051,
				Last:         10080,
				Range:        30,
				Provider:     "C",
				SourceRegion: "Y",
				Region:       0,
			},
			{
				First:        10081,
				Last:         10089,
				Range:        9,
				Provider:     "D",
				SourceRegion: "Z",
				Region:       0,
			},
		},
		UnknownRegion: nil,
		Wrong:         nil,
		Warnings:      nil,
	}

	expected := records{
		Correct: []record{
			{
				First:        10021,
				Last:         10030,
				Range:        10,
				Provider:     "A",
				SourceRegion: "VW",
				Region:       1,
			},
			{
				First:        10031,
				Last:         10050,
				Range:        20,
				Provider:     "B",
				SourceRegion: "X",
				Region:       3,
			},
			{
				First:        10051,
				Last:         10080,
				Range:        30,
				Provider:     "C",
				SourceRegion: "Y",
				Region:       0,
			},
			{
				First:        10081,
				Last:         10089,
				Range:        9,
				Provider:     "D",
				SourceRegion: "Z",
				Region:       0,
			},
		},
		UnknownRegion: []record{
			{
				First:        10051,
				Last:         10080,
				Range:        30,
				Provider:     "C",
				SourceRegion: "Y",
				Region:       0,
			},
			{
				First:        10081,
				Last:         10089,
				Range:        9,
				Provider:     "D",
				SourceRegion: "Z",
				Region:       0,
			},
		},
		Wrong: nil,
		Warnings: []string{
			"too many regions found (3,4) for record (10051-10080; C, Y)",
			"couldn't find region for record (10081-10089; D, Z)",
		},
	}

	recs.fillRegions(regs)

	if !assert.Equal(t, expected, recs) {
		return
	}
}

func TestRecord_report(t *testing.T) {
	recs := records{
		Correct: []record{
			{},
		},
		UnknownRegion: []record{
			{
				First:        100,
				Last:         199,
				Range:        99,
				Provider:     "A",
				SourceRegion: "B",
				Region:       1,
			},
		},
		Wrong: map[string][]string{
			"oops": {"1", "2"},
		},
		Warnings: []string{"foo", "bar"},
	}
	expected, err := helpers.Get("testdata/report")
	if !assert.Nil(t, err) {
		return
	}

	if !assert.Equal(t, string(expected), recs.report(1)) {
		return
	}
}
