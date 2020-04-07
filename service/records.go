package service

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"def2sql/regions"

	"github.com/pkg/errors"
)

type records struct {
	Correct       []record            `json:"-"`
	UnknownRegion []record            `json:"unknown_regions,omitempty"`
	Wrong         map[string][]string `json:"wrong_records,omitempty"`
	Warnings      []string            `json:"warnings,omitempty"`
}

type record struct {
	First        int    `db:"first"`
	Last         int    `db:"last"`
	Range        int    `db:"-"`
	Provider     string `db:"provider"`
	SourceRegion string `db:"source_region"`
	Region       int    `db:"region"`
}

const fieldCount = 6

func (r record) String() string {
	return fmt.Sprintf("%d-%d; %s, %s",
		r.First, r.Last, r.Provider, r.SourceRegion)
}

func (r *records) readFromCSV(data []byte) {
	reader := csv.NewReader(bytes.NewReader(data))
	reader.LazyQuotes = true
	reader.Comma = ';'
	reader.TrimLeadingSpace = true
	reader.FieldsPerRecord = fieldCount

	for {
		rec, err := reader.Read()
		if err == io.EOF {
			break
		}

		r.addRecord(rec, err)
	}
}

func (r *records) addRecord(record []string, err error) {
	if err != nil {
		r.addWrongRecord(record, err)
		return
	}

	rec, err := parseRecord(record)
	if err != nil {
		r.addWrongRecord(record, err)
		return
	}

	r.Correct = append(r.Correct, rec)
}

func (r *records) addWrongRecord(rec []string, err error) {
	msg := err.Error()
	r.Wrong[msg] = append(r.Wrong[msg], strings.Join(rec, ";"))
}

func parseRecord(data []string) (record, error) {
	first, err := strconv.Atoi(fmt.Sprintf("%s%s", data[0], data[1]))
	if err != nil {
		return record{}, err
	}
	last, err := strconv.Atoi(fmt.Sprintf("%s%s", data[0], data[2]))
	if err != nil {
		return record{}, err
	}
	rng, err := strconv.Atoi(fmt.Sprintf("%s", data[3]))
	if err != nil {
		return record{}, err
	}

	if last-first != rng-1 {
		return record{}, errors.New("wrong range")
	}

	return record{
		First:        first,
		Last:         last,
		Range:        rng,
		Provider:     data[4],
		SourceRegion: data[5],
	}, nil
}

func (r *records) checkRange() {
	r.sort()

	last := 0
	for _, rec := range r.Correct {
		if rec.First <= last {
			r.addWarning(fmt.Sprintf("last value %d of previous record exceeds or equals initial value %d", last, rec.First))
			continue
		}

		last = rec.Last
	}
}

func (r *records) sort() {
	sort.Slice(r.Correct, func(i, j int) bool {
		return r.Correct[i].First < r.Correct[j].First
	})
}

func (r *records) addWarning(warn string) {
	r.Warnings = append(r.Warnings, warn)
}

type codes map[int]struct{}

func (r *records) fillRegions(reg regions.Regions) {
	for i, rec := range r.Correct {
		c := codes(reg.Match(rec.SourceRegion))

		switch len(c) {
		case 1:
			r.Correct[i].Region = c.firstOrDefault()
			continue
		case 0:
			r.addWarning(fmt.Sprintf("couldn't find region for record (%s)", rec))
		default:
			r.addWarning(fmt.Sprintf("too many regions found (%s) for record (%s)", c, rec))
		}

		r.addRecordWithoutRegion(rec)
	}
}

func (r *records) addRecordWithoutRegion(rec record) {
	r.UnknownRegion = append(r.UnknownRegion, rec)
}

func (c codes) firstOrDefault() int {
	for code := range c {
		return code
	}

	return 0
}

func (c codes) String() string {
	str := make([]string, 0, len(c))
	for _, code := range c.toSortedSlice() {
		str = append(str, strconv.Itoa(code))
	}

	return strings.Join(str, ",")
}

func (c codes) toSortedSlice() []int {
	res := make([]int, 0, len(c))
	for code := range c {
		res = append(res, code)
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i] < res[j]
	})

	return res
}

func (r *records) report(inserted int) string {
	report := fmt.Sprintf("correct records amount: %d", len(r.Correct))
	report = fmt.Sprintf("%s\ninserted %d records", report, inserted)

	b, err := json.MarshalIndent(r, "", "    ")
	if err != nil {
		return fmt.Sprintf("%s\n%+v", report, err)
	}

	if string(b) == "{}" {
		return report
	}

	return fmt.Sprintf("%s\n%s", report, b)
}
