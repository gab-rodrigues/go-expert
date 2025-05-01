package main

import (
	"fmt"
	"time"
)

func worker(id int, ch <-chan int) {
	for x := range ch {
		fmt.Printf("Worker %d received %d\n", id, x)
		time.Sleep(1 * time.Second)
	}
}

func main() {
	data := make(chan int)
	workersQuantity := 10000
	// criar muitas threads em Go é muito mais tranquilo do que em outras linguagens. Gasto de memória muito baixo.
	for i := 0; i < workersQuantity; i++ {
		go worker(i, data)
	}

	for i := 0; i < 100000; i++ {
		data <- i
	}
}
