package session

import (
	"log"
)

func (s *Session) Begin() (err error) {
	log.Println("transaction begin")
	if s.tx, err = s.db.Begin(); err != nil {
		log.Printf("Begin error:[%v]", err)
		return
	}
	return
}

func (s *Session) Commit() (err error) {
	log.Println("transaction commit")
	if err = s.tx.Commit(); err != nil {
		log.Printf("Commit error:[%v]", err)
		return
	}
	return
}

func (s *Session) RollBack() (err error) {
	log.Println("transaction rollback")
	if err = s.tx.Rollback(); err != nil {
		log.Printf("RollBack error:[%v]", err)
		return
	}
	return
}
