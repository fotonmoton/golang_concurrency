package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/fotonmoton/golang_concurrency/crawler"
	safecounter "github.com/fotonmoton/golang_concurrency/safeCounter"
)

/*

Concurrency
1. Sequential execution
2. Concurrent execution
3. Parallel execution
4. Go will not wait any go routine
5. Sync mechanisms in go
6. WaitGroup
7. Mutex
6. Goroutines
7. Difference between goroutines and threads

*/

func say(s string, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < 5; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Println(s)
	}
}

func HelloSequential() {
	wg := &sync.WaitGroup{}
	wg.Add(2)
	say("hello", wg)
	say("world", wg)
}

func HelloConcurrent() {

	wg := &sync.WaitGroup{}
	wg.Add(3)

	go say("world", wg)
	go say("hello", wg)

	wg.Wait()
}

func Loops() {
	f := func(n int) {
		for i := 0; i < 10; i++ {
			fmt.Println(n, ":", i)
			amt := time.Duration(rand.Intn(250))
			time.Sleep(time.Millisecond * amt)
		}
	}

	for i := 0; i < 10; i++ {
		go f(i)
	}
	var input string
	fmt.Scanln(&input)
}

func sum(s []int, c chan int) {
	sum := 0
	for _, v := range s {
		sum += v
	}
	c <- sum // send sum to c
}

func Sum() {

	s := []int{7, 2, 8, -9, 4, 0}

	c := make(chan int, 2)

	go sum(s[:3], c)
	go sum(s[3:], c)
	x, y := <-c, <-c // receive from c

	fmt.Println(x, y, x+y)
}

func BufferedChannels() {
	ch := make(chan int, 2)
	ch <- 1
	ch <- 2
	fmt.Println(<-ch)
	fmt.Println(<-ch)
}

func fibonacciEasy(n int, c chan int) {
	x, y := 0, 1
	for i := 0; i < n; i++ {
		amt := time.Duration(rand.Intn(500))
		time.Sleep(time.Millisecond * amt)
		c <- x
		x, y = y, x+y
	}
	close(c)
}

func ClosingChannel() {
	c := make(chan int, 10)
	go fibonacciEasy(cap(c), c)
	for i := range c {
		fmt.Println(i)
	}

	// close(c)

}

func fibonacciSelect(results, quit chan int) {
	x, y := 0, 1
	for {
		select {
		case results <- x:
			x, y = y, x+y
		case <-quit:
			fmt.Println("quit")
			return
		}
	}
}

func reader(c, quit chan int) {
	for i := 0; i < 10; i++ {
		fmt.Println(<-c)
	}
	quit <- 0
}

func SelectingOnChanel() {
	results := make(chan int)
	quit := make(chan int)

	go reader(results, quit)
	fibonacciSelect(results, quit)
}

func DefaultSelection() {
	tick := time.Tick(100 * time.Millisecond)
	boom := time.After(500 * time.Millisecond)
	for {
		select {
		case <-tick:
			fmt.Println("tick.")
		case <-boom:
			fmt.Println("BOOM!")
			return
		default:
			fmt.Println("    .")
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func Mutex() {
	c := safecounter.NewCounter()

	for i := 0; i < 1000; i++ {
		go c.Inc("somekey")
	}

	time.Sleep(time.Second)
	fmt.Println(c.Value("somekey"))
}

func Crawler() {
	crawler.Crawl("https://golang.org/", 4, crawler.FakeFetcher{
		"https://golang.org/": &crawler.FakeResult{
			"The Go Programming Language",
			[]string{
				"https://golang.org/pkg/",
				"https://golang.org/cmd/",
			},
		},
		"https://golang.org/pkg/": &crawler.FakeResult{
			"Packages",
			[]string{
				"https://golang.org/",
				"https://golang.org/cmd/",
				"https://golang.org/pkg/fmt/",
				"https://golang.org/pkg/os/",
			},
		},
		"https://golang.org/pkg/fmt/": &crawler.FakeResult{
			"Package fmt",
			[]string{
				"https://golang.org/",
				"https://golang.org/pkg/",
			},
		},
		"https://golang.org/pkg/os/": &crawler.FakeResult{
			"Package os",
			[]string{
				"https://golang.org/",
				"https://golang.org/pkg/",
			},
		},
	})
}

func main() {
	// HelloSequential()
	// HelloConcurrent()
	// Sum()
	// Loops()
	// BufferedChannels()
	// ClosingChannel()
	// SelectingOnChanel()
	// DefaultSelection()
	// Mutex()
	// Crawler()
}
