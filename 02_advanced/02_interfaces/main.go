/*
是什么：interface 定义一组方法，表达“这个值能做什么”；类型只要具有所需方法，就能隐式满足该接口。
作用：让不同类型通过相同的方法被调用；类型断言和类型 switch 用来检查接口中实际保存的类型。
      接口同时保存动态类型和值，所以它装着一个 nil 指针时，接口本身也可能不等于 nil。
使用场景：统一处理不同数据来源、输出目标或业务实现，例如 io.Reader；在需要替换行为时使用小接口。

目标：隐式实现接口、类型断言、类型 switch、接口的动态类型/值。
对照：类似 TS 的结构化类型，但 Go 接口匹配方法集；无需 implements 关键字。
运行：go run ./02_advanced/02_interfaces
观察：person 与 service 都可作为 Greeter；typed nil interface 为 false，而 nil pointer 为 true。
练习：增加第三种 Greeter 实现；把 Person 的方法改为指针接收者并修正编译错误。
*/
package main

import "fmt"

type Greeter interface {
	Greet() string
}

type Person struct{ Name string }

func (person *Person) Greet() string {
	if person == nil {
		return "Person is nul"
	}
	return "Hello, " + person.Name
}

type Service string

func (service Service) Greet() string { return "Service: " + string(service) }

// 2. 增加第三种 Greeter 实现：Robot
type Robot struct{ Model string }

func (robot Robot) Greet() string { return "Beep boop, I am " + robot.Model }

func describe(value any) { // any 是 interface{}
	switch typed := value.(type) { // 类型分支 switch ... .(type) 语法糖
	case string:
		fmt.Println("string:", typed)
	case int:
		fmt.Println("int:", typed)
	default:
		fmt.Printf("other: %T\n", typed)
	}
}

func main() {
	greeters := []Greeter{&Person{Name: "Ada"}, Service("users"), Robot{Model: "T-800"}}
	fmt.Println(greeters)
	for _, greeter := range greeters {
		fmt.Println(greeter.Greet())
		if person, ok := greeter.(*Person); ok { // greeter.(Person) 类型断言
			fmt.Println("asserted person:", person.Name)
		}
	}
	describe("Go")
	describe(42)
	describe(true)
	var pointer *Person
	var greeter Greeter = pointer
	fmt.Printf("%T\n", greeter)
	fmt.Println("nil pointer:", pointer == nil, "typed nil interface:", greeter == nil) // 判断 greeter == nil 的唯一标准是：必须这两个格子(类型指针，数据指针)都为nil，这个接口才算nil
	// 此时不要调用 greeter.Greet()：值接收者需要解引用 nil 指针，会 panic。
	fmt.Println("Safe call on nil:", greeter.Greet())
}
