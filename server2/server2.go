package server2

import (
	"net"
	"net/rpc"

	"golang-rpc-example/core"
)

// Server stores the information about current server.
type Server struct {
	Name     string
	Port     uint
	Address  string
	listener net.Listener
}

// Init sets up the task handler and starts listenting.
func (s *Server) Init() (err error) {
	// Create an instance of task handler for current server.
	taskHandler := &core.TaskHandler{
		Name: s.Name,
	}

	// Register the task handler with RPC.
	rpc.Register(taskHandler)

	// Start listenting on specified port.
	s.listener, err = net.Listen("tcp", core.PortToString(s.Port))
	s.Address = s.listener.Addr().String()

	// And handle error if it occurs.
	if err != nil {
		return
	}

	// Close listener when this function finishes executing.
	defer s.listener.Close()

	rpc.Accept(s.listener)
	return
}
