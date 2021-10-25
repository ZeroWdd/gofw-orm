package schema

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"gitee.com/wudongdongfw/gofw-orm/dialect"
)

type Field struct {
	Name string `json:"name"` // 字段名称
	Type string `json:"type"` // 字段类型
	Tag  string `json:"tag"`  // 字段标签
	Size string `json:"size"` // 字段大小
}

type Schema struct {
	Model      interface{}
	Name       string
	Fields     []*Field
	FieldNames []string
	FieldMap   map[string]*Field
}

func (s *Schema) GetField(name string) *Field {
	return s.FieldMap[name]
}

func Parse(dest interface{}, d dialect.Dialect) (schema *Schema) {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema = &Schema{
		Model:    dest,
		Name:     modelType.Name(),
		FieldMap: make(map[string]*Field),
	}
	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)
		field := newField(dest, d, p)
		schema.Fields = append(schema.Fields, field)
		schema.FieldNames = append(schema.FieldNames, field.Name)
		schema.FieldMap[modelType.Field(i).Name] = field
	}
	return
}

func (s *Schema) RecordValues(dest interface{}) []interface{} {
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	var fieldValues []interface{}
	for i := 0; i < len(s.Fields); i++ {
		fieldValues = append(fieldValues, destValue.Field(i).Interface())
	}
	return fieldValues
}

func newField(dest interface{}, d dialect.Dialect, p reflect.StructField) *Field {
	field := &Field{
		Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
	}
	v, ok := p.Tag.Lookup("fw-orm")
	if !ok {
		errStr, _ := fmt.Printf("model must have a tag:[fw-orm], model:[%v]", dest)
		log.Println(errStr)
		panic(errStr)
	}
	for _, s := range strings.Split(v, ";") {
		i2 := strings.Split(s, ":")
		if len(i2) != 2 {
			errStr, _ := fmt.Printf("model must have a tag:[fw-orm], model:[%v]", dest)
			log.Println(errStr)
			panic(errStr)
		}
		switch i2[0] {
		case "name":
			field.Name = i2[1]
		case "type":
			field.Type = i2[1]
		case "tag":
			field.Tag = i2[1]
		case "size":
			field.Size = i2[1]
		}
	}
	return field
}
