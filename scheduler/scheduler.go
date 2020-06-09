package scheduler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net"
	"net/rpc"

	"golang-rpc-example/core"
)

// ServerInfo stores information about servers.
type ServerInfo struct {
	config  ServerConfiguration
	load    int
	client  *rpc.Client
	isalive bool
}

// Scheduler handles the load balancing. It stores the information about all
type Scheduler struct {
	Port     uint
	Address  string
	Servers  []ServerInfo
	listener net.Listener
}

func getConfiguration(path string) (Config, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return Config{}, errors.New("no configuration was found")
	}

	config := Config{}

	err = json.Unmarshal([]byte(file), &config)

	return config, nil
}

func (scheduler *Scheduler) reestablishConnections() (err error) {
	for i, s := range scheduler.Servers {
		fullAddress := s.config.Address + core.PortToString(s.config.Port)
		scheduler.Servers[i].client, err = rpc.Dial("tcp", fullAddress)

		if err != nil {
			// Couldn't connect to one of the servers!
			fmt.Printf("Could not connect to the server %s\n", s.config.Name)
			continue
		}

		// Connection successful.
		fmt.Println("Sucessfully connected to server", s.config.Name)
		scheduler.Servers[i].isalive = true
	}

	return
}

// Init initialises new scheduler and makes connections to all the server nodes.
func (scheduler *Scheduler) Init() (err error) {
	// Get configuration for the available servers.
	var config Config
	config, err = getConfiguration("schedulerConfig.json")
	if err != nil {
		println(err.Error())
		return
	}

	// Initialise array of ServerInfo structs, instead of just
	// ServerConfiguration structs
	for _, c := range config.AvailableServers {
		scheduler.Servers = append(scheduler.Servers, ServerInfo{
			config: c,
		})
	}

	// Prepare server connections
	scheduler.reestablishConnections()

	// Register self with the rpc, so we can call the method execution.
	// The current method will not be registered, because it does not
	// conform to go's rpc requirements for regestering methods, which is what we want.
	rpc.Register(scheduler)

	// And start listenting.
	scheduler.listener, err = net.Listen("tcp", core.PortToString(scheduler.Port))
	scheduler.Address = scheduler.listener.Addr().String()

	if err != nil {
		// TODO: Handle error.
	}

	// Close all connections automatically, when this function returns.
	defer scheduler.close()

	// Wait for listener to get a request.
	rpc.Accept(scheduler.listener)

	return
}

func (scheduler *Scheduler) close() (err error) {
	// Close connections to all servers.
	for _, s := range scheduler.Servers {
		s.client.Close()
	}
	return
}

// ExecuteTaskName is needed for the client side. Client runs this function
// through an rpc, then Scheduler class calls whichever task is needed.
const ExecuteTaskName = "Scheduler.ExecuteTask"

// ExecuteTask is used for transparency with the client. The client requests
// the task to be executed, but doesn't know which node exactly will execute
// the task. This function handles all the decisions for the
func (scheduler *Scheduler) ExecuteTask(request *core.Request, response *core.Response) (err error) {
	minRes := math.MaxInt32
	var lbs int // Least busy server

	for i, s := range scheduler.Servers {
		if s.isalive && s.load < minRes {
			minRes = s.load
			lbs = i
		}
	}

	// Print information about which server was chosen for the task and other info.
	println(
		"Delegating task",
		request.ID,
		"to server",
		scheduler.Servers[lbs].config.Name,
		"with",
		request.Res,
		"units,",
		"load before:",
		scheduler.Servers[lbs].load,
	)

	scheduler.Servers[lbs].load += request.Res

	res := new(core.Response)

	// Assign the task to the least busy server. Note that arguments from the user and response are piped.
	err = scheduler.Servers[lbs].client.Call(core.PerformTaskName, request, res)
	if err != nil {
		println(err.Error())
		return
	}

	*response = *res

	scheduler.Servers[lbs].load -= request.Res

	fmt.Printf("Server %s is done executing task %d with %d units.\n", scheduler.Servers[lbs].config.Name, request.ID, request.Res)

	return
}
