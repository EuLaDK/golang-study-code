package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	// 忽略对服务端证书的校验
	// tr := &http.Transport{
	// 	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	// }
	// client := &http.Client{Transport: tr}
	// resp, err := client.Get("https://localhost:8081")

	// 加载证书
	pool := x509.NewCertPool()
	caCertPath := "../server/server.crt"

	caCrt, err := os.ReadFile(caCertPath)
	if err != nil {
		fmt.Println("ReadFile err:", err)
		return
	}
	pool.AppendCertsFromPEM(caCrt)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{RootCAs: pool},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get("https://localhost:8081")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Println(string(body))
}
