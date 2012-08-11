package main

import (
	"flag"
	"fmt"
	"github.com/jmckaskill/gospdy"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"time"
)

var VERBOSE = 0

// Connect loop
func connect(start, done chan bool, url string) {

	for {
		<-start

		tr := &spdy.Transport{}
		client := &http.Client{Transport: tr}

		r, err := client.Get(url)

		if err != nil {
			log.Fatal(err)
			done <- true
			return
		}

		defer r.Body.Close()

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			log.Fatal(err)
			done <- true
			return
		}

		// Print the body
		if VERBOSE > 0 {
			fmt.Printf(string(body))
		}

		done <- true
	}
}


// Goroutine to keep consuming <-done
func requests_done(done, end chan bool, n, c int) {

	start_time := time.Now().UnixNano()
	last_now := start_time

	for i := 0; i < n; i++ {
		<-done
		if i % c == 0 && i != 0 {
			now := time.Now().UnixNano()
			interval := (now - last_now) / 1e6
			fmt.Printf("%d reqs done. + %d msecs\n", i, interval)
			last_now = now
		}
	}

	fmt.Printf("%d reqs done.\n", n)

	total_time := (time.Now().UnixNano() - start_time) / 1e6
	fmt.Printf("\ntotal time: %d msecs\n", total_time)

	end <- true
}


func main() {

	url := flag.String("url", "", "url to connect to: 'https://www.google.com'")
	n := flag.Int("n", 1, "total number of requests")
	c := flag.Int("c", 1, "number of parallel requests")

	verbose := flag.Int("v", 0, "verbosity")

	flag.Parse()

	VERBOSE = *verbose

	// Enable debug in gospdy code
	if VERBOSE > 2 {
		spdy.Log = func(format string, args ...interface{}) {
			fmt.Printf(format, args...)
		}
	}

	start := make(chan bool, *n)
	done := make(chan bool, *n)


	end := make(chan bool, *n)

	// just one request
	if *n == 1 && *c == 1 {

		if VERBOSE == 0 {
			VERBOSE = 1
		}

		go connect(start, done, *url)
		start <- true
		<-done

		return
	}

	// Multiple Requests

	// Setting Max Procs to the Number of CPU Cores
	fmt.Printf("Max procs %d\n", runtime.GOMAXPROCS(runtime.NumCPU()))
	fmt.Printf("Max procs %d\n\n", runtime.GOMAXPROCS(0))

	go requests_done(done, end, *n, *c)

	for i := 0; i < *c; i++ {
		go connect(start, done, *url)
		// start some goroutines immediately
		start <- true
	}

	for i := *c; i < *n; i++ {
		// fill in the chan so everybody can work
		start <- true
	}

	// wait for all the requests to be terminated
	<-end
}
