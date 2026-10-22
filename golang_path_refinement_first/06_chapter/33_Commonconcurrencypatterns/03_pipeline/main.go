package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"sync"
	"time"
)

func newNumGenerator(start, count int) <-chan int {
	c := make(chan int)
	go func() {
		for i := start; i < start+count; i++ {
			c <- i
		}
		close(c)
	}()
	return c
}

func filterOdd(in int) (int, bool) {
	if in%2 != 0 {
		return 0, false
	}
	return in, true
}

func square(in int) (int, bool) {
	return in * in, true
}

func spawn(f func(int) (int, bool), in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for v := range in {
			r, ok := f(v)
			if ok {
				out <- r
			}
		}
		close(out)
	}()
	return out
}

// 管道模型的两种扩展模式
// 1、扇出模式：多个功能相同的goroutine从同一个channel读取数据并处理，直到该channel关闭，这种情况被称为“扇出”。使用扇出模式可以在一组goroutine中均衡分配工作量，从而更均衡的利用CPU
// 2、扇入模式：处理程序面对不止一个输入channel。我们把所有输入channel的数据汇聚到一个统一的输入channel，然后处理程序在从这个channel中读取数据并处理，直到该channel因所有输入channel关闭而关闭。这种情况被称为扇入
func spawnGroup(name string, num int, f func(int) (int, bool), in <-chan int) <-chan int {
	groupOut := make(chan int)
	var outSlice []chan int
	for i := 0; i < num; i++ {
		out := make(chan int)
		go func(i int) {
			name := fmt.Sprintf("%s-%d:", name, i)
			fmt.Printf("%s begin to work...\n", name)
			// 多个goroutine共同 range 同一个in通道！
			for v := range in {
				r, ok := f(v)
				if ok {
					out <- r
				}
			}
			// in通道关闭、数据读完之后，关闭自己这个worker的输出通道
			close(out)
			fmt.Printf("%s work done\n", name)
		}(i)
		outSlice = append(outSlice, out)
	}

	go func() {
		var wg sync.WaitGroup
		for _, out := range outSlice {
			wg.Add(1)
			go func(out <-chan int) {
				for v := range out {
					groupOut <- v
				}
				wg.Done()
			}(out)
		}
		wg.Wait()
		close(groupOut)
	}()
	return groupOut
}

// 超时与取消模式最佳实践 使用context库
type result struct {
	value string
}

func first(servers ...*httptest.Server) (result, error) {
	c := make(chan result)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	queryFunc := func(i int, server *httptest.Server) {
		url := server.URL
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Printf("query goroutine-%d:http NewRequest error: %s\n", i, err)
			return
		}
		req = req.WithContext(ctx)

		log.Printf("query goroutine-%d:send request...\n", i)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("query goroutine-%d:get return error:%s\n", i, err)
			return
		}
		log.Printf("query goroutine-%d:get response\n", i)
		defer resp.Body.Close()
		body, _ := io.ReadAll(req.Body)
		c <- result{
			value: string(body),
		}
		return
	}
	for i, serv := range servers {
		go queryFunc(i, serv)
	}
	select {
	case r := <-c:
		return r, nil
	case <-time.After(500 * time.Millisecond):
		return result{}, errors.New("timeout")
	}
}

func fakeWeatherServer(name string, interval int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s receiver a http request\n", name)
		time.Sleep(time.Duration(interval) * time.Millisecond)
		w.Write([]byte(name + ":ok"))
	}))
}

func main() {
	in := newNumGenerator(1, 20)
	out := spawn(square, spawn(filterOdd, in))
	for v := range out {
		println(v)
	}

	in2 := newNumGenerator(1, 20)
	out2 := spawnGroup("square", 2, square, spawnGroup("filterOdd", 3, filterOdd, in2))
	time.Sleep(3 * time.Second)

	for v := range out2 {
		fmt.Println(v)
	}

	result, err := first(fakeWeatherServer("open-weather-1", 200), fakeWeatherServer("open-weather-2", 1000),
		fakeWeatherServer("open-weather-3", 600))
	if err != nil {
		log.Println("invoke first error:", err)
		return
	}
	fmt.Println(result)
	time.Sleep(10 * time.Second)
}
