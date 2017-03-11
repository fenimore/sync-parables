package main

import (
	"fmt"
	"sync"
	"time"
)

// Barber
const (
	sleeping = iota
	checking
	cutting
)

// Customer
const (
	going = iota
	_     // checking
	_     // cutting
	waiting
)

type Barber struct {
	name    string
	state   int
	wokenUp chan struct{}
}

type Customer struct {
	name  string
	state int
}

func NewBarber() (b *Barber) {
	return &Barber{
		name:  "Sam",
		state: sleeping,
	}
}

func NewCustomer() (c *Customer) {
	return &Customer{
		name:  "George",
		state: going,
	}
}

// Barber thread
func barber(b *Barber, wr chan *Customer, wakers chan *Customer) {
	for {
		b.state = checking
		// checking the waiting room
		fmt.Printf("Checking, %d, for customer room: %d\n", b.state, len(wr))
		time.Sleep(time.Millisecond * 10)
		select {
		case c := <-wr:
			HairCut(c, b)
		default:
			fmt.Printf("Sleeping Barber\n")
			b.state = sleeping
			c := <-wakers
			fmt.Printf("Woken by %p\n", c)
			HairCut(c, b)
		}
	}
}

// customer goroutine
// just fizzles out if it's full, otherwise the customer
// is passed along to the channel handling it's haircut etc
func customer(c *Customer, b *Barber, wr chan<- *Customer, wakers chan<- *Customer) {
	// arrive
	time.Sleep(time.Millisecond * 150)
	// Check on barber
	defer fmt.Printf("Customer %p comes in to: %d, room: %d\n", c, b.state, len(wr))
	switch b.state {
	case sleeping:
		select {
		case wakers <- c:
		default:
			wr <- c // wakers is full
		}
		return
	case cutting, checking:
		select {
		case wr <- c:
		default:
			// full, leave shop
			wg.Done()
		}
	}
}

func HairCut(c *Customer, b *Barber) {
	b.state = cutting
	c.state = cutting
	// cut some hair
	fmt.Printf("Cutting  %p's hair\n", c)
	time.Sleep(time.Millisecond * 50)
	wg.Done()
}

var wg *sync.WaitGroup // Amount of potentional customers

func main() {
	//lock = new(sync.Mutex)
	b := NewBarber()
	b.name = "Sam"
	WaitingRoom := make(chan *Customer, 5) // 5 chairs
	Wakers := make(chan *Customer, 1)      // only one waker at a time
	go func() {
		barber(b, WaitingRoom, Wakers)
	}()
	wg = new(sync.WaitGroup)
	n := 10
	wg.Add(10)
	// Spawn customers
	for i := 0; i < n; i++ {
		time.Sleep(time.Millisecond * 20)
		c := NewCustomer()
		go customer(c, b, WaitingRoom, Wakers)
	}

	wg.Wait()
}
