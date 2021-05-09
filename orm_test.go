package gofw_orm

import (
	"testing"

	"gitee.com/wudongdongfw/gofw-orm/engine"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Id int64 `fworm:"primary key"`
	// TODO: When a string is converted to a database type, the type length needs to be set
	// name string
}

func TestEngine(t *testing.T) {
	e, _ := engine.NewEngine("mysql", "root:root@tcp(192.168.229.136:3306)/orm")
	session := e.NewSession()
	model := session.Model(&User{})
	_ = model.DropTable()
	_ = model.CreateTable()

}
