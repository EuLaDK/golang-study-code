// 是什么：slices 提供切片的常用操作，cmp 帮助比较值，sort 提供传统排序和查找工具。
// 作用：排序、浅拷贝和查找；比较函数决定先后顺序，二分查找通过不断缩小区间查找已经排好序的数据。
// 使用场景：把结果按年龄、姓名或 ID 排列，以及在有序列表中查找；不希望原数据改变时先复制切片。
//
// 目标：排序、复制、查找切片；现代 Go 优先 slices/cmp，认识传统 sort。
// 运行（仓库根目录）：go run ./04_stdlib/07_sort_slices_cmp
// 观察：Sort 原地修改；Clone 后排序不影响原切片；BinarySearch 需要已排序数据。
// 练习：按年龄降序排序；同龄时仍按姓名升序。
package main

import (
	"cmp"
	"fmt"
	"slices"
	"sort"
)

type person struct {
	Name string
	Age  int
}

func main() {
	original := []int{3, 1, 2}
	ordered := slices.Clone(original)
	slices.Sort(ordered)
	index, found := slices.BinarySearch(ordered, 2)
	fmt.Println("original:", original, "sorted:", ordered, "search:", index, found)
	people := []person{{"Bo", 25}, {"Ada", 25}, {"Cat", 20}}
	slices.SortFunc(people, func(a, b person) int {
		return cmp.Or(cmp.Compare(a.Age, b.Age), cmp.Compare(a.Name, b.Name))
	})
	fmt.Println(people)
	fmt.Println("first >= 2:", sort.SearchInts(ordered, 2))
	// JS 默认 sort 按字符串顺序；slices.Sort([]int) 按数值顺序。
}
