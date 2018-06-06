package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

type host struct {
	name string
	port int
}

func main() {

	// Parse flags
	var list string
	var port int
	var timeout int
	flag.IntVar(&port, "port", 80, "port to check")
	flag.IntVar(&timeout, "timeout", 500, "timeout for check in milliseconds")
	flag.StringVar(&list, "file", "", "name of input file")
	flag.Parse()

	lines, err := readLines(list)
	checkErr(err)

	results := make(chan string, len(lines))
	hosts := listToHostStruct(lines, port)

	var wg sync.WaitGroup

	// Dispatch Goroutines to check for liveliness
	for _, h := range hosts {
		go isUpConc(h, results, timeout, &wg)
		wg.Add(1)
	}

	// Wait for all goroutines to return
	wg.Wait()
	// Close channel
	close(results)

	for line := range results {
		fmt.Println(line)
	}
}

/*
/ Takes a path to a file, returns the contents of the file as a slice of strings and an error object
*/
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	checkErr(err)

	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

/*
/ Takes a host (hostname or IP address) and a port, returns a bool response if that host is accepting connections
/ on that port
*/

func isUp(h host, t int) bool {
	remote := h.name + ":" + strconv.Itoa(h.port)
	timeout := time.Duration(t) * time.Millisecond

	var status bool

	conn, err := net.DialTimeout("tcp", remote, timeout)
	if err != nil {
		status = false
	} else {
		status = true
		defer conn.Close()
	}

	return status
}

/*
/ Handles error objects and exits if needed
*/
func checkErr(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

/*
/ Concurrent function that returns status of isUp via a channel
*/

func isUpConc(h host, results chan string, t int, wg *sync.WaitGroup) {

	result := isUp(h, t)
	var status string
	if result {
		status = h.name + " up"
	} else {
		status = h.name + " down"
	}
	wg.Done()
	results <- status
}

func listToHostStruct(list []string, port int) []host {
	var hosts []host

	for _, item := range list {
		hosts = append(hosts, host{
			name: item,
			port: port,
		})
	}

	return hosts
}
