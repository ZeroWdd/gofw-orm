package session

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"gitee.com/wudongdongfw/gofw-orm/schema"
)

func (s *Session) Model(value interface{}) *Session {
	if s.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(s.refTable.Model) {
		s.refTable = schema.Parse(value, s.dialect)
	}
	return s
}

func (s *Session) RefTable() *schema.Schema {
	if s.refTable == nil {
		log.Println("RefTable Model is not set")
	}
	return s.refTable
}

func (s *Session) CreateTable() error {
	table := s.RefTable()
	var columns []string
	for _, field := range table.Fields {
		columns = append(columns, fmt.Sprintf("%s %s %s",
			field.Name, field.Type, field.Tag))
	}
	join := strings.Join(columns, ", ")
	str := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s);", table.Name, join)
	_, err := s.Raw(str).Exec()
	return err
}

func (s *Session) DropTable() error {
	_, err := s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s;", s.RefTable().Name)).Exec()
	return err
}

//func (s *Session) HasTable() bool {
//	sql, values := s.dialect.TableExistSql(s.RefTable().Name)
//	row := s.Raw(sql, values...).Query()
//	var tmp string
//	_ = row.Scan(&tmp)
//	return tmp == s.RefTable().Name
//}
