// 是什么：maps 是处理 map 集合的标准库工具，提供复制键值、比较内容等操作。
// 作用：Clone 创建浅拷贝，Copy 将来源键值写入目标，Equal 比较键值；浅拷贝不会递归复制值中引用的数据。
// 使用场景：复制配置或对比缓存数据；需要固定输出顺序时仍要收集并排序键，共享修改时仍要同步。
//
// 目标：maps.Clone/Equal/Copy；map 的浅拷贝和无序遍历。
// 运行（仓库根目录）：go run ./04_stdlib/08_maps
// 观察：普通赋值共享 map；Clone 独立复制键值，但不会深拷贝值中的引用。
// 练习：把值改成 []string，观察 Clone 后修改切片元素是否互相影响。
package main

import (
	"fmt"
	"maps"
	"slices"
)

func main() {
	original := map[string]int{"go": 1, "http": 2}
	copyMap := maps.Clone(original)
	copyMap["go"] = 3
	fmt.Println("equal:", maps.Equal(original, copyMap), "original go:", original["go"])
	maps.Copy(copyMap, map[string]int{"json": 4})
	keys := make([]string, 0, len(copyMap))
	for key := range copyMap {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	for _, key := range keys {
		fmt.Println(key, copyMap[key])
	}
	// 不依赖 map 遍历顺序；稳定输出先收集键再排序。并发写 map 需要同步。
}
