# Sleeping Barber

## Problem

A **synchronization** problem for multiple threads -- that of a **Barber** cutting hair and a stream of **Customers**. The Barber _thread_ goes from checking on customers to either cutting hair or going to sleep (if there is no one in the waiting room.

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

The problem boils down to synchronizing change of state of the barber (sleeping/checking/cutting) and of the customer (checking, waiting, cutting).


## Solution

Typically the solution relies on a **Mutex** lock for ensuring only one thread can change state at a time -- so that way the barber never checks fro customers when the customer is checking for the barber.

### Using channels to handle state

Rather than depending `sync.Mutex` _for the customer_, however, my solution handles customer state by passing customers into **channels**, `chan *Customer` -- that is, **Communicating Sequential Processess**, or CSP techniques. The customer enters and checks the barber's state, and then passes into a _buffered_ channel according to the state

```
func customer(c *Customer, b *Barber, wr chan<- *Customer, wakers chan<- *Customer) {
    b.Lock()                   // Arrive and Check on Barber
    defer b.Unlock()
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
    case cutting, checking:
        select {
        case wr <- c:         // Go to waiting room if Barber is cutting
        default:              // if full waiting, leave shop
        }
    }
}
```

The barber thread, when sleeping, blocks on the `wakers` channel, so when a customer sees the barber sleepers, they enter `wakers` -- `wakers` is a **buffered** channel of `1` length. `select` statements will fallthrough to `default` if the channel to enter, either the waker or waiting room `chan` is full.

```
func barber(b *Barber, wr chan *Customer, wakers chan *Customer) {
    for {
        b.Lock()
        b.state = checking     // barber goes to check the waiting room
        b.customer = nil       // current served customer
        b.Unlock()
        time.Sleep(time.Millisecond * 100)
        select {
        case c := <-wr:        // cuts hair of first person in queue
            HairCut(c, b)
        default:               // if waiting room is empty
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
Checking waiting room: 0
Sleeping Barber ZzzzZzz - <nil>
Customer c020 comes in to: sleeping barber, room: 0, wake: 0 - customer: <nil>
Sleeping barber c020, room: 0, wake: 0
Woken by c020
Cutting  c020's hair
Customer 4000 comes in to:  cutting barber, room: 0, wake: 0 - customer: c020
Customer 4010 comes in to:  cutting barber, room: 1, wake: 0 - customer: c020
Checking waiting room: 1
Customer 2010 comes in to: checking barber, room: 1, wake: 0 - customer: <nil>
Customer 4080 comes in to: checking barber, room: 3, wake: 0 - customer: <nil>
Cutting  4000's hair
Customer 40d0 comes in to:  cutting barber, room: 2, wake: 0 - customer: 4000
Customer 2070 comes in to:  cutting barber, room: 4, wake: 0 - customer: 4000
Checking waiting room: 4
Customer 20d0 comes in to: checking barber, room: 4, wake: 0 - customer: <nil>
Customer 2110 comes in to: checking barber, room: 6, wake: 0 - customer: <nil>
Cutting  4010's hair
Customer 2160 comes in to:  cutting barber, room: 5, wake: 0 - customer: 4010
Checking waiting room: 7
Cutting  2010's hair
Checking waiting room: 6
Cutting  4080's hair
Checking waiting room: 5
Cutting  40d0's hair
Checking waiting room: 4
Cutting  2070's hair
Checking waiting room: 3
Cutting  20d0's hair
Checking waiting room: 2
Cutting  2110's hair
Checking waiting room: 1
Cutting  2160's hair
Checking waiting room: 0
No more customers for the day
```
