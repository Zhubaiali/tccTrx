package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(3)

	go printFunc(&wg, "cat")
	go printFunc(&wg, "fish")
	go printFunc(&wg, "dog")

	wg.Wait()
}

func printFunc(wg *sync.WaitGroup, s string) {
	defer wg.Done()
	for i := 0; i < 100; i++ {
		fmt.Println(s)
	}
}
