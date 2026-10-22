package main

import "fmt"

// 在Go中，只有在函数和方法内部才能使用defer，这些函数被称为deferred函数。defer将他们注册到其所在goroutine用于存放deferred函数的栈数据结构中，这些deferred函数将在执行defer的函数退出前被按后进先出的顺序调度执行
// defer 常见用法
// 1.拦截panic

func makeSlice(n int) []byte {
	defer func() {
		if recover() != nil {
			panic("程序错误")
		}
	}()
	return make([]byte, n)
}

func bar() {
	fmt.Println("raise a panic")
	panic(-1)
}

func foo() {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("recovered from a panic", e)
		}
	}()
	bar()
}

// defer修改函数返回值
func foo1(a, b int) (x, y int) {
	defer func() {
		x = x * 5
		y = y * 10
	}()
	x = a + 5
	y = b + 10
	return // 有具名返回值直接return就行，等价于 return x,y
}

// 输出调试信息
func trace(s string) string {
	fmt.Println("entering:", s)
	return s
}

func un(s string) {
	fmt.Println("leaving:", s)
}

func a() {
	defer un(trace("a"))
	fmt.Println("in a")
}

func b() {
	defer un(trace("b"))
	fmt.Println("in b")
	a()
}

// defer关键字后面的表达式是在将deferred函数注册到deferred函数栈的时候进行求值的
func foo2() {
	for i := 0; i <= 3; i++ {
		defer fmt.Println(i)
	}
}

func foo3() {
	for i := 0; i <= 3; i++ {
		defer func(n int) {
			fmt.Println(n)
		}(i)
	}
}

func foo4() {
	for i := 0; i <= 3; i++ {
		defer func() {
			fmt.Println(i)
		}()
	}
}

func foo5() {
	s1 := []int{1, 2, 3}
	defer func(a []int) {
		fmt.Println(a)
	}(s1)

	s1 = []int{3, 2, 1} // 这个是在内存里面再开辟一个新地址并把这个新的地址传给s1
	_ = s1
}

func foo6() {
	s1 := []int{1, 2, 3}
	defer func(p *[]int) {
		fmt.Println(*p)
	}(&s1)
	s1 = []int{3, 2, 1}
	_ = s1
}
func foo5_modified() {
	s1 := []int{1, 2, 3}

	defer func(a []int) {
		fmt.Println(a) // 这里会输出 [666, 2, 3] ！！！
	}(s1)

	// 注意看！这里没有换仓库（没有 s1 = ...）
	// 而是拿着便签 s1 找到仓库，把里面的第 0 块砖改了！
	s1[0] = 666
}

func main() {
	foo()
	fmt.Println("main exit normally")
	x, y := foo1(1, 2)
	fmt.Println("x=", x, "y=", y)
	b()
	fmt.Println("foo2 result:")
	foo2()
	fmt.Println("\nfoo3 result:")
	foo3()
	fmt.Println("\nfoo4 result:")
	foo4()
	fmt.Println("\n-------------------")
	foo5()
	foo5_modified()
	foo6()
}

// 提出的疑问 还是没太理解不是[2]int{3,4}像这样就是定死的数组长度才传的是值吗就是说会新建一个副本，切片都是指向外面的地址
// 我完全理解你现在的困惑点在哪里了！
// 你脑子里的逻辑是这样的：
// 数组 [2]int：笨重，传参时把里面的数据一个一个全部复制一遍（所以是传值）。
// 切片 []int：灵活，传参时只传一个指向底层的地址（所以它应该是指针，改了就全改）。
// 你这个想法 90% 是对的！但你漏掉了一个致命的细节：
// “指向底层数组的地址” 这句话没错，但这个【地址】本身，也是作为一个值被复制过去的！
// 我们把它们拆开来看，用最直白的图解来对比：
// 对比一：数组 vs 切片 传参到底复制了什么？
// 1. 数组 [3]int{1, 2, 3} 传参：
// 复制的是整个仓库的所有货物。
// 原数组：有 3 块砖。
// 传参后：在内存里重新搬来 3 块一模一样的砖给函数。成本很大。
// 2. 切片 []int{1, 2, 3} 传参：
// 确实不会复制砖！它只复制写着仓库门牌号的小便签。
// 原切片 s1：是一张便签，上面写着 仓库地址：0x100。
// 传参给 a：把便签上的字，抄了一份新的便签给 a，上面也写着 仓库地址：0x100。
// 这就是你说的“它们都指向外面的地址”，到这一步完全没错！
// 那为什么 foo5 还是输出了原来的值？
// 核心矛盾爆发在这里！ 请看这一行代码：
// code
// Go
// s1 = []int{3, 2, 1}
// 很多人以为这行代码是：把原来仓库里的砖砸了，换成 3, 2, 1。
// 完全不是！
// []int{3, 2, 1} 在 Go 里的意思是：
// 在内存的另一个地方，开辟了一个全新的、装有 3, 2, 1 的新仓库（地址是：0x999）。
// 然后，把你手里的原便签 s1 擦掉，写上新地址 0x999。
// 现在请看两张便签的现状：
// 你的便签 s1：地址变成了 0x999（新仓库）
// defer 手里的便签 a：还留着当初抄过去的旧地址 0x100（旧仓库，里面依然是 1, 2, 3）
// 所以 defer 执行的时候，顺着它手里的便签 0x100 找过去，打印出来的当然还是 [1, 2, 3]！
// 如果怎么写，才会影响到 defer？
// 只有你不换仓库，而是真的顺着地址去改仓库里的数据：
// code
// Go
// func foo5_modified() {
// 	s1 := []int{1, 2, 3}

// 	defer func(a []int) {
// 		fmt.Println(a) // 这里会输出 [666, 2, 3] ！！！
// 	}(s1)

// 	// 注意看！这里没有换仓库（没有 s1 = ...）
// 	// 而是拿着便签 s1 找到仓库，把里面的第 0 块砖改了！
// 	s1[0] = 666
// }
// 因为便签 a 和便签 s1 指向同一个仓库，所以你顺着 s1 去把仓库里的砖给改了，a 回头看这个仓库的时候，砖确实变了！
// 一句话让你彻底通透：
// 切片确实是指向外面地址的。
// 但是你写的 s1 = []int{3, 2, 1} 不是在修改那个地址里的内容，而是建了新切片，换掉了 s1 这个变量保存的地址。
// 而 defer 早就把你原来的那个旧地址给复制拿走了，你后来爱换成啥地址，关人家 defer 屁事呢！
