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
			time.Sleep(time.Duration(1000) * time.Millisecond)
		}
	}()

	return c // return the channel to the caller
}

func main() {
	c := boring("Joe")
	timeout := time.After(5 * time.Second) // We kill the loop after 5sec total

	for {
		select {
		case s := <-c:
			fmt.Println(s)

		case <-timeout:
			fmt.Println("You talk too much")
			return
		}
	}
}

func one() {
	c := boring("Joe")

	for {
		select {
		case s := <-c:
			fmt.Println(s)

		case <-time.After(time.Duration(rand.Intn(1e3)+100) * time.Millisecond): // We kill the loop if we don't receive a message from the boring("Joe") within 1sec per iteration
			fmt.Println("You're too slow")
			return
		}
	}
}
