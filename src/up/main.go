package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

type host struct {
	name string
	port string
}

func main() {

	list := os.Args[1]
	port := os.Args[2]

	lines, err := readLines(list)
	checkErr(err)

	results := make(chan string, len(lines))

	var wg sync.WaitGroup

	// Dispatch Goroutines to check for liveliness

	hosts := listToHostStruct(lines, port)

	for _, h := range hosts {
		go isUpConc(h, results, &wg)
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

func isUp(h host) bool {
	remote := h.name + ":" + h.port
	timeout := time.Duration(500) * time.Millisecond

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

func isUpConc(h host, results chan string, wg *sync.WaitGroup) {

	result := isUp(h)
	var status string
	if result {
		status = h.name + " up"
	} else {
		status = h.name + " down"
	}
	wg.Done()
	results <- status
}

func listToHostStruct(list []string, port string) []host {
	var hosts []host

	for _, item := range list {
		hosts = append(hosts, host{
			name: item,
			port: port,
		})
	}

	return hosts
}
