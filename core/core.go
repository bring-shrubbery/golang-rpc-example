package core

import (
	"fmt"
	"strconv"
	"time"
)

// Response stores response data.
type Response string

// Request stores request data.
type Request struct {
	ID  int
	Res int
}

// TaskHandler stores functions that server should be able to execute.
type TaskHandler struct {
	Name          string
	Busy          bool
	ResourceCap   int
	resourceUsage int
}

// PerformTaskName stores the full name of the PerformTask function.
const PerformTaskName = "TaskHandler.PerformTask"

// PerformTask will perform task, given the request, and response variables.
func (h *TaskHandler) PerformTask(args *Request, response *Response) error {
	h.Busy = true
	fmt.Printf("Received task %d with %d units!\n", args.ID, args.Res)
	h.resourceUsage = h.resourceUsage + args.Res

	// Pretend to do stuff for 'args.Res' seconds.
	time.Sleep(time.Duration(args.Res) * time.Second)

	h.resourceUsage = h.resourceUsage - args.Res

	*response = Response(fmt.Sprintf("%s: Task %d performed using %d units", h.Name, args.ID, args.Res))
	h.Busy = false
	return nil
}

// PortToString converts port number in uint format to a string format, appends colon in front.
func PortToString(port uint) string {
	return ":" + strconv.Itoa(int(port))
}
