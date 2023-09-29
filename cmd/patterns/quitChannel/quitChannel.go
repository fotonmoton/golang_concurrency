package main

import (
	"fmt"
	"time"
)

func boring(msg string, quit chan string) <-chan string {
	c := make(chan string)

	go func() { // we launch the goroutine from inside the function
		for i := 0; ; i++ {
			select {
			case c <- fmt.Sprintf("%s, %d", msg, i):
				// Do nothing
			case quitMsg := <-quit:
				fmt.Println("finishing...: ", quitMsg)
				time.Sleep(5 * time.Second)
				fmt.Println("finished!")
				quit <- "all clean up is done!"
				return // Parent routine tells us to finish, so we return from the goroutine
			}
		}
	}()

	return c // return the channel to the caller
}

func main() {
	quit := make(chan string)
	c := boring("Joe", quit)
	for i := 5; i >= 0; i-- {
		fmt.Println(<-c)
	}

	// change to make it wait
	quit <- "go home" // Tell the routine to finish
	fmt.Println("Work is done: ", <-quit)
}
