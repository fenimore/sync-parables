package barber

import (
	"sync"
	"testing"
	"time"
)

func TestChecking(t *testing.T) {
	b := NewBarber()
	b.name = "Sam"
	WaitingRoom := make(chan *Customer, 15) // 5 chairs
	Wakers := make(chan *Customer, 1)       // only one waker at a tim
	go barber(b, WaitingRoom, Wakers)
	time.Sleep(time.Millisecond * 100)
	wg = new(sync.WaitGroup)
	wg.Add(8)
	time.Sleep(time.Millisecond * 50)
	c := new(Customer)
	go customer(c, b, WaitingRoom, Wakers)
	time.Sleep(time.Millisecond * 10)
	c = new(Customer)
	go customer(c, b, WaitingRoom, Wakers)
	time.Sleep(time.Millisecond * 300)
	for i := 0; i < 3; i++ {
		time.Sleep(time.Millisecond * 10)
		c = new(Customer)
		go customer(c, b, WaitingRoom, Wakers)
	}
	for i := 0; i < 3; i++ {
		time.Sleep(time.Millisecond * 200)
		c = new(Customer)
		go customer(c, b, WaitingRoom, Wakers)
	}

	wg.Wait()
}

func TestTenCustomers(t *testing.T) {
	//lock = new(sync.Mutex)
	b := NewBarber()
	b.name = "Sam"
	WaitingRoom := make(chan *Customer, 15) // 5 chairs
	Wakers := make(chan *Customer, 1)       // only one waker at a tim
	go barber(b, WaitingRoom, Wakers)
	time.Sleep(time.Millisecond * 100)
	wg = new(sync.WaitGroup)
	n := 10
	wg.Add(10)
	for i := 0; i < n; i++ {
		time.Sleep(time.Millisecond * 50)
		c := new(Customer)
		go customer(c, b, WaitingRoom, Wakers)
	}

	time.Sleep(time.Millisecond * 50)
	c := new(Customer)
	go customer(c, b, WaitingRoom, Wakers)
	time.Sleep(time.Millisecond * 10)
	c = new(Customer)
	go customer(c, b, WaitingRoom, Wakers)
	time.Sleep(time.Millisecond * 300)
	for i := 0; i < 3; i++ {
		time.Sleep(time.Millisecond * 10)
		c = new(Customer)
		go customer(c, b, WaitingRoom, Wakers)
	}
	for i := 0; i < 3; i++ {
		time.Sleep(time.Millisecond * 200)
		c = new(Customer)
		go customer(c, b, WaitingRoom, Wakers)
	}
	wg.Wait()
}
