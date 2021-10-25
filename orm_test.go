package gofw_orm

import (
	"errors"
	"fmt"
	"log"
	"testing"

	"gitee.com/wudongdongfw/gofw-orm/engine"
	"gitee.com/wudongdongfw/gofw-orm/session"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Id       int64  `fw-orm:"name:id;size:10;tag:PRIMARY KEY AUTO_INCREMENT"`
	Name     string `fw-orm:"name:name;size:20"`
	Age      int    `fw-orm:"name:age;size:10"`
	NickName string `fw-orm:"name:nick_name;size:20"`
}

var model *session.Session

func TestCreateTable(t *testing.T) {
	e, _ := engine.NewEngine("mysql", "root:root@tcp(192.168.146.128:3306)/orm")
	model = e.NewSession().Model(&User{})
	_ = model.CreateTable()
}

var (
	user1 = &User{Name: "Tom", Age: 16, NickName: "Tom_Tom"}
	user2 = &User{Name: "Jerry", Age: 18}
	user3 = &User{Name: "Linda", Age: 16, NickName: "Linda_Linda"}
)

func (u *User) BeforeInsert(s *session.Session) error {
	log.Println("BeforeInsert")
	return nil
}

func (u *User) AfterInsert(s *session.Session) error {
	log.Println("AfterInsert")
	return nil
}

func TestInsert(t *testing.T) {
	TestCreateTable(t)
	_, _ = model.Insert(user1, user2, user3)
}

func (u *User) AfterQuery(s *session.Session) error {
	log.Println("AfterQuery")
	return nil
}

func TestFind(t *testing.T) {
	TestCreateTable(t)
	var users []*User
	if err := model.Find(&users); err != nil {
		log.Fatal("failed to query all")
	}
	fmt.Println(users[0])
}

func TestFirst(t *testing.T) {
	TestCreateTable(t)
	user := &User{}
	if err := model.Where("name = ?", "Tom").First(user); err != nil {
		log.Fatal("failed to query all")
	}
	fmt.Println(user)
}

func TestUpdate(t *testing.T) {
	TestCreateTable(t)
	user := &User{NickName: "test", Age: 66}
	_, _ = model.Where("name = ?", "Tom").Update(user)
}

func TestDelete(t *testing.T) {
	TestCreateTable(t)
	_, _ = model.Where("name = ?", "Linda").Delete()
}

func TestCount(t *testing.T) {
	TestCreateTable(t)
	count, err := model.Where("name = ?", "Tom").Count()
	if err != nil {
		log.Fatal("failed to Count")
	}
	fmt.Println(count)
}

func TestTransaction(t *testing.T) {
	e, _ := engine.NewEngine("mysql", "root:root@tcp(192.168.146.128:3306)/orm")
	e.Transaction(func(s *session.Session) (interface{}, error) {
		s.Model(&User{})
		user := &User{NickName: "test", Age: 77}
		_, _ = s.Where("name = ?", "Tom").Update(user)
		return nil, errors.New("error")
	})
}
