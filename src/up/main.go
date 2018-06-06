package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	update "github.com/inconshreveable/go-update"
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

	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(0)
	}

	if os.Args[1] == "update" {
		systemOS := runtime.GOOS
		arch := runtime.GOARCH

		err := doUpdate("https://tools.liftedkilt.xyz/api/download/" + systemOS + "/" + arch + "/up")
		checkErr(err)
		// If no error
		fmt.Println("Update was successful")
		os.Exit(0)
	} else if os.Args[1] == "version" {
		available, latest := updateAvailable()
		local := getVersion()

		if available {
			fmt.Println("Local:", local, " Latest:", latest)
		} else {
			fmt.Println("No update available. Current version is:", local)
		}
		os.Exit(0)

	}

	lines, err := readLines(list)
	checkErr(err)

	results := make(chan string, len(lines))
	hosts := listToHostStruct(lines, port)

	var wg sync.WaitGroup

	// Dispatch Goroutines to check for liveliness
	for _, h := range hosts {
		go isUpConcurrent(h, results, timeout, &wg)
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

func isUpConcurrent(h host, results chan string, t int, wg *sync.WaitGroup) {
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

func doUpdate(url string) error {
	resp, err := http.Get(url)
	checkErr(err)
	defer resp.Body.Close()

	err = update.Apply(resp.Body, update.Options{})
	checkErr(err)

	return err
}

func getVersion() string {
	return "1.0.0"
}

func updateAvailable() (bool, string) {
	resp, err := http.Get("https://raw.githubusercontent.com/liftedkilt/tools/master/src/up/version")
	checkErr(err)
	defer resp.Body.Close()

	raw, _ := ioutil.ReadAll(resp.Body)

	latest := string(raw)
	current := getVersion()

	if latest != current {
		return true, latest
	}

	return false, latest

}
