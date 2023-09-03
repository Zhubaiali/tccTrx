package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(3)

	// 创建一个互斥锁
	var mutex sync.Mutex

	go printFunc2(&wg, &mutex, "cat")
	go printFunc2(&wg, &mutex, "fish")
	go printFunc2(&wg, &mutex, "dog")

	wg.Wait()
}

func printFunc2(wg *sync.WaitGroup, mutex *sync.Mutex, s string) {
	defer wg.Done()
	for i := 0; i < 100; i++ {
		// 加锁
		mutex.Lock()
		fmt.Println(s)
		// 解锁
		mutex.Unlock()
	}
}
