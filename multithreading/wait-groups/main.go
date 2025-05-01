package main

import (
	"fmt"
	"sync"
	"time"
)

func task(value string, wg *sync.WaitGroup) {
	for i := 0; i < 10; i++ {
		fmt.Printf("%d: Task %s is running\n", i, value)
		time.Sleep(1 * time.Second)
		wg.Done()
	}
}

// Thread 1
func main() {

	// Adiciona quantidade de tarefas
	// Informa que uma operação foi finalizada
	// Esperar até que as operações sejam finalizadas
	// Sincronizar operações meapeadas
	wg := &sync.WaitGroup{}
	wg.Add(25)
	defer wg.Wait()

	// Thread 2
	go task("A", wg)
	// Thread 3
	go task("B", wg)
	// Thread 4
	go func() {
		for i := 0; i < 5; i++ {
			fmt.Printf("%d: Task %s is running\n", i, "anonymous")
			time.Sleep(1 * time.Second)
			wg.Done()
		}
	}()
}
