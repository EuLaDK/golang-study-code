package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

func subTwoTimeHasMonotonic() {
	t1 := time.Now()
	time.Sleep(time.Second * 5)
	t2 := time.Now()
	diff := t2.Sub(t1)
	fmt.Println(t2.After(t1))
	fmt.Printf("[hasMonotonic = 1] t2 - t1 = %v\n", diff)
}

func subTwoTimeNoMonotonic() {
	t1 := time.Date(2026, 10, 18, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2026, 10, 18, 12, 0, 0, 0, time.UTC)
	diff := t2.Sub(t1)
	fmt.Printf("[hasMonotonic = 0] t2 - t1 = %v\n", diff)
}

// Timer的创建
func create_timer_by_afterfunc() {
	_ = time.AfterFunc(1*time.Second, func() {
		fmt.Println("time created by afterfunc fired!")
	})
}

func create_timer_by_newtimer() {
	timer := time.NewTimer(2 * time.Second)
	select {
	case <-timer.C:
		fmt.Println("timer create by newtimer fired!")
	}
}

func create_timer_by_after() {
	select {
	case <-time.After(2 * time.Second):
		fmt.Println("time created by after fired!")
	}
}

func consume(c <-chan bool, timer *time.Timer) bool {
	if !timer.Stop() {
		select {
		case <-timer.C:
		default:
		}
	}
	timer.Reset(5 * time.Second)

	select {
	case b := <-c:
		if b == false {
			log.Printf("recv true, continue")
			return true
		}
		log.Printf("recv true, return")
		return false
	case <-timer.C:
		log.Printf("timer expired")
		return true
	}
}

func main() {
	// 获取当前时间
	timeNow := time.Now()
	fmt.Println(timeNow) // time.Now函数获取的当前时间中则既包含挂钟时间，也包含单调时间

	// 获取特定时区的当前时间 显示加载时区信息
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		fmt.Println("加载时间出错", err)
		return
	}

	t1 := timeNow.In(loc) // 转换成美国东部纽约时间表示
	fmt.Println(t1)

	// 时间的比较 直接用==和!=来做不合适 time包提供了专门的Equal方法
	// 当两个Time实例均带有单调时间数据时（hasMonotonic都为1），那么直接比较两者的单调时间是否相等，否则分别比较两个时间的整秒部分和非整秒部分。如果两个部分分别相等，那么两个时间相同，否则不同
	fmt.Println(timeNow.Equal(t1))
	// Time类型还提供了Before和After方法，用于判断两个即时时间的先后关系
	fmt.Println(timeNow.After(t1))
	fmt.Println(timeNow.Before(t1))
	// time包还可以用来对两个即时时间进行时间运算，其中最主要的运算就是由Sub方法提供的差值运算（Since和Until方法均是基于Sub方法实现的）
	subTwoTimeHasMonotonic()
	subTwoTimeNoMonotonic()

	// 格式化时间 go使用更简单的参考时间来格式化时间
	fmt.Println(time.Now().Format("2006年01月02日 15时04分05秒"))

	// 定时器的使用 time包提供了两类定时器：一次性定时器Timer和重复定时器Ticker。顾名思义，Timer只能进行一次定时器到期事件触发（fire），而Ticker则可以按一定时间间隔多次触发定时器到期事件
	create_timer_by_afterfunc()
	create_timer_by_newtimer()
	create_timer_by_after()
	// Timer的资源释放  在定时器被最小堆移除并触发事件后，其占用的内存资源、channel等都会在后续被垃圾收集器回收

	// 停止定时器 stop方法 和重新使用定时器 reset方法
	// Go建议只对如下两种定时器调用reset方法
	// 1-已经停止了的定时器（Stopped）
	// 2-已经触发过且Timer.C中的数据已经被读空
	// Go官方文档给出了推荐的使用模式
	// if !t.Stop(){
	// 	<-t.C
	// }
	// t.Reset(d)
	c := make(chan bool)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		for i := 0; i < 5; i++ {
			time.Sleep(time.Second * 7)
			c <- false
		}
		time.Sleep(time.Second * 7)
		c <- true
		wg.Done()
	}()

	go func() {
		timer := time.NewTimer(time.Second * 5)
		for {
			if b := consume(c, timer); !b {
				wg.Done()
				return
			}
		}
	}()
	wg.Wait()
}
