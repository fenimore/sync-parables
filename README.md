# Concurrency Practice Problem Sets

- [x] Dining Philosophers
- [x] Sleeping Barber
- [x] Cigarrette Smokers

See the directory readmes

## Dining Philosophers

`N` philosophers eat around a table together, with `N` chopsticks, one between each pair of philosophers. The philosophers pick up one chopstick, another, eat, and then put the chopsticks down and think. Eventually, each philosopher picks up the chopstick to their left and **deadlocks**, waiting for their right chopstick. Unable to eat, they starve.

## Sleeping Barber

A barber checks on the waiting room and then either cuts hair or goes to sleep (if there is no one in the waiting room). Concurrently, the customer checks on the barber and if the barber is sleeping, the barber wakes up.

If the customer checks on the barber when the barber is checking on an _empty_ waiting room, the barber would go back to sleep and the customer would go wait, possibly **forever**.

## Cigarette Smokers

There are three smokers around a table, each with unlimited supply of either paper, tobacco, or paper. A fourth party, with an unlimited supply of everything, chooses at random a smoker, and put on the table the supplies needed for a cigarrette. The chosen smoker smokes, and the process should repeat indefinitely.
