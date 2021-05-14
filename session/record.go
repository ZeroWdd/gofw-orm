package session

import (
	"errors"
	"reflect"

	"gitee.com/wudongdongfw/gofw-orm/clause"
)

func (s *Session) Insert(values ...interface{}) (int64, error) {
	recordValues := make([]interface{}, 0)

	for _, value := range values {
		dest := reflect.Indirect(reflect.ValueOf(value))
		table := s.Model(reflect.New(dest.Type().Elem()).Interface()).RefTable()
		s.clause.Set(clause.INSERT, table.Name, table.FieldNames)
		recordValues = append(recordValues, table.RecordValues(value))
	}

	s.clause.Set(clause.VALUES, recordValues...)
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)

	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (s *Session) Find(args interface{}) error {
	destSlice := reflect.Indirect(reflect.ValueOf(args))
	destType := destSlice.Type().Elem()
	//fmt.Println(destType)
	//fmt.Println(destType.Elem())
	//fmt.Println(reflect.New(destType).Interface())
	//fmt.Println(reflect.New(destType).Elem().Interface())
	//fmt.Println(reflect.New(destType.Elem()).Interface())
	table := s.Model(reflect.New(destType.Elem()).Interface()).RefTable()

	//fmt.Println(destSlice)                      					// []
	//fmt.Println(destSlice.Type())               					// []*User
	//fmt.Println(destSlice.Type().Elem())        					// *User
	//fmt.Println(destSlice.Type().Elem().Elem()) 					// User
	//
	//fmt.Println(reflect.New(destType.Elem()))                    	// &{0  0 }
	//fmt.Println(reflect.New(destType.Elem()).Elem())             	// {0  0 }
	//fmt.Println(reflect.New(destType.Elem()).Elem().Interface()) 	// {0  0 }

	s.clause.Set(clause.SELECT, table.Name, table.FieldNames)
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	rows, err := s.Raw(sql, vars...).QueryRows()
	if err != nil {
		return err
	}

	for rows.Next() {
		//dest := reflect.New(destType.Elem()).Elem()
		dest := reflect.New(destType.Elem())
		var values []interface{}
		for i := 0; i < len(table.Fields); i++ {
			values = append(values, dest.Elem().Field(i).Addr().Interface())
		}
		//fmt.Println(values) 反射得到的values切片是创建的dest各个字段的地址 [0xc0000a65c0 0xc0000a65d0]
		if err = rows.Scan(values...); err != nil {
			return err
		}
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	return rows.Close()
}

func (s *Session) First(args interface{}) error {
	dest := reflect.Indirect(reflect.ValueOf(args))
	destSlice := reflect.New(reflect.SliceOf(dest.Type())).Elem()

	if err := s.Limit(0, 1).Find(destSlice.Addr().Interface()); err != nil {
		return err
	}

	if destSlice.Len() == 0 {
		return errors.New("First not found ")
	}
	dest.Set(destSlice.Index(0))
	return nil
}

func (s *Session) Where(query interface{}, args ...interface{}) *Session {
	var vars []interface{}
	s.clause.Set(clause.WHERE, append(append(vars, query), args...)...)
	return s
}

func (s *Session) Limit(args ...interface{}) *Session {
	s.clause.Set(clause.LIMIT, args...)
	return s
}

func (s *Session) OrderBy(args ...interface{}) *Session {
	s.clause.Set(clause.ORDERBY, args)
	return s
}
