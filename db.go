package main

import "database/sql"

type Counter struct {
	Clicks   int
	FullURL  string
	ShortURL string
}

func createCounter(db *sql.DB, x *Counter) error {
	q := `
	insert into shortcodes
		(fullurl, shorturl)
	values
		($1, $2);
	`
	_, err := db.Exec(q, x.FullURL, x.ShortURL)
	return err
}

func getCounter(db *sql.DB, shortCode string) (*Counter, error) {
	q := `
		select fullurl, clicks
		from shortcodes
		where shorturl = $1;
	`
	c := Counter{ShortURL: shortCode}

	err := db.QueryRow(q, shortCode).Scan(&c.FullURL, &c.Clicks)
	return &c, err
}

func getCounters(db *sql.DB) ([]Counter, error) {
	q := `
		select fullurl, shorturl, clicks
		from shortcodes;
	`
	var counters []Counter
	rows, err := db.Query(q)
	if err != nil {
		return counters, err
	}
	defer rows.Close()
	for rows.Next() {
		var c Counter
		if err := rows.Scan(&c.FullURL, &c.ShortURL, &c.Clicks); err != nil {
			return counters, err
		}
		counters = append(counters, c)
	}
	return counters, rows.Err()
}

func incrementClicks(db *sql.DB, c *Counter) error {
	q := `
		update shortcodes
		set clicks = clicks + 1
		where shorturl = $1;
	`
	_, err := db.Exec(q, c.ShortURL)
	return err

}
