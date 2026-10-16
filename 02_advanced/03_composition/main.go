/*
是什么：组合是把一个类型作为另一个类型的字段；嵌入是省略字段名的一种组合写法。
作用：让外层类型直接访问内层的部分字段和方法，这叫提升；外层定义同名方法时，直接调用会使用外层方法。
使用场景：在已有 User 数据上增加 Admin 的角色信息，或复用公共行为；仍可通过 admin.User 显式访问内层值。

目标：通过嵌入组合类型，理解方法提升与显式访问。
对照：Go 没有 extends；嵌入使内部类型的方法可从外部访问，并不会产生继承关系。
运行：go run ./02_advanced/03_composition
观察：可直接调用 admin.Display()；管理员字段仍是自身字段；显式 admin.User.Display() 结果相同。
练习：为 Admin 添加同名 Display 方法，比较 admin.Display() 与 admin.User.Display()。
*/
package main

import "fmt"

type User struct {
	Name string
}

func (user User) Display() string { return "user: " + user.Name }

type Admin struct {
	User // 没有写名字直接嵌入
	Role string
}

// 当你把 User 嵌入进 Admin 后，Go 编译器会默默帮你在外层做一次“代理转发”：
// 你访问 admin.Name
//  编译器自动帮你翻译成：admin.User.Name
// 你调用 admin.Display()
//  编译器自动帮你翻译成：admin.User.Display()
// 你完全可以显式地调用 admin.User.Display()，两者的执行结果和底层汇编是一模一样的。

func (admin Admin) Display() string {
	return admin.User.Display() + " admin" // 这里必须要显示指定，否则会直接栈溢出
}

// ：因为 User 是匿名嵌入到 Admin 中的，Go 编译器会自动把内部 User 的所有方法和属性“提升”到外层的 Admin 上。admin.Display() 本质上就是 admin.User.Display() 的语法糖！
func main() {
	admin := Admin{User: User{Name: "Ada"}, Role: "moderator"}
	fmt.Println(admin)
	fmt.Println(admin.Display())                                    //  编译器自动帮你翻译成：admin.User.Display()
	fmt.Println("promoted field:", admin.Name, "role:", admin.Role) // 编译器自动帮你翻译成：admin.User.Name
	fmt.Println("explicit:", admin.User.Display())
	// Admin 不是 User 的子类；接收 User 的函数需要显式传 admin.User。
	printUser(admin.User)
	fmt.Println("这是admin的Display方法", admin.Display())
}

func printUser(user User) {
	fmt.Println("received User:", user.Name)
}
