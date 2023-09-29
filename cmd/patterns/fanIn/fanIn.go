package main

import (
	"fmt"
	"math/rand"
	"time"
)

func boring(msg string) <-chan string { // returns receive-only channel of strings
	c := make(chan string)

	go func() { // we launch the goroutine from inside the function
		for i := 0; ; i++ {
			c <- fmt.Sprintf("%s, %d", msg, i)
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
		}
	}()

	return c // return the channel to the caller
}

func fanIn(joe, ann <-chan string) <-chan string {
	c := make(chan string)

	go func() {
		for {
			c <- <-joe
		}
	}()

	go func() {
		for {
			c <- <-ann
		}
	}()

	return c
}

func fanInSelect(joe, ann <-chan string) <-chan string {
	allMessages := make(chan string)

	go func() {
		for {
			select {
			case joeMessage := <-joe:
				allMessages <- joeMessage
			case annMessage := <-ann:
				allMessages <- annMessage
			}
		}
	}()

	return allMessages
}

func main() {
	allMessages := fanIn(boring("Joe"), boring("Ann"))

	for i := 0; i < 10; i++ {
		fmt.Println(<-allMessages)
	}

	fmt.Println("You're both boring; I'm leaving")
}
