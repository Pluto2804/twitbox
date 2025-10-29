package model

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
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

type User struct {
	Id              int
	Name            string
	Email           string
	Hashed_Password []byte
	Created         time.Time
}
type UserModel struct {
	DB *sql.DB
}

func (um *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	stmt := `INSERT INTO users(name,email,hashed_password,created)
	       VALUES(?,?,?,UTC_TIMESTAMP()) `
	_, err = um.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		var mySqlErr *mysql.MySQLError
		if errors.As(err, &mySqlErr) {
			if mySqlErr.Number == 1062 && strings.Contains(mySqlErr.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err

	}
	return nil
}
func (um *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte
	stmt := `SELECT id,hashed_password FROM users WHERE email = ? `
	err := um.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	return id, nil
}
func (um *UserModel) Exists(id int) (bool, error) {
	var exists bool

	stmt := `SELECT EXISTS(SELECT true from users WHERE id=?)`
	err := um.DB.QueryRow(stmt, id).Scan(&exists)
	return exists, err

}
