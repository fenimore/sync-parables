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

type Table struct {
	sync.Mutex
	paper  chan int
	grass  chan int
	match  chan int
	signal chan int
}

func arbitrate(t *Table) {
	for {
		next := rand.Intn(3)
		fmt.Printf("Table chooses %d\n", next)
		t.signal <- next
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
	}
}

func smoker(t *Table, name string, smokes int) {
	var chosen = -1
	has := map[int]bool{
		paper: paper == smokes,
		grass: grass == smokes,
		match: match == smokes,
	}
	for {
		has[smokes] = true // smokes -> infinite smoke
		select {
		case signal := <-t.signal:
			if chosen == signal {
				t.signal <- signal
				continue
			} else {
				chosen = signal
			}
		case item := <-t.paper:
			if smokes == item || chosen != smokes {
				t.paper <- item
				continue
			}
			// consume supply
			has[item] = true
		case item := <-t.grass:
			if smokes == item || chosen != smokes {
				t.paper <- item
				continue
			}
			// consume supply
			has[item] = true
		case item := <-t.match:
			if smokes == item || chosen != smokes {
				t.paper <- item
				continue
			}
			// consume supply
			has[item] = true
		}

		if has[grass] && has[paper] && has[match] {
			time.Sleep(time.Millisecond * 10)
			// Finish consuming
			has[paper], has[match], has[grass] = false, false, false
			has[smokes] = true // infinite supply ;)
		}
	}
}

func main() {
	fmt.Println("What")
}
