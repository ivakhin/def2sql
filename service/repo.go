package service

import (
	"fmt"

	"github.com/gocraft/dbr/v2"
	"github.com/pkg/errors"
)

type repo struct {
	db    *dbr.Session
	table string
}

func (r *repo) truncate() error {
	_, err := r.db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", r.table))
	return errors.WithStack(err)
}

func (r *repo) create() error {
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			first BIGINT PRIMARY KEY NOT NULL,
			last BIGINT UNIQUE NOT NULL,
			provider TEXT,
			source_region TEXT,
			region INT NOT NULL)`, r.table)

	_, err := r.db.Exec(query)
	return errors.WithStack(err)
}

const chunkSize = 1000

func (r *repo) insert(recs []record) (int, error) {
	count := 0
	stmt := r.db.InsertInto(r.table).Columns("first", "last", "provider", "source_region", "region")

	for i := 0; i < len(recs); i++ {
		stmt.Record(recs[i])
		if i%chunkSize == 0 || i == len(recs)-1 {
			_, err := stmt.Exec()
			if err != nil {
				return count, errors.WithStack(err)
			}
			stmt = r.db.InsertInto(r.table).Columns("first", "last", "provider", "source_region", "region")
		}
		count++
	}

	return count, nil
}
