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

func main() {

	url := flag.String("url", "", "url to connect to")
	n := flag.Int("n", 1, "total number of requests")
	c := flag.Int("c", 1, "number of parallel requests")

	verbose := flag.Int("v", 0, "verbosity")

	flag.Parse()

	VERBOSE = *verbose

	// Enable debug in gospdy code
	if VERBOSE > 10 {
		spdy.Log = func(format string, args ...interface{}) {
		    fmt.Printf(format, args...)
		}
	}

	// Setting Max Procs to the Number of CPU Cores
	fmt.Printf("Max procs %d\n", runtime.GOMAXPROCS(runtime.NumCPU()))
	fmt.Printf("Max procs %d\n", runtime.GOMAXPROCS(0))

	loop := *n / *c

	done := make(chan bool, *n)
	start := make(chan bool, *n)

	b := time.Now().UnixNano()

	for i := 0; i < *c; i++ {
		go connect(start, done, *url)
		start <- true // start some goroutines immediately
	}
	for i := *c; i < *n; i++ {
		start <- true // fill in the chan so everybody can work
	}

	for i := 0; i < *n; i++ {
		<-done
		if i%loop == 0 && i != 0 {
			fmt.Printf("%d reqs end in %d msecs\n", i, (time.Now().UnixNano()-b)/1e6)
		}
	}
	e := time.Now().UnixNano()

	fmt.Printf("time used: %dmsecs\n", (e-b)/1e6)
}
