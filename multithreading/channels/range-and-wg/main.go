package main

import (
	"fmt"
	"sync"
)

// Thread 1
func main() {
	ch := make(chan int)
	wg := &sync.WaitGroup{}
	wg.Add(10)
	defer wg.Wait()

	go publish(ch)
	go reader(ch, wg)

	// Com o wait group, a thread principal é segurada (pelo wait) e assim, as duas outras funções podem ser executadas em paralelo
}

// <- ao lado esquerdo do canal significa que o canal apenas irá receber dados (chamado receive only)
func reader(ch <-chan int, wg *sync.WaitGroup) {
	for i := range ch {
		fmt.Printf("Received: %d\n", i)
		wg.Done()
	}
}

// <- ao lado direito do canal significa que o canal apenas irá enviar dados (chamado send only)
func publish(ch chan<- int) {
	for i := 0; i < 10; i++ {
		ch <- i
	}
	close(ch)
}
