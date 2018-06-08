package condors

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", "root:@cloudsql(condors-fanclub:us-central1:condors)/")
	if err != nil {
		log.Fatalf("Could not open db: %v", err)
	}
}

// An Observation is a data point: This User has witnessed this many condors
// at this date, in the geographic area.
type Observation struct {
	ID        int
	Date      time.Time
	Username  string
	Region    string
	NbCondors int
}

func queryObservations(c context.Context, year int) ([]Observation, error) {
	rows, err := db.QueryContext(c, `
		SELECT id
		FROM condors.observation
		WHERE YEAR(date)=?
	`, year)
	if err != nil {
		return nil, err
	}
	var ids []int
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	observations := make([]Observation, len(ids))
	for i, id := range ids {
		observations[i], err = queryObservation(c, id)
		if err != nil {
			return nil, err
		}

	}
	return observations, nil
}

func queryObservation(c context.Context, id int) (Observation, error) {
	var obs Observation
	rows, err := db.QueryContext(c, `
		SELECT id, date, region, user, nbcondors
		FROM condors.observation
		WHERE id=?
	`, id)
	if err != nil {
		return obs, err
	}
	if !rows.Next() {
		return obs, fmt.Errorf("No observation found with id %d", id)
	}
	var date mysql.NullTime
	err = rows.Scan(
		&obs.ID,
		&date,
		&obs.Region,
		&obs.Username,
		&obs.NbCondors,
	)
	if err != nil {
		return obs, err
	}
	obs.Date = date.Time
	if rows.Next() {
		return obs, fmt.Errorf("Multiple observations found with id %d", id)
	}
	return obs, nil
}
