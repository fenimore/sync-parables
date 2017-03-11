// Dining philosophers
// N philosophers sitting around around table.
// There is only N chopstick on the table, each
// between two philosphers.
// When a philosopher wants to eat, she must acquire her
// left and right chopstick
package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

// Philosopher eats food with a left and right chopstick.
type Philosopher struct {
	name  string
	id    int         // index, for logs
	left  *sync.Mutex // a chopstick
	right *sync.Mutex // a chopstick
	food  int
}

func (p *Philosopher) eat() {
	p.left.Lock()
	fmt.Printf("    %s picks up left: %p\n", p.name, p.left)
	p.right.Lock()
	fmt.Printf("    %s picks up right: %p\n", p.name, p.right)
	//fmt.Printf("%d|%s eats %d%% | %p %p\n", p.id, p.name, p.food/fullness, p.left, p.right)
	p.food += 1
	time.Sleep(time.Millisecond * 100) // eat
	p.left.Unlock()                    // Put down chopstick
	p.right.Unlock()
	fmt.Printf("    %s eats %d/%d\n", p.name, p.food, fullness)
}

func (p *Philosopher) String() string {
	return fmt.Sprintf("%s:, Left: %p Right: %p", p.name, p.left, p.right)
}

var wg *sync.WaitGroup
var fullness int = 100

func main() {
	var n = 5
	if len(os.Args) > 1 {
		in, err := strconv.Atoi(os.Args[1])
		if err != nil {
			return
		}
		n = in
	}
	wg = new(sync.WaitGroup)
	wg.Add(n)
	philosphers := make([]*Philosopher, n)
	chopsticks := make([]*sync.Mutex, n)

	for i := 0; i < n; i++ {
		chopsticks[i] = new(sync.Mutex)
		philosphers[i] = new(Philosopher)
		philosphers[i].id = i

	}

	philosphers[0].name = "Butler"
	philosphers[1].name = "bhooks"
	philosphers[2].name = "Simone"
	philosphers[3].name = "Bingen"
	philosphers[4].name = "Arendt"
	for i := 0; i < n; i++ {
		right := i + 1
		left := i
		// Last philosopher uses first philosopher's chopstick.
		if right == n {
			right = 0
		}
		if i == n-1 {
			philosphers[i].left = chopsticks[left]
			philosphers[i].right = chopsticks[right]
		} else { // Hierarchy Solution
			philosphers[i].right = chopsticks[left]
			philosphers[i].left = chopsticks[right]
		}

		fmt.Println("start: ", i, philosphers[i])
		go Dine(philosphers[i], i)
	}

	wg.Wait()
	fmt.Println("Philosphers are full")
}

func Dine(p *Philosopher, idx int) {
	time.Sleep(time.Millisecond * 10)
	for p.food != fullness {
		p.eat()
	}
	fmt.Printf("%s is finished\n", p.name)
	wg.Done()
}
