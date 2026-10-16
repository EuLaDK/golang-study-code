/*
是什么：struct 把多个有名字的字段组成一个类型；方法是带接收者参数、与某个类型关联的函数。
作用：用值接收者处理对象副本，用指针接收者修改原对象；结构体字段可以表达业务数据的固定形状。
使用场景：定义 User、Order 或配置类型，并把相关操作写成方法；这是 Go 组织业务数据与行为的常用方式。

目标：结构体、字段、方法、指针接收者与值接收者。
对照：结构体不是 JS class；方法与类型分开声明，Go 没有构造函数关键字。
运行：go run ./02_advanced/01_structs_methods
观察：value receiver 不改变 Ada，pointer receiver 将名称改为 Grace，复制后的结构体独立。
练习：增加邮箱字段和修改邮箱的方法；说明为何修改方法使用指针接收者。
*/
package main

import "fmt"

type User struct {
	Name  string // 大写名称可被其他包访问；小写名称仅包内可见。
	age   int
	email string
}

func (user User) Greeting() string {
	return "Hello, " + user.Name
}

func (user User) RenameCopy(name string) {
	user.Name = name
}

func (user *User) UpdateEmail(email string) {
	user.email = email
}

func (user *User) Rename(name string) {
	user.Name = name
}

func main() {
	user := User{Name: "Ada", age: 28, email: "123.com"}
	user.RenameCopy("Grace")
	fmt.Println("value receiver:", user.Name)
	user.Rename("Grace") // 对可取地址的变量，Go 自动使用 &user 调用指针接收者方法。
	fmt.Println("pointer receiver:", user.Name)
	copied := user
	copied.Name = "Lin"
	fmt.Println(user.Greeting(), "copy:", copied.Name)
	var zero User
	fmt.Printf("zero struct: name=%q age=%d\n", zero.Name, zero.age)
	user.UpdateEmail("2585671067@qq.com")
	fmt.Println("email:", user.email)
}
