package schema

import (
	"fmt"
	"testing"

	"gitee.com/wudongdongfw/gofw-orm/dialect"
)

// schema_test.go
type User struct {
	Id   int64  `fw-orm:"name:id;size:10;tag:PRIMARY KEY AUTO_INCREMENT"`
	Name string `fw-orm:"name:name;size:20"`
	Age  int    `fw-orm:"name:age;size:10"`
}

var TestDial, _ = dialect.GetDialect("mysql")

func TestParse(t *testing.T) {

	schema := Parse(&User{}, TestDial)
	if schema.Name != "User" || len(schema.Fields) != 3 {
		t.Fatal("failed to parse User struct")
	}
}

func TestSchema_RecordValues(t *testing.T) {
	u := &User{
		Id:   1,
		Name: "sss",
	}
	schema := Parse(&User{}, TestDial)
	values := schema.RecordValues(&u)

	fmt.Println(values)
}
