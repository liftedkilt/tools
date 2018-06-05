package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

func main() {
	host := os.Args[1]
	port := os.Args[2]

	lines, err := readLines(host)
	checkErr(err)

	results := make(chan string, len(lines))

	var wg sync.WaitGroup

	// Dispatch Goroutines to check for liveliness
	for _, value := range lines {
		go isUpConc(value, port, results, &wg)
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

func isUp(host, port string) bool {
	remote := host + ":" + port
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

func isUpConc(host string, port string, results chan string, wg *sync.WaitGroup) {

	result := isUp(host, port)
	var status string
	if result {
		status = host + " up"
	} else {
		status = host + " down"
	}
	wg.Done()
	results <- status
}
