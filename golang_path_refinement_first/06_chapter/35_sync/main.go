package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// Go标准库sync.Mutex的定义
type Mutex struct {
	state int32  // 表示当前互斥锁的状态
	sema  uint32 // 用于控制锁状态的信号量
}

// 在使用sync包中类型时，推荐通过闭包方式或传递类型实例（或包裹该类型的类型实例）的地址或指针的方式进行
// 读写锁还是互斥锁
// 读写锁适合应用在具有一定并发量且读多写少的场合，在有大量并发读的情况下，多个goroutine可以同时持有读锁，从而减少在竞争中等待的时间；而互斥锁即便是读请求，同一时刻也只能有一个goroutine持有锁，其他goroutine只能阻塞在加锁操作上等待被调度

// 条件变量sync.Cond
type signal struct{}

var ready bool

func worker(i int) {
	fmt.Printf("worker %d: is working...\n", i)
	time.Sleep(1 * time.Second)
	fmt.Printf("worker %d: works done\n", i)
}

func spawnGroup(f func(i int), num int, groupSignal *sync.Cond) <-chan signal {
	c := make(chan signal)
	var wg sync.WaitGroup

	for i := 0; i < num; i++ {
		wg.Add(1)
		go func(i int) {
			groupSignal.L.Lock()
			for !ready {
				// Wait方法在goroutine挂起前会进行Unlock操作
				groupSignal.Wait()
			}
			groupSignal.L.Unlock()
			fmt.Printf("worker %d: start to work...\n", i)
			f(i)
			wg.Done()
		}(i + 1)
	}
	go func() {
		wg.Wait()
		c <- signal(struct{}{})
	}()
	return c
}

// sync.Once 单例模式
type Foo struct{}

var once sync.Once
var instance *Foo

func GetInstance(id int) *Foo {
	defer func() {
		if e := recover(); e != nil {
			log.Printf("goroutine-%d: caught a panic: %s", id, e)
		}
	}()
	log.Printf("goroutine-%d: enter GetInstance\n", id)
	once.Do(func() {
		instance = &Foo{}
		time.Sleep(3 * time.Second)
		log.Printf("goroutine-%d: the addr of instance is %p\n", id, instance)
		panic("panic in once.Do function")
	})
	return instance
}

// 使用sync.Pool 减轻垃圾回收压力
// sync.Pool的特点
// 1、它是goroutine并发安全的，可以被多个goroutine同时使用
// 2、放入该缓存池中的数据对象的生命是暂时的，随时都可能被垃圾回收掉
// 3、缓存池中的数据对象是可以重复利用的，这样可以在一定程度上降低数据对象重新分配的频度，减轻GC的压力
// 4、sync.Pool为每个P（goroutine调度模型中的P）单独建立一个local缓存池，进一步降低高并发下对锁的争抢

// sync.Pool处理长时间内存因为处理一些数据得不到释放的解决方案
// 1、限制要放回缓存池中的数据对象大小
// 2、建立多级缓存池

func main() {
	fmt.Println("start a group of workers...")
	groupSignal := sync.NewCond(&sync.Mutex{})
	c := spawnGroup(worker, 5, groupSignal)

	time.Sleep(5 * time.Second) // 模拟ready前的准备工作
	fmt.Println("the group of workers start to work...")

	groupSignal.L.Lock()
	ready = true
	groupSignal.Broadcast() // 唤醒所有等待的 goroutine
	groupSignal.L.Unlock()

	<-c
	fmt.Println("the group of workers work done!")

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			inst := GetInstance(i)
			log.Printf("goroutine-%d: the addr of instance returned is %p\n", i, inst)
			wg.Done()
		}(i + 1)
	}
	time.Sleep(5 * time.Second)
	inst := GetInstance(0)
	log.Printf("goroutine-0: the addr of instance returned is %p\n", inst)
	wg.Wait()
	log.Printf("all goroutines exit\n")
}
