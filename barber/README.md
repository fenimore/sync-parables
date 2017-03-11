# Sleeping Barber

A **synchronization** problem for concurrent programs.

> run `go test -v` for example scenarios

## Problem

A barber checks on the waiting room and then either cuts hair or goes to sleep (if there is no one in the waiting room). Concurrently, the customer checks on the barber and if the barber is sleeping, the barber wakes up.

If the customer checks on the barber when the >barber is checking on an _empty_ waiting room, the barber would go back to sleep and the customer would go wait, possibly forever.

## Solution

The solution relies on a **Mutex** lock for assuring only one thread can change state at a time -- so that way the barber is never checking for the customers when a customer is checking for the barber (which would cause a **deadlock**). As long as the barber _or_ the customer is checking, the mutex should block the other from doing so.

```
type Barber struct {
    sync.Mutex         // for controlling access to state
    state    int       // sleeping/checking/cutting
    customer *Customer // customer currently being served
}
```

### Using channels to handle state

In addition to `sync.Mutex`, my solution handles the waiting room resource by passing customers into **channels**, `chan *Customer`, like queues with built in mutexes. Channels are safe ways for passing messages between concurrent threads. The customer enters and checks the barber's state, and then passes (itself) into a _buffered_ channel, the waiting room. The customer `switch`es on the barber's state, and then `select`s on a channel: waking the barber up `make(chan *Customer, 1)` or going to the waiting room `make(chan *Customer, 10)`. If the channels are full, the customer leaves.

```

func customer(c *Customer, b *Barber, wr chan<- *Customer, wakers chan<- *Customer) {
    b.Lock()                   // Arrive and Check on Barber
    switch b.state {
    case sleeping:
        select {
        case wakers <- c:      // Go wake up barber if asleep
        default:               // if there is someone already
            select {           // on their way to "waking" the Barber
            case wr <- c:      // go to waiting roomn
            default:           // if full, leave
            }
        }
    case cutting:
        select {
        case wr <- c:         // Go to waiting room if Barber is cutting
        default:              // if full waiting, leave shop
        }
    case checking:            // BOTH goroutines checking at once could result in deadlock
        panic("Customer shouldn't check for the Barber when the barber is checking the waiting room")
    }
    b.Unlock()
}
```

The barber thread, when sleeping, blocks on the `wakers` channel. The barber _gets_ to sleep by having an empty waiting room, when `wr` isn't sending a `*Customer`, the `default` case is selected.

```
func barber(b *Barber, wr chan *Customer, wakers chan *Customer) {
    for {
        b.Lock()
        defer b.Unlock()
        b.state = checking     // barber goes to check the waiting room
        b.customer = nil       // current served customer
        time.Sleep(time.Millisecond * 100)
        select {
        case c := <-wr:        // cuts hair of first person in queue
            HairCut(c, b)      // unlocks during cut
            b.Unlock()         // barber is cutting
        default:               // if waiting room is empty
            b.state = sleeping
            b.customer = nil
            b.Unlock()         // go to sleep on chair Zzzz
            c := <-wakers      // block, wait for waker to arrive
            b.Lock()
            HairCut(c, b)
            b.Unlock()
        }
    }
}
```

## Terminating the program

This could go on forever, but instead one can `add()` to a `sync.WaitGroup` struct for every customer, and `wg.Done()` after the customer leaves.

## Example output

The `Customer` is printed as the last four characters of their _memory address_ (pointer).

```
=== RUN   TestChecking
Checking waiting room: 0
Sleeping Barber ZzzzZzz - <nil>
Customer 12140 checks sleeping barber | room: 0, w 0 - customer: <nil>
Woken by 12140
Cutting  12140's hair
Customer 12150 checks  cutting barber | room: 0, w 0 - customer: 12140
Checking waiting room: 1
Cutting  12150's hair
Checking waiting room: 0
Sleeping Barber ZzzzZzz - <nil>
Customer 121d0 checks sleeping barber | room: 0, w 0 - customer: <nil>
Customer 74e30 checks sleeping barber | room: 0, w 0 - customer: <nil>
Customer 74e40 checks sleeping barber | room: 0, w 1 - customer: <nil>
Woken by 121d0
Cutting  121d0's hair
Checking waiting room: 1
Cutting  74e40's hair
Customer 74e50 checks  cutting barber | room: 0, w 1 - customer: 74e40
Checking waiting room: 1
Cutting  74e50's hair
Customer c8150 checks  cutting barber | room: 0, w 1 - customer: 74e50
Checking waiting room: 1
Cutting  c8150's hair
Customer 12270 checks  cutting barber | room: 0, w 1 - customer: c8150
Checking waiting room: 1
Cutting  12270's hair
Checking waiting room: 0
Sleeping Barber ZzzzZzz - <nil>
Woken by 74e30
Cutting  74e30's hair
Checking waiting room: 0

No more customers for the day
```
