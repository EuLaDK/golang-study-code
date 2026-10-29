package main

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"time"
)

// 反射三大原则
// 1、反射直接的入口：经由接口（interface{}）类型变量值进入反射的世界并获得对应的反射对象
// 2、反射世界的出口：反射对象通过化身为一个接口（interface{}）类型变量值的形式走出反射世界
// 3、修改反射对象的前提：反射对象对应的reflect.Value必须是可设置的（Settable）

func ConstructQueryStmt(obj interface{}) (stmt string, err error) {
	// 仅支持struct或struct指针类型
	typ := reflect.TypeOf(obj)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	if typ.Kind() != reflect.Struct {
		err = errors.New("only struct is supported")
		return
	}

	buffer := bytes.NewBufferString("")
	buffer.WriteString("SELECT ")

	if typ.NumField() == 0 {
		err = fmt.Errorf("the type[%s] has no fields", typ.Name())
		return
	}

	for i := 0; i < typ.NumField(); i++ {
		filed := typ.Field(i)

		if i != 0 {
			buffer.WriteString(", ")
		}
		column := filed.Name
		if tag := filed.Tag.Get("orm"); tag != "" {
			column = tag
		}
		buffer.WriteString(column)
	}
	stmt = fmt.Sprintf("%s FROM %s", buffer.String(), typ.Name())
	return
}

type Product struct {
	ID        uint32
	Name      string
	Price     uint32
	LeftCount uint32 `orm:"left_count"`
	Batch     string `orm:"batch_number"`
	Updated   time.Time
}

type Person struct {
	ID      string
	Name    string
	Age     uint32
	Gender  string
	Addr    string `orm:"address"`
	Updated time.Time
}

func main() {
	stmt, err := ConstructQueryStmt(&Product{})
	if err != nil {
		fmt.Println("construct query stmt for Product error:", err)
		return
	}
	fmt.Println(stmt)

	stmt, err = ConstructQueryStmt(Person{})
	if err != nil {
		fmt.Println("construct query stmt for person error:", err)
		return
	}
	fmt.Println(stmt)
}
