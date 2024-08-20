package main

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/sync/semaphore"
	"log"
	"net/http"
	"runtime"
	"sync"
)

type Task struct {
	ID        int    `json:"id"`
	UserID    string `json:"user_id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func main() {
	wg := &sync.WaitGroup{}
	sem := semaphore.NewWeighted(10)
	//sem := make(chan bool, 10)
	for i := 1; i <= 100; i++ {
		// i:= i this also works
		fmt.Println(runtime.NumGoroutine())
		//sem <- true
		if err := sem.Acquire(context.Background(), 1); err != nil {
			log.Fatal(err)
		}
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			//defer func() {
			//	<-sem
			//}()
			defer sem.Release(1)
			t, err := getTask(j)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(fmt.Sprintf(`%d:%s`, j, t.Title))
		}(i)
	}
	wg.Wait()
}

func getTask(i int) (Task, error) {
	var t Task
	res, err := http.Get(fmt.Sprintf("https://jsonplaceholder.typicode.com/todos/%d", i))
	if err != nil {
		return Task{}, err
	}
	defer func() {
		if err = res.Body.Close(); err != nil {
			log.Println("failed to close the response body", err)
		}
	}()
	if err = json.NewDecoder(res.Body).Decode(&t); err != nil {
		return Task{}, err
	}
	return t, nil
}
