// 是什么：这个文件包含固定边界测试和模糊测试；种子是探索的起始输入，f.Fuzz 会接收工具生成的新输入。
// 作用：检查合法 ID 解析后再转回字符串、再次解析，值保持一致；非法输入返回 error 是正常结果。
// 使用场景：为外部字符串解析器探索前导零、超大数值或其他难以手工列全的组合。
// 运行：go test ./06_testing/04_fuzz 会运行固定测试与种子；主动探索用 -run '^$' -fuzz '^FuzzParseID$' -fuzztime 3s。
//
package main

import (
	"strconv"
	"testing"
)

func TestParseID(t *testing.T) {
	for _, input := range []string{"", "0", "-1", "1.5", "18446744073709551616"} {
		if _, err := parseID(input); err == nil {
			t.Errorf("parseID(%q) should fail", input)
		}
	}
	if got, err := parseID("42"); err != nil || got != 42 {
		t.Fatalf("got %d, err=%v; want 42, nil", got, err)
	}
}

func FuzzParseID(f *testing.F) {
	for _, seed := range []string{"1", "42", "0007", "18446744073709551615", "0", "-1", "", "abc"} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, raw string) {
		id, err := parseID(raw)
		if err != nil {
			return // 不合法的外部输入应返回错误；它本身不是测试失败。
		}
		if id == 0 {
			t.Fatal("successful parse returned zero")
		}
		canonical := strconv.FormatUint(id, 10)
		again, err := parseID(canonical)
		if err != nil || again != id {
			t.Fatalf("round trip %q -> %d -> %q -> %d, err=%v", raw, id, canonical, again, err)
		}
	})
}
