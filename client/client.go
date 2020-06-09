package client

import (
	"bufio"
	"fmt"
	"net/rpc"
	"os"
	"strconv"
	"strings"

	"golang-rpc-example/core"
	"golang-rpc-example/scheduler"
)

// Client contains data about the client.
type Client struct {
	ServerURL string
	Port      uint
	client    *rpc.Client
}

// Init initialises new client. Pretty much the only function that can be executed here.
func (c *Client) Init() (err error) {
	var latestID int

	// First, initialise the connection to the server.
	fullAddress := c.ServerURL + core.PortToString(c.Port)

	// Attempt to connect and save into client parameter of Client struct.
	c.client, err = rpc.Dial("tcp", fullAddress)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		// Request a task size.
		fmt.Println("Enter number of units that next task requires: ")
		unitsStr, err := reader.ReadString('\n')
		unitsStr = strings.Trim(unitsStr, "\n ")
		units, numerr := strconv.Atoi(unitsStr)

		// Handle errors.
		if numerr != nil {
			fmt.Println("\u001b[31mPlease enter a number!\u001b[0m")
			continue
		}

		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		// Make request to the scheduler via an immidiately invoked goroutine,
		// so that next request could be made immidiately after without waiting.
		go func() {
			// Create request T<ID,Res>
			request := &core.Request{
				ID:  latestID,
				Res: units,
			}

			// Prepare place to store response.
			response := new(core.Response)

			// Note: We're requesting scheduler's function, but passing the request
			//       and response as if we're talking to the server executing the task directly.
			err = c.client.Call(scheduler.ExecuteTaskName, request, response)
			if err != nil {
				fmt.Println("Call error:", err.Error())
				return
			}

			// Print response to the user
			fmt.Println(*response)
		}()

		latestID++
	}
}
