package main

import (
	"fmt"
	"unsafe"
)

// map[key_type]value_type
// key和value的类型可以相同，也可以不同
// map[string]string  map[int]string 如果两个map的key类型相同，value类型也相同，那么可以说他们是同一种map类型
// map对value类型没有限制，但是对key的类型有严格要求，因为map类型要保证key的唯一性，Go要求，key的类型必须支持 == 和 != 两种比较操作符
// 函数类型、map 类型自身和切片类型不能作为 map 的 key 类型

// 结构体
// 在 Go中定义一个新类型的方法有两种
// 第一种 类型定义
type MyInt int64 //  type T S 这里的S可以是任意的已知的类型，包括Go的底层原生类型，自己定义的其它类型
type MyType []string
type ( // 这样也可以给类型起一个别名
	T1 int
	T2 string
	T3 bool
)

// 定义结构体类型
type T struct {
	Field1 T1
	Field2 T2
	Field3 T3
}

// 案例定义一个书的结构体类型
type Book struct {
	Name    string         `json:"name"`
	Pages   int            `json:"page"`
	Indexes map[string]int `json:"indexes"`
}

// 定义一个空结构体类型
type Empty struct{}

// 使用其它结构体作为自定义结构体中字段的类型
type Person struct {
	Name  string
	Phone string
	Addr  string
}

type Student struct {
	Age  int
	Info Person
}

type TT struct {
	// t TT // 不可以这样不可以又以自己为类型
	st *TT           // 可以以指针
	sl []TT          // 可以以切片
	sm map[string]TT // 可以
}

func main() {
	// map 是 Go 提供的一种抽象数据类型，用于表示一组无序的键值对

	// map 变量声明方式
	var m map[string]int // 一个map[string]int类型的变量 和切片一样默认值为nil
	// map无法实现零值可用就是nil的情况下不可去访问key的value值
	fmt.Println(m == nil)
	// 使用复合字面值初始化map类型变量
	// m := map[int]string{}
	m1 := map[int][]string{
		1: {"vall_1", "vall_2"},
		4: {"tom", "jerry"},
	}
	fmt.Println(m1[1])

	// 使用make为map类型变量进行初始化
	m2 := make(map[int]string) // 未指定初始容量
	m2[1] = "tom"
	m3 := make(map[int]string, 8) // 指定初始容量为8
	m3[1] = "jerry"
	fmt.Println(m2, m3)
	// 插入键值对
	m2[4] = "hello"
	// 如果key有值就会覆盖之前的值
	m2[1] = "111"
	fmt.Println(m2, "m2", m2[5])
	// 获取键值对数量len函数 不能对 map 类型变量调用 cap 函数来获取当前容量
	fmt.Println(len(m2))

	// 通过comma ok的惯用法对某个key进行查询
	v, ok := m2[6]
	if !ok {
		fmt.Println("6不在map-m2中")
	} else {
		fmt.Println(v)
	}
	// 删除数据 内置函数delete
	// 第一个参数是 map 类型变量，第二个参数是要删除的键
	delete(m2, 4) // 删除 key==4
	fmt.Println(m2)
	// 遍历键值 在Go中遍历键值只有for range一种方法
	m3[2] = "tom"
	m3[3] = "marry"
	for k, v := range m3 {
		fmt.Println(k, v)
	}
	// 如果只想用某一个 可以用 _ 代替 例如 for _, v := range m3
	// map是无序的，每次遍历都可能出现新的结果

	// map的传递开销，map每次传递的是指针传递而非值传递，所以开销比较小，这也是为什么两个变量指向一个map的话修改一个会影响另一个

	// 结构体
	book := Book{
		Name:  "计算机组成原理",
		Pages: 345,
		Indexes: map[string]int{
			"序言": 1,
			"结束": 345,
		},
	}
	fmt.Println(book)
	var b Book
	b.Name = "玄霸九天"
	b.Pages = 1927
	fmt.Println(b)
	// 定义一个空结构体类型
	var s Empty
	fmt.Println(unsafe.Sizeof(s))
	// 空的结构体类型变量大小为0，即内存占用为0，由于空结构体类型的变量不占用内存空间，他们常被用“事件”在goroutine之间进行通信
	// 例如以下代码
	// var c = make(chan Empty) // 声明一个元素类型为Empty的channel
	// c <- Empty{}             // 向channel 写入一个事件 这种方式创建的以空结构体为元素的channel，是Go中实现最小的goroutine的通信方式之一

	stu := Student{
		Age: 22,
		Info: Person{
			Name:  "小明",
			Phone: "1234567",
			Addr:  "临渊城许家",
		},
	}
	fmt.Println(stu, stu.Info.Phone)
}
