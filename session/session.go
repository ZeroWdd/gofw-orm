package session

import (
	"database/sql"
	"log"
	"strings"
)

type Session struct {
	db      *sql.DB
	sql     strings.Builder
	sqlVars []interface{}
}

func New(db *sql.DB) *Session {
	return &Session{db: db}
}

func (s *Session) DB() *sql.DB {
	return s.db
}

func (s *Session) Raw(sql string, vars ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, vars...)
	return s
}

func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
}

func (s Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	log.Printf("Exec sql:[%s], values:[%s]", s.sql.String(), s.sqlVars)
	if result, err = s.db.Exec(s.sql.String(), s.sqlVars...); err != nil {
		log.Printf("Exec error:[%v]", err)
	}
	return
}

func (s Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	log.Printf("QueryRows sql:[%s], values:[%s]", s.sql.String(), s.sqlVars)
	if rows, err = s.db.Query(s.sql.String(), s.sqlVars...); err != nil {
		log.Printf("QueryRows error:[%v]", err)
	}
	return
}
