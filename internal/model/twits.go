package model

import (
	"database/sql"
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
	return 0, nil
}
func (tm *TwitModel) Get(id int) (*Twit, error) {
	return nil, nil
}
func (tm *TwitModel) Latest() ([]*Twit, error) {
	return nil, nil
}
