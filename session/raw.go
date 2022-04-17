package session

import (
	"database/sql"
	"log"
	"strings"

	"github.com/ZeroWdd/gofw-orm/clause"
	"github.com/ZeroWdd/gofw-orm/dialect"
	"github.com/ZeroWdd/gofw-orm/schema"
)

type Session struct {
	db       *sql.DB
	tx       *sql.Tx
	dialect  dialect.Dialect
	refTable *schema.Schema

	// TODO: sql and sqlVars can extract struct
	sql     strings.Builder
	sqlVars []interface{}

	// sql构造器
	clause clause.Clause
}

func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{
		db:      db,
		dialect: dialect,
	}
}

type CommonDB interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

func (s *Session) DB() CommonDB {
	if s.tx != nil {
		return s.tx
	}
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
	s.clause = clause.Clause{}
}

func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	log.Printf("Exec sql:[%s], values:[%v]", s.sql.String(), s.sqlVars)
	if result, err = s.DB().Exec(s.sql.String(), s.sqlVars...); err != nil {
		log.Printf("Exec error:[%v]", err)
	}
	return
}

func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	log.Printf("QueryRow sql:[%s], values:[%s]", s.sql.String(), s.sqlVars)
	return s.DB().QueryRow(s.sql.String(), s.sqlVars...)
}

func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	log.Printf("QueryRows sql:[%s], values:[%s]", s.sql.String(), s.sqlVars)
	if rows, err = s.DB().Query(s.sql.String(), s.sqlVars...); err != nil {
		log.Printf("QueryRows error:[%v]", err)
	}
	return
}
