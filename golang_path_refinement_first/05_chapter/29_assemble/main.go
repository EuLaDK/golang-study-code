package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"strings"
)

// 接口为连接点水平组合的集中形式
// 1、基本形式 func YourFuncName(param YourInterfaceType) // 通过类型连接
// 在标准库里面的应用
// $GOROOT/src/io/ioutil/ioutil.go
// func readAll(r io.Reader, capacity int64) (b []byte, err error)
// $GOROOT/src/io/io.go
// func Copy(dst Writer, src Reader) (written int64, err error)

// 2、包裹函数 func YourWrapperFunc(param YourInterfaceType) YourInterfaceType 它接受接口类型参数，并返回与其参数类型相同的返回值
// Go标准库中一个典型的包裹函数io.LimitReader：
// // $GOROOT/src/io/io.go
// func LimitReader(r Reader, n int64) Reader { return &LimitedReader{r, n} }
// type LimitedReader struct {
// R Reader
// N int64
// }
// func (l *LimitedReader) Read(p []byte) (n int, err error) {
// ...
// }

func CapReader(r io.Reader) io.Reader {
	return &capitalizedReader{r: r}
}

type capitalizedReader struct {
	r io.Reader
}

func (r *capitalizedReader) Read(p []byte) (int, error) {
	n, err := r.r.Read(p)
	if err != nil {
		return 0, err
	}

	q := bytes.ToUpper(p)
	for i, v := range q {
		p[i] = v
	}
	return n, err
}

// 适配器函数类型
// 适配器函数类型（adapter function type）是一个辅助水平组合实
// 现的“工具”类型。强调一下，它是一个类型。它可以将一个满足特定函
// 数签名的普通函数显式转换成自身类型的实例，转换后的实例同时也是
// 某个单方法接口类型的实现者
// type Handler interface {
// ServeHTTP(ResponseWriter, *Request)
// }
// type HandlerFunc func(ResponseWriter, *Request)
// func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
// f(w, r)
// }
// // chapter5/sources/horizontal-composition-3.go
// func greetings(w http.ResponseWriter, r *http.Request) {
// fmt.Fprintf(w, "Welcome!")
// }

// 4中间件
// “中间件”（middleware）这个词的含义可大可小，在Go Web编程
// 中，它常常指的是一个实现了http.Handler接口的http.HandlerFunc类型
// 实例。实质上，这里的中间件就是包裹函数和适配器函数类型结合的
// 产物

// func validateAuth(s string) error {
// 	if s != "123456" {
// 		return fmt.Errorf("%s", "bad auth token")
// 	}
// 	return nil
// }
// func greetings(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Welcome!")
// }
// func logHandler(h http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		t := time.Now()
// 		log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t)
// 		h.ServeHTTP(w, r)
// 	})
// }
// func authHandler(h http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		err := validateAuth(r.Header.Get("auth"))
// 		if err != nil {
// 			http.Error(w, "bad auth param", http.StatusUnauthorized)
// 			return
// 		}
// 		h.ServeHTTP(w, r)
// 	})
// }

// 我们看到所谓中间件（如logHandler、authHandler）本质上就是一
// 个包裹函数（支持链式调用），但其内部利用了适配器函数类型
// （http.HandlerFunc）将一个普通函数（如例子中的几个匿名函数）转
// 换为实现了http.Handler的类型的实例，并将其作为返回值返回

func main() {
	r := strings.NewReader("hello,gopher!\n")
	r1 := CapReader(io.LimitReader(r, 4))
	if _, err := io.Copy(os.Stdout, r1); err != nil {
		log.Fatal(err)
	}
}
