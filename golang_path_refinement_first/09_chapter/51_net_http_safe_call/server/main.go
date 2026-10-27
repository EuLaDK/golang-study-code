package main

import (
	"fmt"
	"net/http"
)

// HTTPS比HTTP增加了一个安全传输层。安全传输层通常采用SSL或TLS协议实现。这一层负责HTTP协议传输的内容加密、通信双方身份验证等 HTTPS协议就是在安全传输层上运行的HTTP协议
// HTTPS 安全传输层的工作机制

// 安全传输层建立连接的过程也称为“握手阶段”，这个阶段设计四轮通信
// 1、ClinetHello（客户端-> 服务端）客户端向服务端发出建立安全层传输连接，构建加密通信通道的请求。在这个请求中，客户端会向服务端提供本地最新TLS版本、支持的加密算法组合的集合

// 2、ServerHello & Server certificate & ServerKeyExchange（服务端->客户端）这一轮通信分为三个重要步骤
// 2.1、第一步服务端收到客户端发来的ClietnHello请求后，用客户端发来的信息与自己本地支持的TLS版本、加密算法组合的集合做比较，选出一个TLS版本和一个合适的加密算法组合，然后生成一个随机数，一起打包到ServerHello中返回给客户端
// 2.2、第二步 服务器会将自己的服务端公钥证书发送给客户端（Server certificate），这个服务端公钥证书身兼两大职责：客户端对服务端身份的验证以及后续双方会话密钥的协商和生成。如果服务端要验证客户端身份（可选的），那么这里服务端还会发送一个CertificateRequest的请求给客户端，要求对客户端的公钥证书进行验证
// 2.3、第三步 发送开启双方会话密钥协商的请求（ServerKeyExchange）。对称加密算法的性能要高出几个数量级，因此HTTPS在开始真正传输应用层的用户数据之前，选择了在非对称加密算法的帮助下协商一个基于对称加密算法的密钥。
// 在密钥协商环节，通常会使用到Diffie-Hellman（DH）密钥交换算法，这是一种密钥协商的协议，支持通信双发在不安全的通道上生成对称加密算法所需的共享密钥。因此在这个步骤的请求中，服务端会向客户端发送密钥交换算法的相关参数信息
// 最后，服务端以Server Finihed(又称为ServerDone) 作为该轮通信的结束标志

// 3、ClinetKeyExchange & ClientChangeCipher & Finished (客户端 -> 服务端)
// 客户端在收到服务端的公钥证书后会对服务端的身份进行验证（也可以不验证），如果验证失败，则此次连接失败。如果验证通过，那么客户端将从证书中提取出服务端的公钥，用于加密后续协商密钥时发送给服务端的信息
// 如果服务端要求验证客户端（接到服务端发送的CertificateRequest请求），那么客户端还需通过ClientCertificate将自己的公钥证书发送给服务端进行验证
// 收到服务端对称加密共享密钥协商的请求后，客户端根据之前的随机数、确定的加密算法组合以及服务端发来的参数计算出最终的会话密钥，然后将服务端单独计算出会话密钥所需的信息用服务端的公钥加密后以ClientKeyExchange请求发送给服务端
// 随后客户端用ClientChangeCipher通知服务端从现在开始发送的消息都是加密的，最后有一个Finished消息来验证商法的对称加密共享密钥协商是否成功。验证方法就是通过协商好的新共享密钥和对称加密算法进行一段特定内容的加密，并以服务端是否能够正确解密该请求报文作为密钥协商成功与否的判定标准
// 而被加密的这段特定内容包含的是连接至今的全部报文内容。Finished报文作为该轮通信的结束标志，也是客户端发出的第一条使用协商密钥加密的信息

// 4、ServreChangeCipher & Finished（服务端 -> 客户端）
// 服务端收到客户端发过来的ClientKeyExchange中的参数后，也将单独计算出会话密钥。之后和客户端一样，服务端用ServerChangeCipher通知客户端从现在开始发送的消息都是加密的
// 最后服务端用一个Finished消息跟在ServerChangeCipher后面，既用于标识该轮握手结束，也用于验证对方计算出来的共享密钥是否有效。这也是服务端发出的第一条使用协商密钥加密的信息。

// 一旦HTTPS安全传输层的连接成功建立起来，后续双方通信的内容（应用层的HTTP协议）就会在一个进过加密处理的安全通道中得以传输

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!\n")
	})

	fmt.Println("server start at https://localhost:8081")
	// ListenAndServeTLS(addr, certFile, keyFile, handler)
	err := http.ListenAndServeTLS("localhost:8081", "server.crt", "server.key", nil)
	if err != nil {
		fmt.Printf("server start failed: %v\n", err)
	}
}
