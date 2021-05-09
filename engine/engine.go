package engine

import (
	"database/sql"
	"log"

	"gitee.com/wudongdongfw/gofw-orm/dialect"
	"gitee.com/wudongdongfw/gofw-orm/session"
)

type Engine struct {
	db      *sql.DB
	dbName  string
	dialect dialect.Dialect
}

func NewEngine(driver string, source string) (e *Engine, err error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Printf("NewEngine sql.Open is err:[%v]", err)
		return
	}
	if err = db.Ping(); err != nil {
		log.Printf("NewEngine db.Ping is err:[%v]", err)
		return
	}
	dia, ok := dialect.GetDialect(driver)
	if !ok {
		log.Printf("NewEngine dia.GetDialect is not found :[%s]", driver)
		return
	}
	e = &Engine{
		db:      db,
		dbName:  driver,
		dialect: dia,
	}
	log.Printf("NewEngine connect success, driver is [%s]", driver)
	return
}

func (e *Engine) Close() {
	if err := e.db.Close(); err != nil {
		log.Printf("Failed to close the Database")
		return
	}
	log.Printf("Closed Database success")
}

func (e *Engine) NewSession() *session.Session {
	return session.New(e.db, e.dialect)
}
