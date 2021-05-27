package session

import (
	"log"
	"reflect"
)

const (
	BeforeQuery  = "BeforeQuery"
	AfterQuery   = "AfterQuery"
	BeforeUpdate = "BeforeUpdate"
	AfterUpdate  = "AfterUpdate"
	BeforeDelete = "BeforeDelete"
	AfterDelete  = "AfterDelete"
	BeforeInsert = "BeforeInsert"
	AfterInsert  = "AfterInsert"
)

func (s *Session) CallMethod(name string) {
	params := []reflect.Value{reflect.ValueOf(s)}
	fun := reflect.ValueOf(s.RefTable().Model).MethodByName(name)
	if fun.IsValid() {
		if v := fun.Call(params); len(v) > 0 {
			if err, ok := v[0].Interface().(error); ok {
				log.Printf("CallMethod error:[%v]", err)
			}
		}
	}
}
