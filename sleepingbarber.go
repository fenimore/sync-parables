package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	sleeping = iota
	checking
	cutting
)

var stateLog = map[int]string{
	0: "sleeping",
	1: "checking",
	2: " cutting",
}
var wg *sync.WaitGroup // Amount of potentional customers

type Barber struct {
	sync.Mutex
	name     string
	state    int
	customer *Customer
}

type Customer struct {
	name string
}

func (c *Customer) String() string {
	return fmt.Sprintf("%p", c)[8:]
}

func NewBarber() (b *Barber) {
	return &Barber{
		name:  "Sam",
		state: sleeping,
	}
}

// Barber goroutine
// Checks for customers
// Sleeps - wait for wakers to wake him up
func barber(b *Barber, wr chan *Customer, wakers chan *Customer) {
	for {
		b.Lock()
		b.state = checking
		b.customer = nil
		b.Unlock()
		// checking the waiting room
		fmt.Printf("Checking waiting room: %d\n", len(wr))
		time.Sleep(time.Millisecond * 100)
		select {
		case c := <-wr:
			HairCut(c, b)
		default: // Waiting room is empty
			fmt.Printf("Sleeping Barber ZzzzZzz - %s\n", b.customer)
			b.Lock()
			b.state = sleeping
			b.customer = nil
			b.Unlock()
			c := <-wakers
			fmt.Printf("Woken by %s\n", c)
			HairCut(c, b)
		}
	}
}

func HairCut(c *Customer, b *Barber) {
	b.Lock()
	b.state = cutting
	b.customer = c
	b.Unlock()
	// cut some hair
	fmt.Printf("Cutting  %s's hair\n", c)
	time.Sleep(time.Millisecond * 100)
	b.Lock()
	b.state = cutting
	b.customer = nil
	b.Unlock()
	wg.Done()
}

// customer goroutine
// just fizzles out if it's full, otherwise the customer
// is passed along to the channel handling it's haircut etc
func customer(c *Customer, b *Barber, wr chan<- *Customer, wakers chan<- *Customer) {
	// arrive
	fmt.Printf("Customer %s comes in to: %s barber, room: %d, wake: %d - customer: %s\n",
		c, stateLog[b.state], len(wr), len(wakers), b.customer)
	time.Sleep(time.Millisecond * 50)
	// Check on barber
	b.Lock()
	defer b.Unlock()
	switch b.state {
	case sleeping:
		fmt.Printf("Sleeping barber %p, room: %d, wake: %d\n", c, len(wr), len(wakers))
		select {
		case wakers <- c:
		default:
			select {
			case wr <- c:
			default:
				wg.Done()
			}
		}
	// TODO: fallthrough?
	case cutting:
		select {
		case wr <- c:
		default:
			// full, leave shop
			wg.Done()
		}
	case checking:
		select {
		case wr <- c:
		default:
			// full, leave shop
			wg.Done()
		}

	}
}

func main() {
	//lock = new(sync.Mutex)
	b := NewBarber()
	b.name = "Sam"
	WaitingRoom := make(chan *Customer, 15) // 5 chairs
	Wakers := make(chan *Customer, 1)       // only one waker at a time
	go func() {
		barber(b, WaitingRoom, Wakers)
	}()
	time.Sleep(time.Millisecond * 100)
	wg = new(sync.WaitGroup)
	n := 10
	wg.Add(10)
	// Spawn customers
	for i := 0; i < n; i++ {
		time.Sleep(time.Millisecond * 50)
		c := new(Customer)
		go customer(c, b, WaitingRoom, Wakers)
	}

	wg.Wait()
	fmt.Println("No more customers for the day")
}
