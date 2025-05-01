package main

import (
	"fmt"
	"time"
)

func task(value string) {
	for i := 0; i < 10; i++ {
		fmt.Printf("%d: Task %s is running\n", i, value)
		time.Sleep(1 * time.Second)
	}
}

// Thread 1
func main() {

	// Adiciona quantidade de tarefas
	// Informa que uma operação foi finalizada
	// Esperar até que as operações sejam finalizadas
	//wg := sync.WaitGroup{}
	// Thread 2
	go task("A")
	// Thread 3
	go task("B")
	// Thread 4
	go func() {
		for i := 0; i < 5; i++ {
			fmt.Printf("%d: Task %s is running\n", i, "anonymous")
			time.Sleep(1 * time.Second)
		}
	}()

	time.Sleep(15 * time.Second)

}
