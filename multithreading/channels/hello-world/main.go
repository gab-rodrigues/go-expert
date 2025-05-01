package main

// Thread 1
func main() {
	//channels
	// Fazer comunicação entre threads. Nesse exemplo entre o 1 (main) e o 2 (anonymous function)
	// Segurança para uma thread saber em que momento ela pode trabalhar com determinado dado

	channel := make(chan string) // canal vazio

	// Thread 2
	go func() {
		channel <- "Hello World, dear Gabriel" // canal cheio
	}()

	msg := <-channel // vazio
	println(msg)
}
