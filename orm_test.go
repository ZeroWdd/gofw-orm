package gofw_orm

import (
	"testing"

	"gitee.com/wudongdongfw/gofw-orm/engine"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Id   int64  `fw-orm:"name:id;size:10;tag:PRIMARY KEY AUTO_INCREMENT"`
	Name string `fw-orm:"name:name;size:20"`
}

func TestEngine(t *testing.T) {
	e, _ := engine.NewEngine("mysql", "root:root@tcp(192.168.229.136:3306)/orm")
	session := e.NewSession()
	model := session.Model(&User{})
	_ = model.DropTable()
	_ = model.CreateTable()
}
