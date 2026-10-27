package main

// 表驱动测试的优点
// 1、简单和紧凑
// 2、数据即测试
// 3、结合子测试后，可单独运行某个数据项的测试

// 表驱动测试实践中的注意事项
// 1、表的实现方式 表也可以使用基于自定义结构体的其他集合类型（如map）来实现，map是无序的所以表内数据项的先后不固定
// 2、测试失败时的数据项定位 通过输出数据表项在数据表中的偏移量来辅助定位“元凶”
// 3、Errorf还是Fatalf 如果一个数据项导致的测试失败不会对后续数据项测试结果造成影响，推荐Errof；如果数据项导致的测试失败会直接影响到后续数据项的测试结果，那么可以使用Fatalf让测试尽快结束

import (
	"strings"
	"testing"
)

func TestCompare(t *testing.T) { // 自定义结构体切片
	compareTests := []struct {
		name, a, b string
		i          int
	}{
		{"compareTwoEmptyString", "", "", 0},
		{"compareSecondParamIsEmpty", "a", "", 1},
		{"compareFirstParamIsEmpty", "", "a", -1},
	}

	for _, tt := range compareTests {
		t.Run(tt.name, func(t *testing.T) {
			cmp := strings.Compare(tt.a, tt.b)
			if cmp != tt.i {
				t.Errorf("want %v, but Compare(%q, %q) = %v", tt.i, tt.a, tt.b, cmp)
			}
		})
	}
}

func TestCompareMap(t *testing.T) {
	compareTests := map[string]struct {
		a, b string
		i    int
	}{
		`compareTwoEmptyString`:     {"", "", 0},
		`compareSecondParamIsEmpty`: {"a", "", 1},
		`compareFirstParamIsEmpty`:  {"", "a", -1},
	}
	for name, tt := range compareTests {
		t.Run(name, func(t *testing.T) {
			cmp := strings.Compare(tt.a, tt.b)
			if cmp != tt.i {
				t.Errorf(`want %v, but Compare(%q, %q) = %v`, tt.i, tt.a, tt.b, cmp)
			}
		})
	}
}
