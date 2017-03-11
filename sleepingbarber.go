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

type Barber struct {
	sync.Mutex
	name  string
	state int
}

type Customer struct {
	name string
}

func NewBarber() (b *Barber) {
	return &Barber{
		name:  "Sam",
		state: sleeping,
	}
}

func NewCustomer() (c *Customer) {
	return &Customer{
		name: "George",
	}
}

// Barber thread
func barber(b *Barber, wr chan *Customer, wakers chan *Customer) {
	for {
		b.Lock()
		b.state = checking
		b.Unlock()
		// checking the waiting room
		fmt.Printf("Checking, %s, for customer room: %d\n", stateLog[b.state], len(wr))
		time.Sleep(time.Millisecond * 10)
		select {
		case c := <-wr:
			HairCut(c, b)
		default:
			fmt.Printf("Sleeping Barber\n")
			b.Lock()
			b.state = sleeping
			b.Unlock()
			c := <-wakers
			fmt.Printf("Woken by %p\n", c)
			HairCut(c, b)
		}
	}
}

func HairCut(c *Customer, b *Barber) {
	b.Lock()
	b.state = cutting
	b.Unlock()
	// cut some hair
	fmt.Printf("Cutting  %p's hair\n", c)
	time.Sleep(time.Millisecond * 100)
	wg.Done()
}

// customer goroutine
// just fizzles out if it's full, otherwise the customer
// is passed along to the channel handling it's haircut etc
func customer(c *Customer, b *Barber, wr chan<- *Customer, wakers chan<- *Customer) {
	// arrive
	fmt.Printf("Customer %p comes in to: %s barber, room: %d, wake: %d\n",
		c, stateLog[b.state], len(wr), len(wakers))
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
			fmt.Printf("Someone else is waking the barber, %p goes to waiting room\n", c)
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

var wg *sync.WaitGroup // Amount of potentional customers
var stateLog = map[int]string{
	0: "sleeping",
	1: "checking",
	2: " cutting",
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
	wg = new(sync.WaitGroup)
	n := 10
	wg.Add(10)
	// Spawn customers
	for i := 0; i < n; i++ {
		time.Sleep(time.Millisecond * 50)
		c := NewCustomer()
		go customer(c, b, WaitingRoom, Wakers)
	}

	wg.Wait()
}
