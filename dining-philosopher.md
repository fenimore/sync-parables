# Dining Philosophers

## Problem

N philosophers begin eating, with N chopsticks, one between each pair of philosophers. The Philosophers pick up one chopstick, another, eat, and then put the chopsticks down and think. Eventually, each philosopher picks up the chopstick to their left and **deadlocks**, waiting for the right chopstick. Unable to eat, they starve.


```
Butler:, Left: e2c0 Right: e2c8 // Representing Chopsticks
bhooks:, Left: e2c8 Right: e2d0 // by their reference in memory
Simone:, Left: e2d0 Right: e2d8
Bingen:, Left: e2d8 Right: e2e0
Arendt:, Left: e2e0 Right: e2c0
```

Chopsticks are `sync.Mutex` **locks**, Mutex locks are synchronization tools for concurrent programs which limit access to a resource for multliple threads. This script represents the philosophers:


```
type Philosopher struct {
    name  string
    left  *sync.Mutex // a chopstick
    right *sync.Mutex // a chopstick
    food  int
}

func (p *Philosopher) eat() {
    p.left.Lock()                      // Pick up chopsticks
    p.right.Lock()
    p.food += 1
    time.Sleep(time.Millisecond * 100) // Eat food
    p.left.Unlock()                    // Put down chopsticks
    p.right.Unlock()
}
```

> fatal error: all goroutines are asleep - deadlock!


## Solution

**Resource Hierarchy** assigns a hierarchy, or _partial order_, to the philosophers' chopstick preference: each will reach for the left most chopstick first, _except_ for the last philosopher, who will reach for the right most first. Partial order ranks the chopsticks and only one philosopher will have access to the _highest_ fork._

```
        // left index is i and right is i+1
        // swapping the values for the last philosophers
        // and only one philosopher will have access to the
        // "highest" chopstick
        if i == n-1 {
            philosphers[i].left = chopsticks[left]
            philosphers[i].right = chopsticks[right]
        } else { // partial order
            philosphers[i].right = chopsticks[left]
            philosphers[i].left = chopsticks[right]
        }
```

Output:

```

    Bingen picks up left: 0xc42000e2e0
    Bingen picks up right: 0xc42000e2d8
    Bingen eats 200/200
Bingen is finished
Philosphers are full

```


## Example Failure:

Output:

```
start:  0 Butler:, Left: 0xc42000e2c0 Right: 0xc42000e2c8
start:  1 bhooks:, Left: 0xc42000e2c8 Right: 0xc42000e2d0
start:  2 Simone:, Left: 0xc42000e2d0 Right: 0xc42000e2d8
start:  3 Bingen:, Left: 0xc42000e2d8 Right: 0xc42000e2e0
start:  4 Arendt:, Left: 0xc42000e2e0 Right: 0xc42000e2c0
    Bingen picks up left: 0xc42000e2d8
    Butler picks up left: 0xc42000e2c0
    Bingen picks up right: 0xc42000e2e0
    Butler picks up right: 0xc42000e2c8
    Simone picks up left: 0xc42000e2d0
    Bingen eats 1/200
    Butler eats 1/200
    Bingen picks up left: 0xc42000e2d8
    Bingen picks up right: 0xc42000e2e0
    Butler picks up left: 0xc42000e2c0
    Butler picks up right: 0xc42000e2c8
    Butler eats 2/200
    Bingen eats 2/200
    Butler picks up left: 0xc42000e2c0
    Simone picks up right: 0xc42000e2d8
    Butler picks up right: 0xc42000e2c8
    Arendt picks up left: 0xc42000e2e0
    Butler eats 3/200
    Arendt picks up right: 0xc42000e2c0
    bhooks picks up left: 0xc42000e2c8
    bhooks picks up right: 0xc42000e2d0
    Simone eats 1/200
    Bingen picks up left: 0xc42000e2d8
    bhooks eats 1/200
    bhooks picks up left: 0xc42000e2c8
    Arendt eats 1/200
    Bingen picks up right: 0xc42000e2e0
    Butler picks up left: 0xc42000e2c0
    Simone picks up left: 0xc42000e2d0
    Bingen eats 3/200
    Bingen picks up left: 0xc42000e2d8
    Bingen picks up right: 0xc42000e2e0
    Bingen eats 4/200
    Simone picks up right: 0xc42000e2d8
    Arendt picks up left: 0xc42000e2e0
    Simone eats 2/200
    bhooks picks up right: 0xc42000e2d0
    Bingen picks up left: 0xc42000e2d8
    bhooks eats 2/200
    bhooks picks up left: 0xc42000e2c8
    bhooks picks up right: 0xc42000e2d0
    bhooks eats 3/200
    bhooks picks up left: 0xc42000e2c8
    bhooks picks up right: 0xc42000e2d0
    bhooks eats 4/200
    bhooks picks up left: 0xc42000e2c8
    bhooks picks up right: 0xc42000e2d0
    bhooks eats 5/200
    bhooks picks up left: 0xc42000e2c8
    bhooks picks up right: 0xc42000e2d0
    bhooks eats 6/200
    bhooks picks up left: 0xc42000e2c8
    bhooks picks up right: 0xc42000e2d0
    bhooks eats 7/200
    Simone picks up left: 0xc42000e2d0
    bhooks picks up left: 0xc42000e2c8
fatal error: all goroutines are asleep - deadlock!
```
