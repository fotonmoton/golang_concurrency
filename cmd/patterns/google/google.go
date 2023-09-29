package main

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	Web   = fakeSearch("web")
	Image = fakeSearch("image")
	Video = fakeSearch("video")
)

type Result = string

type Search func(query string) Result

func fakeSearch(kind string) Search {
	return func(query string) Result {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		return Result(fmt.Sprintf("%s result for %q\n", kind, query))
	}
}

func Google(query string) []Result {
	results := make([]Result, 3)

	results = append(results, Web(query)) // Currently, we block in each of these search queries
	results = append(results, Image(query))
	results = append(results, Video(query))

	return results
}

func FasterGoogle(query string) []Result {
	c := make(chan Result)
	results := make([]Result, 3)

	// Spin each of the searches off into their own goroutine and pipe all results back to the channel created above
	go func() { c <- Web(query) }()
	go func() { c <- Image(query) }()
	go func() { c <- Video(query) }()

	// Pull each result out of the channel as it becomes available and append it to the results splice returned
	for i := 0; i < 3; i++ {
		result := <-c
		results = append(results, result)
	}

	return results
}

func TimeoutGoogle(query string, timeoutDuration time.Duration) []Result {
	c := make(chan Result)
	results := make([]Result, 3)

	// Spin each of the searches off into their own goroutine and pipe all results back to the channel created above
	go func() { c <- Web(query) }()
	go func() { c <- Image(query) }()
	go func() { c <- Video(query) }()

	// Now we pull each Result from the channel, but timeout if searches take longer than 80ms
	timeout := time.After(timeoutDuration * time.Millisecond)
	for i := 0; i < 3; i++ {
		select {
		case result := <-c:
			results = append(results, result)
		case <-timeout:
			fmt.Println("Timed out")
			return results
		}
	}

	return results
}

// Returns the first response received from one of the many replica Search functions received
func First(query string, replicas ...Search) Result {
	c := make(chan Result)
	searchReplica := func(i int) { c <- replicas[i](query) } // A function that abstract calling a given Search function with the given query, then piping the result into channel c

	// Iterate over all the given replica Search functions and spin off into their own goroutines
	for i := range replicas {
		go searchReplica(i)
	}

	// Return the first result received only
	return <-c
}

func ReplicasGoogle(query string, timeoutDuration time.Duration) (results []Result) {
	c := make(chan Result)

	go func() { c <- First(query, Web, Web) }()
	go func() { c <- First(query, Image, Image) }()
	go func() { c <- First(query, Video, Video) }()

	timeout := time.After(timeoutDuration * time.Millisecond)
	for i := 0; i < 3; i++ {
		select {
		case result := <-c:
			results = append(results, result)
		case <-timeout:
			fmt.Println("Timed out")
			return
		}
	}

	return results
}

func main() {
	start := time.Now()
	results := Google("prjctr")
	// results := FasterGoogle("prjctr")
	// results := TimeoutGoogle("prjctr", 80)
	// results := ReplicasGoogle("prjctr", 80)
	elapsed := time.Since(start)
	fmt.Println(results)
	fmt.Println(elapsed)
}
