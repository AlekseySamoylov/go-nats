package main

import "fmt"

func main() {
	channel := make(chan string, 1)

	channel <- "Hello"

	msg := <-channel

	fmt.Println(msg)
}
