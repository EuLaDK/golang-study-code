// 是什么：JSON 是前后端常用的数据交换格式，encoding/json 在 JSON 字节和 Go 值之间转换。
// 作用：Marshal 编码，Unmarshal/Decoder 解码；结构体上的 json 标签决定字段名和省略规则，解析后再做业务校验。
// 使用场景：读取 API 请求体、生成响应或处理 JSON 配置；未知字段、额外 JSON 值和过大输入需按接口要求处理。
//
// 目标：结构体标签、导出字段、JSON 编解码与严格请求解析。
// 运行（仓库根目录）：go run ./04_stdlib/13_json
// 观察：字段名受 json 标签控制，空 omitempty 字段省略；未知字段被拒绝。
// 练习：加入 Age 字段，并在解析后校验 0..150，区分格式错误和业务错误。
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Note string `json:"note,omitempty"`
}

func decodeUser(input string) (User, error) {
	var user User
	if len(input) > 4096 {
		return user, errors.New("JSON exceeds 4096 bytes")
	}
	decoder := json.NewDecoder(strings.NewReader(input))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&user); err != nil {
		return user, err
	}
	// Decode 只读取一个 JSON 值；再次读取并要求 EOF，拒绝拼接的第二份数据。
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		return User{}, errors.New("expected exactly one JSON value")
	}
	if user.ID <= 0 || strings.TrimSpace(user.Name) == "" {
		return User{}, errors.New("positive id and non-empty name required")
	}
	return user, nil
}

func main() {
	user, err := decodeUser(`{"id":1,"name":"Ada"}`)
	if err != nil {
		panic(err)
	}
	data, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
	_, err = decodeUser(`{"id":1,"name":"Ada","admin":true}`)
	fmt.Println("unknown field:", err)
	// 小写未导出字段不能被 encoding/json 访问；标签不会突破可见性。
	// 解到 any 时数字默认 float64；大整数优先用明确的 int64/string 字段。
	// DisallowUnknownFields 仍不拒绝重复键或大小写变体；不要视为完整模式校验器。
}
