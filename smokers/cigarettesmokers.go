package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	paper = iota
	grass
	match
)

var smokeMap = map[int]string{
	paper: "paper",
	grass: "grass",
	match: "match",
}

var names = map[int]string{
	paper: "Sandy",
	grass: "Apple",
	match: "Daisy",
}

type Table struct {
	sync.Mutex
	paper  chan int
	grass  chan int
	match  chan int
	signal chan int
}

func arbitrate(t *Table, smokers [3]chan int) {
	for {
		time.Sleep(time.Millisecond * 500)
		next := rand.Intn(3)
		fmt.Printf("Table chooses %s: %s\n", smokeMap[next], names[next])
		switch next {
		case paper:
			t.grass <- 1
			t.match <- 1
		case grass:
			t.paper <- 1
			t.match <- 1
		case match:
			t.grass <- 1
			t.paper <- 1
		}
		for _, smoker := range smokers {
			smoker <- next
		}
		wg.Add(1)
		wg.Wait()
	}
}

func smoker(t *Table, name string, smokes int, signal chan int) {
	var chosen = -1
	has := map[int]bool{
		paper: paper == smokes,
		grass: grass == smokes,
		match: match == smokes,
	}
	for {
		has[smokes] = true // smokes -> infinite smoke
		select {
		case sign := <-signal:
			chosen = sign
		case item := <-t.paper:
			t.Lock()
			if chosen != smokes {
				t.paper <- item
			} else {
				// consume supply
				has[item] = true
			}
			t.Unlock()
			fmt.Println(name, "recveived ", has, smokeMap[item])
		case item := <-t.grass:
			t.Lock()
			if chosen != smokes {
				t.paper <- item
			} else {
				// consume supply
				has[item] = true
			}
			t.Unlock()
			fmt.Println(name, "recveived ", has, smokeMap[item])
		case item := <-t.match:
			t.Lock()
			if chosen != smokes {
				t.paper <- item
			} else {
				// consume supply
				has[item] = true
			}
			t.Unlock()
			fmt.Println(name, "recveived ", has, smokeMap[item], item)
		}

		if has[grass] && has[paper] && has[match] {
			fmt.Printf("%s is smoking, owner of %s\n", name, smokeMap[smokes])
			time.Sleep(time.Millisecond * 10)
			// Finish consuming
			has[paper], has[match], has[grass] = false, false, false
			has[smokes] = true // infinite supply ;)
			wg.Done()
		}
	}
}

const LIMIT = 10

var wg *sync.WaitGroup

func main() {
	wg = new(sync.WaitGroup)
	table := new(Table)
	table.match = make(chan int, LIMIT)
	table.paper = make(chan int, LIMIT)
	table.grass = make(chan int, LIMIT)
	var signals [3]chan int
	// three smokers
	for i := 0; i < 3; i++ {
		signal := make(chan int, 1)
		signals[i] = signal
		go smoker(table, names[i], i, signal)
	}

	arbitrate(table, signals)

}
