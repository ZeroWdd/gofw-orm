package dialect

import (
	"reflect"
)

var dialectMap = map[string]Dialect{}

type Dialect interface {
	DataTypeOf(t reflect.Value) string
	TableExistSql(table string) (string, []interface{})
}

func RegisterDialect(name string, dialect Dialect) {
	dialectMap[name] = dialect
}

func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectMap[name]
	return
}
