package main

import "fmt"

var (
	a = c + b
	b = f()
	c = f()
	d = 3
)

var (
	a1 = c1 + b1
	b1 = ff()
	_  = ff()
	c1 = ff()
	d1 = 3
)

func f() int {
	d++
	return d
}
func ff() int {
	d1++
	return d1
}

func f1() int {
	fmt.Println("calling f")
	return 1
}

func g1(a, b, c int) int {
	fmt.Println("calling g", a, b, c)
	return 2
}

func h1() int {
	fmt.Println("calling h")
	return 3
}

func i1() int {
	fmt.Println("calling i")
	return 1
}

func j1() int {
	fmt.Println("calling j")
	return 1
}

func k1() bool {
	fmt.Println("calling k")
	return true
}

func main() {
	fmt.Println(a, b, c, d)
	fmt.Println(a1, b1, c1, d1)
	var y = []int{11, 12, 13}
	var x = []int{21, 22, 23}
	var c chan int = make(chan int)
	go func() {
		c <- 1
	}()
	y[f1()], _ = g1(h1(), i1()+x[j1()], <-c), k1()
	fmt.Println(y)
}
