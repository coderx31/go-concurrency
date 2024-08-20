package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	started := time.Now()
	foods := []string{"mushroom pizza", "pasta", "kebab", "cake"}
	wg := &sync.WaitGroup{}
	for _, food := range foods {
		wg.Add(1)
		go cook(food, wg)
	}
	wg.Wait()
	fmt.Printf("done in %v \n", time.Since(started))
}

func cook(food string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("cooking %s...\n", food)
	time.Sleep(2 * time.Second)
	fmt.Printf("done cooking %s \n", food)
	fmt.Println("")
}
