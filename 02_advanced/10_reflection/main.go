/*
是什么：反射通过 reflect 包，在程序运行时查看值的类型、字段、标签和可修改性。
作用：在事先不知道具体类型时读取结构信息；修改字段还需满足可设置、导出可见性和类型匹配等规则。
使用场景：通用序列化、校验和框架内部处理；encoding/json 等工具已经封装许多需求，普通业务先使用具体类型或接口。

目标：运行时读取类型、值、结构体标签，理解可设置性。
对照：不像 JS 可以随意增删属性，Go 反射仍遵循静态类型与导出规则。
运行：go run ./02_advanced/10_reflection
观察：两个字段的名字、类型、json 标签和值；指针解引用后的 Name 可设置并改为 Grace。
练习：尝试对非指针的 reflect.Value 修改字段，先用 CanSet 判断；不要直接调用导致 panic。
*/
package main

import (
	"fmt"
	"reflect"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	user := User{Name: "Ada", Age: 28}
	value := reflect.ValueOf(user) // 提取user的反射值；注意是值拷贝，非指针，无法修改原始结构体。实际就是拿到Ada, 28
	// fmt.Println("reflect ValueOf:", reflect.ValueOf(&user))
	typeInfo := value.Type()
	for index := 0; index < typeInfo.NumField(); index++ {
		field := typeInfo.Field(index)
		fmt.Printf("%s type=%s json=%q value=%v\n", field.Name, field.Type, field.Tag.Get("json"), value.Field(index).Interface())
	}
	fmt.Println("value can set:", value.FieldByName("Name").CanSet())
	field := reflect.ValueOf(&user).Elem().FieldByName("Name")
	if field.IsValid() && field.CanSet() && field.Kind() == reflect.String {
		field.SetString("Grace")
	}
	fmt.Println("changed:", user.Name)
	// 这里用于观察机制；实际修改请直接 user.Name = "Grace"，JSON 编解码请用 encoding/json。
}
