package session

import (
	"database/sql"
	"log"
	"strings"

	"gitee.com/wudongdongfw/gofw-orm/clause"
	"gitee.com/wudongdongfw/gofw-orm/dialect"
	"gitee.com/wudongdongfw/gofw-orm/schema"
)

type Session struct {
	db       *sql.DB
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
	s.clause = clause.Clause{}
}

func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	log.Printf("Exec sql:[%s], values:[%v]", s.sql.String(), s.sqlVars)
	if result, err = s.db.Exec(s.sql.String(), s.sqlVars...); err != nil {
		log.Printf("Exec error:[%v]", err)
	}
	return
}

func (s *Session) Query() *sql.Row {
	defer s.Clear()
	log.Printf("Query sql:[%s], values:[%s]", s.sql.String(), s.sqlVars)
	return s.db.QueryRow(s.sql.String(), s.sqlVars...)
}

func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	log.Printf("QueryRows sql:[%s], values:[%s]", s.sql.String(), s.sqlVars)
	if rows, err = s.db.Query(s.sql.String(), s.sqlVars...); err != nil {
		log.Printf("QueryRows error:[%v]", err)
	}
	return
}
