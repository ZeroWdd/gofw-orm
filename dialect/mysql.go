package dialect

import (
	"reflect"
	"time"
)

type mysql struct {
}

func init() {
	RegisterDialect("mysql", &mysql{})
}

func (m *mysql) DataTypeOf(t reflect.Value) string {
	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uintptr:
		return "integer"
	case reflect.Int64, reflect.Uint64:
		return "bigint"
	case reflect.Float32, reflect.Float64:
		return "double"
	case reflect.String:
		return "varchar"
	case reflect.Struct:
		if _, ok := t.Interface().(time.Time); ok {
			return "datetime"
		}
	}
	panic("implement me")
}

func (m *mysql) TableExistSql(table string) (string, []interface{}) {
	panic("implement me")
}
