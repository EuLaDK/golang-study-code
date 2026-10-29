package main

import (
	"fmt"
	"reflect"
)

// reflect.TypeOf和reflect.ValueOf是进入反射世界仅有的两扇“大
// 门”。通过reflect.TypeOf这扇“门”进入反射世界，你将得到一个
// reflect.Type对象，该对象中包含了被反射的Go变量实例的所有类型信
// 息；而通过reflect.ValueOf这扇“门”进入反射世界，你将得到一个
// reflect.Value对象。Value对象是反射世界的核心，不仅该对象中包含了
// 被反射的Go变量实例的值信息，而且通过调用该对象的Type方法，我
// 们还可以得到Go变量实例的类型信息，这与通过reflect.TypeOf获得类型
// 信息是等价的

// 反射世界入口可以获取Go变量实例的类型信息和值信息的关键在于，它们利用了interface{}类型的形式参数对传入的实际参数（Go变量实例）的析构能力。两个入口函数分别将得到的值信息和类型信息存储在reflect.Value对象和reflect.Type对象中

// 通过反射对象，我们还可以调用函数或对象的方法
func Add(i, j int) int {
	return i + j
}

type Calculator struct{}

func (c Calculator) Add(i, j int) int {
	return i + j
}

// 反射世界的出口
// reflect.Value.Interface()是reflect.ValueOf()的逆过程，通过Interface方
// 法我们可以将reflect.Value对象恢复成一个interface{}类型的变量值。

// eflect.Value提供了CanSet、CanAddr及CanInterface等方法来帮助我们判断反射对象是否可设置（Settable）、可寻址、可恢复为一个interface{}类型变量

// 当被反射对象以值类型（T）传递给reflect.ValueOf时，所得到的反射对象（Value）是不可设置和不可寻址的
// 当被反射对象以指针类型（*T或&T）传递给reflect.ValueOf时，通过reflect.Value的Elem方法可以得到代表该指针所指内存对象的Value反射对象。而这个反射对象是可设置和可寻址的，对其进行修改（比如利用Value的SetInt方法）将会像函数的输出参数那样直接修改被反射对象所指向的内存空间的值
// 当传入结构体或数组指针时，通过Field或Index方法得到的代表结构体字段或数组元素的Value反射对象也是可设置和可寻址的。如果结构体中某个字段是非导出字段，则该字段是可寻址但不可设置的
// 当被反射对象的静态类型是接口类型时，该被反射对象的动态类型决定了其进入反射世界后的可设置性。如果动态类型为*T或&T时，那么通过Elem方法获得的反射对象就是可设置和可寻址的
// map类型被反射对象比较特殊，它的key和value都是不可寻址和不可设置的。但我们可以通过Value提供的SetMapIndex方法对map反射对象进行修改，这种修改会同步到被反射的map变量中

func main() {
	// reflect.ValueOf().Type() 等价于 reflect.TypeOf()
	var i int = 5
	val := reflect.ValueOf(i)
	typ := reflect.TypeOf(i)
	fmt.Println(reflect.DeepEqual(typ, val.Type()), typ.Name(), val.Int())

	var b = true
	val = reflect.ValueOf(b)
	typ = reflect.TypeOf(b)
	fmt.Println(typ.Name(), val.Bool())

	var f = 3.14
	val = reflect.ValueOf(f)
	typ = reflect.TypeOf(f)
	fmt.Println(typ.Name(), val.Float())

	var s = "hello, reflection"
	val = reflect.ValueOf(s)
	typ = reflect.TypeOf(s)
	fmt.Println(typ.Name(), val.String())

	var fn = func(a, b int) int {
		return a + b
	}
	val = reflect.ValueOf(fn)
	typ = reflect.TypeOf(fn)
	fmt.Println(typ.Kind(), typ.String())

	fun := reflect.ValueOf(Add)
	var m = 5
	var n = 5
	vals := []reflect.Value{reflect.ValueOf(m), reflect.ValueOf(n)}
	ret := fun.Call(vals)
	fmt.Println(ret[0].Int())

	// 方法调用
	c := reflect.ValueOf(Calculator{})
	mc := c.MethodByName("Add")
	ret = mc.Call(vals)
	fmt.Println(ret[0].Int(), ret, "ret")

	// 如果是值会创建一个新的地址放新值，所以修改新变量的值不会影响原来的变量的值
	var outi = 5
	valout := reflect.ValueOf(i)
	rout := valout.Interface().(int)
	fmt.Println(rout)
	rout = 6
	fmt.Println(outi, rout)

	// 如果是传的指针，那么他们共用一块内存地址，所以修改qout指针会影响outi指针
	valout = reflect.ValueOf(&i)
	qout := valout.Interface().(*int)
	fmt.Printf("%p, %p, %d\n", &i, qout, *qout)
	*qout = 7
	fmt.Println(outi)
}
