package main

import (
	"fmt"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
)

func worker(protocol string, addr string, ports, results chan int) {
	for p := range ports {
		address := strings.Join([]string{addr, ":", strconv.Itoa(p)}, "")
		conn, err := net.Dial(protocol, address)
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- p
	}
}

func main() {
	args := os.Args
	if len(args) < 3 {
		fmt.Println("error using port scanner")
		fmt.Println("usage: go run portscanner.go <protocol> <URL>")
		os.Exit(2)
	}

	ports := make(chan int, 100)
	results := make(chan int, 100)
	var openports []int

	for i := 0; i < cap(ports); i++ {
		go worker(args[1], args[2], ports, results)
	}

	go func() {
		for i := 1; i <= 1024; i++ {
			ports <- i
		}
	}()

	for i := 0; i < 1024; i++ {
		port := <-results
		if port != 0 {
			openports = append(openports, port)
		}
	}

	close(ports)
	close(results)
	sort.Ints(openports)
	for _, port := range openports {
		fmt.Println(port, "open")
	}
}
