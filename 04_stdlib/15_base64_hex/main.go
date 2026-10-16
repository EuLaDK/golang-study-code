// 是什么：Base64 和十六进制是把字节表示为文本的编码方式，可以解码还原，不会隐藏内容。
// 作用：让二进制数据适合放进文本格式；Base64 通常比十六进制更紧凑，URL-safe 变体避免部分 URL 特殊字符。
// 使用场景：传输二进制字段、展示摘要或把安全随机字节表示为令牌；编码本身不提供加密或身份验证。
//
// 目标：将二进制编码为可传输的文本；Base64 与十六进制都不是加密。
// 运行（仓库根目录）：go run ./04_stdlib/15_base64_hex
// 观察：StdEncoding 可能带 = 填充；RawURLEncoding 适合 URL 中的随机令牌。
// 练习：传入非法 Base64，读取 DecodeString 返回的 error。
package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

func main() {
	data := []byte("Go!")
	encoded := base64.StdEncoding.EncodeToString(data)
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		panic(err)
	}
	fmt.Println("base64:", encoded, "decoded:", string(decoded))
	fmt.Println("hex:", hex.EncodeToString(data))
	fmt.Println("URL-safe:", base64.RawURLEncoding.EncodeToString([]byte{251, 255, 254}))
	// 字节转字符串是原样解释；编码为 Base64 则改变表示方式，能被任何人解码。
}
