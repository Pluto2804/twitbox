package model

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"time"
)

type Twit struct {
	Id      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}
type TwitModel struct {
	DB *sql.DB
}

func (tm *TwitModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `INSERT INTO twits(title,content,created,expires)
	        VALUES(?,?,UTC_TIMESTAMP(),DATE_ADD(UTC_TIMESTAMP(),INTERVAL ? DAY))`
	result, err := tm.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}
	//LastInsertId() method on the result to get the ID of our
	//newly inserted record in twits table
	id, err := result.LastInsertId()
	if err != nil {
		log.Fatalf("Something went wrong!-%s", err)
		os.Exit(-1)
	}
	return int(id), nil
}
func (tm *TwitModel) Get(id int) (*Twit, error) {
	stmt := `SELECT id,title,content,created,expires FROM twits 
	         WHERE expires > UTC_TIMESTAMP() AND id=?`
	row := tm.DB.QueryRow(stmt, id)
	s := &Twit{}
	err := row.Scan(&s.Id, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}
func (tm *TwitModel) Latest() ([]*Twit, error) {
	stmt := `SELECT id,title,content,created,expires FROM twits 
	      WHERE expires>UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`
	rows, err := tm.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	twits := []*Twit{}
	for rows.Next() {
		t := &Twit{}
		err := rows.Scan(&t.Id, &t.Title, &t.Content, &t.Created, &t.Expires)
		if err != nil {
			return nil, err
		}
		twits = append(twits, t)

	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return twits, nil

}
