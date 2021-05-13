package session

import (
	"reflect"

	"gitee.com/wudongdongfw/gofw-orm/clause"
)

func (s *Session) Insert(values ...interface{}) (int64, error) {
	recordValues := make([]interface{}, 0)

	for _, value := range values {
		table := s.Model(value).RefTable()
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

func (s *Session) Find(values interface{}) error {
	destSlice := reflect.Indirect(reflect.ValueOf(values))
	destType := destSlice.Type().Elem()

	table := s.Model(reflect.New(destType).Elem().Interface()).RefTable()

	//fmt.Println(destSlice)                                // []
	//fmt.Println(destSlice.Type())                         // []session.User
	//fmt.Println(destSlice.Type().Elem())                  // session.User

	//fmt.Println(reflect.New(destType))                    // 创建空对象指针
	//fmt.Println(reflect.New(destType).Elem())             // 获取空对象value封装
	//fmt.Println(reflect.New(destType).Elem().Interface()) // 空对象转interface

	s.clause.Set(clause.SELECT, table.Name, table.FieldNames)
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	rows, err := s.Raw(sql, vars...).QueryRows()
	if err != nil {
		return err
	}

	for rows.Next() {
		dest := reflect.New(destType).Elem()
		var values []interface{}
		for i := 0; i < len(table.Fields); i++ {
			values = append(values, dest.Field(i).Addr().Interface())
		}
		//fmt.Println(values) 反射得到的values切片是创建的dest各个字段的地址 [0xc0000a65c0 0xc0000a65d0]
		if err := rows.Scan(values...); err != nil {
			return err
		}

		destSlice.Set(reflect.Append(destSlice, dest))
	}
	return rows.Close()
}
