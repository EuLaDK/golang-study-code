package main

import "fmt"

// map的key不能是函数，map，切片

// 使用复合字面值创建map类型变量
var statusText = map[int]string{
	200: "OK",
	201: "Created",
	202: "Accepted",
}

// 使用make函数创建map类型变量
// icookies = make(map[string][]*Cookie)

func foo(m map[string]int) {
	m["key1"] = 20
	m["key2"] = 25
}

// 如果你需要一个稳定的遍历次序，那么一个比较通用的做法是使用另一种数据结构来按需要的次序保存key
func doIteration(sl []int, m map[int]int) {
	fmt.Printf("{")
	for _, k := range sl {
		// 按切片中的元素次序迭代
		v, ok := m[k]
		if !ok {
			continue
		}
		fmt.Printf("[%d, %d]", k, v)
	}
	fmt.Printf("}\n")
}

// 这里延续运行时字段解释 把图看完来看这里
// 真正用来存储键值对数据的是bucket（桶），每个bucket中存储的Hash值低bit位数值相同的元素，默认每个元素个数值为8
// 当某个bucket的8个空槽（slot）都已填满且map尚未达到扩容条件时，运行时会建立ovewflow bucket，并将该overflow bucket挂在上面bucket末尾的overflow指针上，这样两个bucket形成了一个链表结构，该结构的存在将持续到下一次map扩容
// 每个bucket又三部分组成：tophash区域、key存储区域和value存储区域

// map扩容
// Go在运行时的map实现中引入了一个LoadFactor（负载因子），当count > LoadFactor*2或 overflow bucket过多时，运行时会对map进行扩容
// 如果是因为overflow bucket过多导致的扩容，实际上运行时会新建一个和现有规模一样的bucket数组，然后再进行assign和delete操作时进行排空和迁移；如果是因为当前输出数量超出LoadFacto指定的水位的情况
// 那么运行时会建立一个两倍于现有规模的bucket数组，但真正的排空和迁移工作也是在进行assign和delete操作时逐步进行的。原bucket数组会挂在hmap的oldbuckets指针下面，直到原buckets数组中所有数组都迁移到新数组，原buckets数组才会被释放

// 尽量使用cap参数创建map 主要是防止频繁扩容影响性能
func main() {
	fmt.Println(statusText)
	m1 := map[string]int{
		"key1": 200,
		"key2": 250,
	}
	fmt.Println(m1)
	foo(m1)
	fmt.Println(m1)
	m1["key3"] = 30
	v, ok := m1["key3"]
	if !ok {
		fmt.Println("key3不在m1中")
	}
	fmt.Println(v)
	for k, v := range m1 {
		fmt.Println(k, v)
	}
	delete(m1, "key2")
	fmt.Println(m1)
	var sl []int
	m := map[int]int{
		1: 11,
		2: 12,
		3: 13,
	}
	for k, _ := range m {
		sl = append(sl, k)
	}
	for i := 0; i < 3; i++ {
		doIteration(sl, m)
	}
}
