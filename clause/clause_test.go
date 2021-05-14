package clause

import (
	"testing"
)

func TestClause_Build(t *testing.T) {
	var clause Clause
	clause.Set(LIMIT, 3, 4)
	clause.Set(SELECT, "User", []string{"*"})
	clause.Set(WHERE, "name = ? and age = ? ", "Tom", 6)
	clause.Set(ORDERBY, "age ASC")
	sql, vars := clause.Build(SELECT, WHERE, ORDERBY, LIMIT)
	t.Log(sql, vars)
}
