# Cigarette Smokers

# Problem

There are three smokers around a table, each with unlimited supply of either paper, tobacco, or paper. A fourth party, with an unlimited supply of everything, chooses at random a smoker, and put on the table the supplies needed for a cigarrette. The chosen smoker smokes, and the process should repeat indefinitely.

# Solution

This solution sends a signal to all smokers, telling whose turn it is to smoke, and then places the necessary inputs into the `Table` by the tables three _thread safe_ **channels**. The smokers don't reach for anything until the `signal` channel stops blocking, and then only the chosen smoker will access the table.

```
func arbitrator(t *Table, smokers [3]chan int) {
    for {
        time.Sleep(time.Millisecond * 500)
        next := rand.Intn(3)              // choose next smoker
        switch next {
        case paper:
            t.grass <- 1                  // put the proper
            t.match <- 1                  // ingredients
        case grass:                       // on the table
            t.paper <- 1
            t.match <- 1
        case match:
            t.grass <- 1
            t.paper <- 1
        }
        for _, smoker := range smokers {
            smoker <- next               // send the signal
        }
        wg.Add(1)                        // wait for them to light up
        wg.Wait()
    }
}

func smoker(t *Table, name string, smokes int, signal chan int) {
    var chosen = -1
    for {
        chosen = <-signal         // blocks
        if smokes != chosen {
            continue
        }
        select {                  // consume first item
        case <-t.paper:           // blocks
        case <-t.grass:
        case <-t.match:
        }
        select {                  // consume second item
        case <-t.paper:
        case <-t.grass:
        case <-t.match:
        }
        time.Sleep(time.Millisecond * 100)
        wg.Done()                 // aribitrator can move on
    }
}
```


# Example Output

```
Sandy, Apple, Daisy, sit with
paper, grass, match

Table chooses match: Daisy
Table: 1 grass: 1 match: 0
Table: 1 grass: 0 match: 0
Table: 0 grass: 0 match: 0
Daisy smokes a cigarette
Table chooses paper: Sandy
Table: 0 grass: 1 match: 1
Table: 0 grass: 0 match: 1
Table: 0 grass: 0 match: 0
Sandy smokes a cigarette
Table chooses match: Daisy
Table: 1 grass: 1 match: 0
Table: 1 grass: 0 match: 0
Table: 0 grass: 0 match: 0
Daisy smokes a cigarette
Table chooses match: Daisy
Table: 1 grass: 1 match: 0
Table: 1 grass: 0 match: 0
Table: 0 grass: 0 match: 0
Daisy smokes a cigarette
Table chooses grass: Apple
Table: 1 grass: 0 match: 1
Table: 1 grass: 0 match: 0
Table: 0 grass: 0 match: 0
Apple smokes a cigarette
Table chooses paper: Sandy
Table: 0 grass: 1 match: 1
Table: 0 grass: 1 match: 0
Table: 0 grass: 0 match: 0
Sandy smokes a cigarette
```
