// 是什么：math 提供数学函数；math/rand/v2 生成伪随机数，即按算法产生、给定种子可以复现的数列。
// 作用：计算平方根、取整、比较浮点误差，以及生成模拟数据；PCG 是本例使用的伪随机数算法。
// 使用场景：制作测试数据、模拟概率或随机选择样本；访问令牌等需要不可预测的随机值，应使用 crypto/rand。
//
// 目标：math 常用函数，以及可复现的 math/rand/v2 模拟数据。
// 运行（仓库根目录）：go run ./04_stdlib/06_math_rand
// 观察：固定 PCG 种子每次给出相同序列；IntN 的上界不包含在结果中。
// 练习：模拟 100 次掷骰子，统计每个点数出现次数。
package main

import (
	"fmt"
	"math"
	"math/rand/v2"
)

func main() {
	fmt.Println("sqrt:", math.Sqrt(81), "ceil:", math.Ceil(2.1), "round:", math.Round(2.5))
	a, b := 0.1, 0.2
	fmt.Println("float equality:", a+b == 0.3, "within tolerance:", math.Abs((a+b)-0.3) < 1e-9)
	rng := rand.New(rand.NewPCG(1, 2))
	for i := 0; i < 5; i++ {
		fmt.Println("dice:", rng.IntN(6)+1)
	}
	// 金额通常用最小货币单位的整数；容差取决于业务，不能一律套用 1e-9。
	// math/rand 不适合令牌或密码；安全随机数见 21_crypto。
}
