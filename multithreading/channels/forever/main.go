package main

// Thread 1
func main() {
	forever := make(chan bool) //canal vazio

	go func() {
		for i := 0; i < 10; i++ {
			println(i)
		}
		//forever <- true (com isso, o canal fica cheio e não resulta mais em deadlock)
	}()

	<-forever // esperando o canal estar cheio para continuar

	// resultado: deadlock!
}
