/*
是什么：map 是按键查找值的集合，例如 map[string]int 用字符串作为键、整数作为值，类似字典。
作用：保存、查找和删除键值对；value, ok := mapping[key] 中的 ok 可以区分“键不存在”和“值刚好是零值”。
使用场景：统计出现次数、按 ID 查找数据、保存配置；map 遍历顺序没有保证，共享读写时需要同步。

目标：map 的创建、存在性判断、删除和遍历。
对照：更接近 JS Map / Python dict；map 的迭代顺序不保证稳定。
运行：go run ./01_basics/09_maps
观察：missing 返回 0 false；副本变量写入后原 map 可见；按排序后的键稳定输出。
练习：去掉排序直接 range map，解释为什么不能依赖输出顺序。
*/
package main

import (
	"fmt"
	"sort"
)

func main() {
	scores := map[string]int{"Ada": 90, "Lin": 0}
	value, exists := scores["missing"]
	fmt.Println("missing:", value, exists)
	value, exists = scores["Lin"]
	fmt.Println("zero score:", value, exists)

	shared := scores
	shared["Ada"] = 95
	delete(scores, "Lin")
	keys := make([]string, 0, len(scores))
	for key := range scores {
		keys = append(keys, key)
	}
	sort.Strings(keys) // Go range map 不是按插入顺序、不是按键排序顺序遍历。map 是哈希表，每次遍历起点随机,所以用sort排个序输出
	for _, key := range keys {
		fmt.Println(key, scores[key])
	}
	var uninitialized map[string]int
	fmt.Println("nil read:", uninitialized["x"])
	uninitialized = make(map[string]int) // nil map 可读不可写，写入前要初始化。
	uninitialized["x"] = 1
	fmt.Println("initialized:", uninitialized["x"])
}
