// 是什么：SHA-256 把内容算成固定长度摘要；HMAC 用共享密钥计算消息认证值；crypto/rand 生成安全随机字节。
// 作用：摘要可用于比对内容，HMAC 用于检测篡改和验证对方持有密钥，安全随机字节可生成难以猜测的令牌。
// 使用场景：文件完整性比对、回调消息认证和令牌生成；这些机制不隐藏消息，也不能直接替代专用的密码存储方案。
//
// 目标：SHA-256 摘要、HMAC 消息认证、安全随机令牌各自解决什么问题。
// 运行（仓库根目录）：go run ./04_stdlib/21_crypto
// 观察：摘要确定，令牌每次不同；篡改消息会使 HMAC 校验失败。
// 练习：修改 message 内容，观察摘要和认证值如何变化。
package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

func sign(key, message []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(message) // hash.Hash 的 Write 按接口约定永不返回 error。
	return mac.Sum(nil)
}

func main() {
	message := []byte("order=42&amount=100")
	digest := sha256.Sum256(message)
	fmt.Printf("SHA-256: %x\n", digest)
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		panic(err)
	}
	signature := sign(key, message)
	valid := hmac.Equal(signature, sign(key, message))
	tampered := hmac.Equal(signature, sign(key, []byte("order=42&amount=1")))
	if !valid || tampered {
		panic("HMAC self-check failed")
	}
	fmt.Println("valid:", valid, "tampered valid:", tampered)
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		panic(err)
	}
	fmt.Println("random token:", base64.RawURLEncoding.EncodeToString(token))
	// SHA-256 不是加密，也不适合直接存储密码；HMAC 验证使用 hmac.Equal。
	// 本课不实现密码存储；真实项目使用专用密码哈希方案并安全管理密钥。
	// 不输出 key；此处打印 token 只为演示，真实服务不得把访问令牌写进日志。
}
