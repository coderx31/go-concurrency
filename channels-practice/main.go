package main

import (
	"fmt"
	"time"
)

func main() {
	started := time.Now()
	foods := []string{"mushroom pizza", "pasta", "kebab", "cake"}
	results := make(chan bool)
	for _, food := range foods {
		go func(food string) {
			cook(food)
			results <- true
		}(food)
	}
	for i := 0; i < len(foods); i++ {
		<-results
	}
	fmt.Printf("done in %v \n", time.Since(started))
}

func cook(food string) {
	fmt.Printf("cooking %s...\n", food)
	time.Sleep(2 * time.Second)
	fmt.Printf("done cooking %s \n", food)
	fmt.Println("")
}
