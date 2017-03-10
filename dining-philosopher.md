# Dining Philosophers

Example output:

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
