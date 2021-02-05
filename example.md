### goroutine

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	go spinner(100 * time.Millisecond)
	const n = 45
	fibN := fib(n)
	fmt.Printf("\rFibonacci(%d) = %d\n", n, fibN)
}

func spinner(delay time.Duration) {
	for {
		for _, r := range `-\|/` {
			fmt.Printf("\r%c", r)
			time.Sleep(delay)
		}
	}
}

func fib(x int) int {
	if x < 2 {
		return x
	}
	return fib(x-1) + fib(x-2)
}
```



### channel  


```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	c := make(chan int32, 12)
	a := make([]int32, 12)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for v := range c {
			fmt.Println(v)
		}
		wg.Done()
	}()
	for i := 0; i < len(a); i++ {
		a[i] = int32(i + 1)
		c <- a[i]
	}
	close(c)
	wg.Wait()
}
```
