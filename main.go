package main

import (
	"flag"
	"fmt"

	"golang-rpc-example/client"
	"golang-rpc-example/scheduler"
	"golang-rpc-example/server1"
	"golang-rpc-example/server2"
)

var (
	port        = flag.Uint("p", 4000, "port to listen on")
	isServer1   = flag.Bool("s1", false, "run this node as server number 1")
	isServer2   = flag.Bool("s2", false, "run this node as server number 2")
	isScheduler = flag.Bool("scheduler", false, "should this be a scheduler?")
)

// Creates and runs server 1
func runServer1() {
	server := &server1.Server{
		Name: "N1",
		Port: *port,
	}

	server.Init()
}

// Creates and runs server 2.
func runServer2() {
	server := &server2.Server{
		Name: "N2",
		Port: *port,
	}

	server.Init()
}

// Creates and runs scheduler.
func runScheduler() {
	scheduler := &scheduler.Scheduler{
		Port: *port,
	}

	scheduler.Init()
}

// Creates and runs client.
func runClient() {
	client := &client.Client{
		ServerURL: "[::]",
		Port:      *port,
	}

	client.Init()
}

func main() {
	// Parse flags before proceeding.
	flag.Parse()

	// If server1 flag is provided, run server, then quit.
	if *isServer1 {
		fmt.Println("Starting server 1!")
		runServer1()
		return
	}

	// If server2 flag is provided, run server, then quit.
	if *isServer2 {
		fmt.Println("Starting server 2!")
		runServer1()
		return
	}

	// If scheduler flag is provided, run scheduler, then quit.
	if *isScheduler {
		fmt.Println("Starting scheduler!")
		runScheduler()
		return
	}

	// If no flags are provided, run client.
	fmt.Println("Starting client!")
	runClient()
}
