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
	name   string
	state  int
	myTurn chan struct{}
}

type Checker interface {
	Check(chan *Customer)
}

func NewBarber() (b *Barber) {
	return &Barber{
		name:    "Sam",
		state:   sleeping,
		wokenUp: make(chan struct{}),
	}
}

func NewCustomer() (c *Customer) {
	return &Customer{
		name:   "George",
		state:  going,
		myTurn: make(chan struct{}),
	}
}

// Barber thread
func barber(b *Barber, wr chan *Customer, wakers chan *Customer) {
	for {
		// checking the waiting room
		fmt.Println("Checking for customer")
		time.Sleep(time.Millisecond * 10)
		select {
		case c := <-wr:
			HairCut(c, b)
		default:
			fmt.Printf("Sleeping Barber\n")
			b.state = sleeping
			c := <-wakers
			HairCut(c, b)
			b.state = checking
		}
	}
}

// customer goroutine
// just fizzles out if it's full, otherwise the customer
// is passed along to the channel handling it's haircut etc
func customer(c *Customer, b *Barber, wr chan<- *Customer, wakers chan<- *Customer) {
	// arrive
	time.Sleep(time.Millisecond * 100)
	// Check on barber
	fmt.Printf("Customer comes in to: %d\n", b.state)
	switch b.state {
	case sleeping:
		select {
		case wakers <- c:
		default:
			// wakers is full
		}
		return
	case cutting:
		select {
		case wr <- c:
		default:
			// full, leave shop
		}
	case checking:
		select {
		case wr <- c:
		default:
			// full, leave shop
		}
	}
	wg.Done()
}

// Customer Methods
func (c *Customer) Check(wr chan *Customer) bool {
	select {
	case wr <- c:
		return true
	default:
		//Waiting room is full, leave
		return false
	}

}

func HairCut(c *Customer, b *Barber) {
	b.state = cutting
	c.state = cutting
	// cut some hair
	fmt.Printf("Cutting %p's hair\n", c)
	time.Sleep(time.Millisecond * 50)
}

// var lock *sync.Mutex
var wg *sync.WaitGroup

func main() {
	//lock = new(sync.Mutex)
	b := NewBarber()
	b.name = "Sam"
	WaitingRoom := make(chan *Customer, 5) // 5 chairs
	Wakers := make(chan *Customer, 1)      // only one waker at a time?
	go func() {
		barber(b, WaitingRoom, Wakers)
	}()
	wg = new(sync.WaitGroup)
	n := 10
	wg.Add(10)
	// Spawn customers
	for i := 0; i < n; i++ {
		c := NewCustomer()
		go customer(c, b, WaitingRoom, Wakers)
	}

	wg.Wait()
}
