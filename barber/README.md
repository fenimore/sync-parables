# Sleeping Barber

> run `go test -v` for example scenarios

## Problem

A **synchronization** problem for multiple threads -- a Barber checks on the waiting room and then cuts hair or goes to sleep (if there is no one in the waiting room). Concurrently, the **Customer** checks on the Barber; if the Barber is sleeping, the Barber wakes up, is  The Barber _thread_ goes from checking on customers to either cutting hair or going to sleep (if there is no one in the waiting room.

```
type Barber struct {
    name     string
    sync.Mutex         // for controlling access to state
    state    int       // sleeping/checking/cutting
    customer *Customer
}
```

From [Wikipedia](https://en.wikipedia.org/wiki/Sleeping_barber_problem),

> The problems are all related to the fact that the actions by both the barber and the customer (checking the waiting room, entering the shop, taking a waiting room chair, etc.) all take an unknown amount of time. For example, a customer may arrive and observe that the barber is cutting hair, so he goes to the waiting room. While he is on his way, the barber finishes the haircut he is doing and goes to check the waiting room. Since there is no one there (the customer not having arrived yet), he goes back to his chair and sleeps. The barber is now waiting for a customer and the customer is waiting for the barber.

The problem boils down to synchronizing change of state of the barber (sleeping, checking, cutting) and of the customer (checking, waiting, cutting).


## Solution

The solution relies on a **Mutex** lock for assuring only one thread can change state at a time -- so that way the barber is never checking for the customers when a customer is checking for the barber (which would cause a **deadlock**). As long as the barber _or_ the customer is checking, the mutex should block the other from doing so.

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
            b.Unlock()
            HairCut(c, b)
        default:               // if waiting room is empty
            b.Unlock()
            b.Lock()           // go to sleep on chair Zzzz
            b.state = sleeping
            b.customer = nil
            b.Unlock()
            c := <-wakers      // block, wait for waker to arrive
            HairCut(c, b)
        }
    }
}
```

## Example output

The `Customer` is printed as the last four characters of their _memory address_ (pointer).

```
=== RUN   TestChecking
Checking waiting room: 0
Sleeping Barber ZzzzZzz - <nil>
Customer 4020 checks sleeping barber room: 0, w 0 - customer: <nil>
Sleeping barber 4020, room: 0, wake: 0
Woken by 4020
Cutting  4020's hair
Customer a000 checks  cutting barber room: 0, w 0 - customer: 4020
Checking waiting room: 1
Cutting  a000's hair
Checking waiting room: 0
Sleeping Barber ZzzzZzz - <nil>
Customer a0a0 checks sleeping barber room: 0, w 0 - customer: <nil>
Sleeping barber a0a0, room: 0, wake: 0
Woken by a0a0
Customer 0f40 checks sleeping barber room: 0, w 0 - customer: <nil>
Sleeping barber 0f40, room: 0, wake: 0
Customer 4030 checks sleeping barber room: 0, w 1 - customer: <nil>
Sleeping barber 4030, room: 0, wake: 1
Cutting  a0a0's hair
Checking waiting room: 1
Cutting  4030's hair
Customer 4040 checks  cutting barber room: 0, w 1 - customer: 4030
Checking waiting room: 1
Cutting  4040's hair
Customer 10c0 checks  cutting barber room: 0, w 1 - customer: 4040
Checking waiting room: 1
Cutting  10c0's hair
Customer 40b0 checks  cutting barber room: 0, w 1 - customer: 10c0
Checking waiting room: 1
Cutting  40b0's hair
Checking waiting room: 0
Sleeping Barber ZzzzZzz - <nil>
Woken by 0f40
Cutting  0f40's hair
Checking waiting room: 0
No more customers today
```
