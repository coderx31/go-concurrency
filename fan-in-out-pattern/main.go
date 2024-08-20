package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

type City struct {
	Name       string
	Population int
}

func main() {
	f, err := os.Open("./fan-in-out-pattern/worldcities.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			log.Println("failed to close file", err)
		}
	}()

	rows := genRows(f)

	// fan out pattern
	// more than one worker competing to consume the same channel

	//       __ worker_1
	// rows /
	//      \__ worker_2
	//       __ worker_n...
	ur1 := upperCityName(filterOutMinPopulation(rows))
	ur2 := upperCityName(filterOutMinPopulation(rows))
	ur3 := upperCityName(filterOutMinPopulation(rows))

	// fan in pattern consolidates the output from multiple channel into one
	//
	// worker_1 ___
	// worker_2 ___\ output
	// worker_3 ___/
	for c := range fanIn(ur1, ur2, ur3) {
		fmt.Println(c)
	}
}

func genRows(r io.Reader) chan City {
	out := make(chan City)
	go func() {
		reader := csv.NewReader(r)
		_, err := reader.Read()
		if err != nil {
			log.Fatal(err)
		}
		for {
			row, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}

			populationInt, err := strconv.Atoi(row[9])
			if err != nil {
				continue
			}
			out <- City{
				Name:       row[1],
				Population: populationInt,
			}
		}
		close(out)
	}()
	return out
}

func upperCityName(cities <-chan City) <-chan City {
	out := make(chan City)
	go func() {
		for c := range cities {
			out <- City{
				Name:       strings.ToUpper(c.Name),
				Population: c.Population,
			}
		}
		close(out)
	}()
	return out
}

func fanIn(channels ...<-chan City) <-chan City {
	out := make(chan City)
	wg := &sync.WaitGroup{}
	wg.Add(len(channels))
	for _, c := range channels {
		go func(city <-chan City) {
			for r := range city {
				out <- r
			}
			wg.Done()
		}(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func filterOutMinPopulation(cities <-chan City) <-chan City {
	out := make(chan City)
	go func() {
		for c := range cities {
			if c.Population > 40000 {
				out <- City{
					Name:       strings.ToUpper(c.Name),
					Population: c.Population,
				}
			}
		}
		close(out)
	}()
	return out
}
