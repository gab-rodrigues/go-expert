package main

import "time"

func main() {
	c1 := make(chan int)
	c2 := make(chan int)

	go func() {
		time.Sleep(1 * time.Second)
		c1 <- 1
	}()

	go func() {
		time.Sleep(2 * time.Second)
		c2 <- 2
	}()

	select {
	case msg1 := <-c1: // rabbitmq or api 1
		println("received", msg1)

	case msg2 := <-c2: // kafka or api 2
		println("received", msg2)

	case <-time.After(3 * time.Second):
		println("timeout")
	}
}
