package service

import (
	"fmt"
	"github.com/gocraft/dbr/v2"

	"def2sql/config"
	"def2sql/exceptions"
	"def2sql/helpers"
	"def2sql/regions"

	"github.com/pkg/errors"

	_ "github.com/go-sql-driver/mysql"
)

type Service struct {
	dataSource []string
	exceptions exceptions.Exceptions
	regions    regions.Regions
	repo       repo
}

func New(c *config.Config) (*Service, error) {
	e, err := exceptions.Load(c.Exceptions)
	if err != nil {
		return nil, err
	}

	r, err := regions.Load(c.Regions)
	if err != nil {
		return nil, err
	}

	db, err := dbr.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s)/%s", c.DB.User, c.DB.Password, c.DB.Host, c.DB.Name),
		nil)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &Service{
		dataSource: c.DataSource,
		exceptions: e,
		regions:    r,
		repo: repo{
			db:    db.NewSession(nil),
			table: c.DB.Table,
		},
	}, nil
}

func (s *Service) Run() error {
	data, err := s.get()
	if err != nil {
		return err
	}
	recs := s.read(data)
	inserted, err := s.save(recs)
	println(recs.report(inserted))

	return err
}

func (s *Service) get() ([][]byte, error) {
	rawData := make([][]byte, 0)
	for _, src := range s.dataSource {
		data, err := helpers.Get(src)
		if err != nil {
			return nil, err
		}
		rawData = append(rawData, data)
	}

	return rawData, nil
}

func (s *Service) read(rawData [][]byte) records {
	rec := makeEmptyRecords()

	for _, data := range rawData {
		data = s.exceptions.Apply(data)
		rec.readFromCSV(data)
	}

	rec.fillRegions(s.regions)
	rec.checkRange()

	return rec
}

func makeEmptyRecords() records {
	return records{
		Wrong: map[string][]string{},
	}
}

func (s *Service) save(r records) (int, error) {
	err := s.repo.create()
	if err != nil {
		return 0, err
	}

	err = s.repo.truncate()
	if err != nil {
		return 0, err
	}

	return s.repo.insert(r.Correct)
}
