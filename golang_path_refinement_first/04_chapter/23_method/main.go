package main

import (
	"fmt"
	"time"
)

// Go 方法的一般声明形式如下：
// func (receiver T/*T) MethodName(参数列表)(返回值列表){
// 	// 方法体
// }
// 上面方法声明中的T称为receiver的基类型。通过receiver，上述方法被绑定到类型T上。换句话说，上述方法是类型T的一个方法，我们可以通过类型T或*T的实例调用该方法
// var t T
// t.MedthodName(参数列表)
// var pt *T = &t
// pt.MedthodName(参数列表)
// Go方法具有如下特点
// 1、方法的首字母是否大写决定了该方法是不是导出方法
// 2、方法定义要和类型定义放在同一个包内。由此我们可以推出：不能为原生类型（如int、float64等）添加方法，只能为自定义类型定义方法，不能横跨Go包为其他包内的自定义类型定义方法
// 3、每个方法只能有一个receiver参数，不支持多receiver参数列表或变长receiver参数，一个方法只能绑定一个基类型
// 4、receiver参数的基类型本身不能是置身类型或接口类型

type T struct {
	a int
}

func (t T) Get() int {
	return t.a
}

func (t *T) Set(x int) int {
	t.a = x
	return t.a
}

// 上面的函数等价于下面这两个
func Get(t T) int {
	return t.a
}
func Set(t *T, a int) int {
	t.a = a
	return t.a
}

// Go方法本质：一个以方法所绑定类型实例为第一个参数的普通函数

// 选择正确的receiver类型
// 函数和方法的等价变换公式：func (t T) M1() <=> M1(t T)  func(t *T)M2() <=>M2(t *T)
// 当receiver参数的类型为T时，选择值类型的receiver
// 当receiver参数的类型为*T时，选择指针类型的receiver

type T1 struct {
	a int
}

func (t T1) M1() {
	t.a = 100
}

func (t *T1) M2() {
	t.a = 200
}

type field struct {
	name string
}

func (p *field) print() {
	fmt.Println(p.name)
}

func main() {
	var t T
	fmt.Println(t.Get())
	fmt.Println(t.Set(1))

	// 所以可以这样调用
	// 这种直接以类型名T调用方法的表达方式被称为方法表达式
	// （Method Expression）。类型T只能调用T的方法集合（Method Set）
	// 中的方法，同理，*T只能调用*T的方法集合中的方法
	var t1 T
	fmt.Println(T.Get(t1))
	fmt.Println((*T).Set(&t1, 1))

	var t2 T1
	fmt.Println(t2.a)
	t2.M1()
	fmt.Println(t2.a)
	t2.M2() // 选择以*T作为receiver参数类型时，T的M2方法等价为M2(t *T)。我们传递给M2函数的t是T类型实例的地址，这样M2函数体中对参数t做的任何修改都会反映到原T类型实例上
	fmt.Println(t2.a)

	// 案例
	data1 := []*field{{"one"}, {"two"}, {"three"}}
	for _, v := range data1 {
		go v.print()
	}

	data2 := []field{{"four"}, {"five"}, {"six"}}
	for _, v := range data2 {
		go v.print()
	}
	time.Sleep(3 * time.Second)
	nums := []int{10, 20, 30}
	for _, v := range nums {
		fmt.Printf("值: %d, 内存地址: %p\n", v, &v)
	}
}
