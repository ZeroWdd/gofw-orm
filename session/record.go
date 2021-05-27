package session

import (
	"errors"
	"log"
	"reflect"

	"gitee.com/wudongdongfw/gofw-orm/clause"
)

func (s *Session) Insert(values ...interface{}) (int64, error) {
	s.CallMethod(BeforeInsert)
	recordValues := make([]interface{}, 0)
	for _, value := range values {
		dest := reflect.Indirect(reflect.ValueOf(value))
		table := s.Model(reflect.New(dest.Type().Elem()).Interface()).RefTable()
		s.clause.Set(clause.INSERT, table.Name, table.FieldNames)
		recordValues = append(recordValues, table.RecordValues(value))
	}

	s.clause.Set(clause.VALUES, recordValues...)
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)
	s.CallMethod(AfterInsert)
	return s.runSQL(sql, vars)
}

func (s *Session) Find(args interface{}) error {
	s.CallMethod(BeforeQuery)
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
	s.CallMethod(AfterQuery)
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

// func (s *Session) Update(values ...interface{}) (int64, error) {
// 	m, ok := values[0].(map[string]interface{})
// 	if !ok {
// 		m = make(map[string]interface{})
// 		length := len(values)
// 		for i := 0; i < length; i += 2 {
// 			m[values[i].(string)] = values[i+1]
// 		}
// 	}

// 	s.clause.Set(clause.UPDATE, s.RefTable().Name, m)
// 	sql, vars := s.clause.Build(clause.UPDATE, clause.WHERE)

// 	return s.runSQL(sql, vars)
// }

func (s *Session) Update(values ...interface{}) (int64, error) {
	s.CallMethod(BeforeUpdate)
	dest := reflect.Indirect(reflect.ValueOf(values[0]))
	log.Printf("dest.type:[%v]\n", dest)
	log.Printf("dest.NumField:[%v]\n", dest.Elem().NumField())
	m := make(map[string]interface{})
	schema := s.refTable
	for i := 0; i < dest.Elem().NumField(); i++ {
		// log.Printf("dest.Field.type:[%v]\n", dest.Field(i).Type())
		// log.Printf("dest.Field.kind:[%v]\n", dest.Field(i).Kind())
		// log.Printf("dest.Field.Name:[%v]\n", dest.Field(i))
		// log.Printf("dest.Field.name:[%v]\n", dest.Field(i).Interface())
		// log.Printf("dest.Field.name:[%v]\n", dest.Type().Field(i).Name)
		val := dest.Elem().Field(i).Interface()
		switch dest.Elem().Field(i).Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
			reflect.Uintptr, reflect.Int64, reflect.Uint64, reflect.Float32, reflect.Float64:
			if val == 0 {
				continue
			}
		case reflect.String:
			if val == "" {
				continue
			}
		case reflect.Struct:

		}
		name := schema.GetField(dest.Elem().Type().Field(i).Name).Name
		m[name] = val
	}
	s.clause.Set(clause.UPDATE, s.RefTable().Name, m)
	sql, vars := s.clause.Build(clause.UPDATE, clause.WHERE)
	s.CallMethod(AfterUpdate)
	return s.runSQL(sql, vars)
}

func (s *Session) Delete() (int64, error) {
	s.CallMethod(BeforeDelete)
	s.clause.Set(clause.DELETE, s.RefTable().Name)
	sql, vars := s.clause.Build(clause.DELETE, clause.WHERE)
	s.CallMethod(AfterDelete)
	return s.runSQL(sql, vars)
}

func (s *Session) Count() (int64, error) {
	s.clause.Set(clause.COUNT, s.RefTable().Name)
	sql, vars := s.clause.Build(clause.COUNT, clause.WHERE)

	row := s.Raw(sql, vars...).QueryRow()
	var tmp int64
	if err := row.Scan(&tmp); err != nil {
		return 0, err
	}
	return tmp, nil
}

func (s *Session) runSQL(sql string, vars []interface{}) (int64, error) {
	exec, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return exec.RowsAffected()
}
