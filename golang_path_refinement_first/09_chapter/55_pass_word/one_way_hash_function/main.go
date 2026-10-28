package main

import (
	"crypto"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
)

// 单向散列函数（one-way hash function）是一个接受不定长输入但产生定长输出的函数，这个定长输出被称为“摘要”（digest）或“指纹”（fingerprint）

// SHA256
func sum256(data []byte) string {
	sum := sha256.Sum256(data)
	return fmt.Sprintf("%x", sum)
}

func main() {
	s := "I love go programming language!!"
	fmt.Println(sum256([]byte(s)))

	// 使用SHA-256单向散列函数的HMAC示例
	key := []byte("12345678123456781234567812345678")
	// 要传递的消息
	message := []byte("I love go programming language!!")

	// 创建hmac实例（使用SHA-256单向散列函数）
	win := hmac.New(sha256.New, key)
	win.Write(message)

	// 计算mac值
	m := win.Sum(nil)
	ms := fmt.Sprintf("%x", m) // mac 到 string
	fmt.Printf("mac值 = %s\n", ms)

	// 使用RSA-PSS进行消息签和验证的示例
	// 1. 生成私钥
	privatekey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	publicKey := privatekey.PublicKey

	// 2. 待签名消息
	msg := []byte("I love go programming language!!")

	// 3. 计算摘要
	digest := sha256.Sum256(msg)

	// 4. 用私钥签名 (注意补齐 5 个参数：digest[:] 和 nil)
	sign, err := rsa.SignPSS(rand.Reader, privatekey, crypto.SHA256, digest[:], nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("签名：%x\n", sign)

	// 5. 用公钥验证签名 (注意补齐 5 个参数：最后一个 opts 传 nil)
	err = rsa.VerifyPSS(&publicKey, crypto.SHA256, digest[:], sign, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("签名验证成功！")

	// 随机数 需要多长的输入字节长度即可
	c := 32
	b := make([]byte, c)
	_, errb := rand.Read(b)
	if errb != nil {
		panic(errb)
	}
	fmt.Printf("随机数%x\n", b)
}
