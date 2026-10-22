package main

import (
	"log"
	"sync"
	"time"
)

// 带缓冲channel T为channel中元素的类型，capacity为带缓冲channel的缓冲区容量
// c := make(chan T, capacity)
// 对一个带缓冲channel，在缓冲区无数据或有数据但未满的情况下，对其进行发送操作的goroutine不会阻塞；
// 在缓冲区已满的情况下，对其发送操作的goroutine会阻塞；
// 在缓冲区为空的情况下，对其进行接收操作的goroutine亦会阻塞

// 带缓冲channel收发性能好于无缓冲区channel

// 带缓冲channel用做计数信号量
var active = make(chan struct{}, 3)
var jobs = make(chan int, 10)

// len(channel)的应用
// 如果s是字符串（string）类型，len(s)返回字符串中的字
// 节数；
// 如何s是[n]T或*[n]T的数组类型，len(s)返回数组的长度
// n；
// 如果s是[]T的切片（slice）类型，len(s)返回切片的当前长
// 度；
// 如果s是map[K]T的map类型，len(s)返回map中已定义的
// key的个数；
// 如果s是chan T类型，那么len(s)针对channel的类型不同，
// 有如下两种语义：
// 当s为无缓冲channel时，len(s)总是返回0；
// 当s为带缓冲channel时，len(s)返回当前channel s中尚
// 未被读取的元素个数。

func main() {
	go func() {
		for i := 0; i < 8; i++ {
			jobs <- (i + 1)
		}
		close(jobs)
	}()
	var wg sync.WaitGroup

	for j := range jobs {
		wg.Add(1)
		go func(j int) {
			active <- struct{}{}
			log.Printf("handle job: %d\n", j)
			time.Sleep(2 * time.Second)
			<-active
			wg.Done()
		}(j)
	}
	wg.Wait()
}
