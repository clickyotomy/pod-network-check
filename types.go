package main

import "time"

// node holds the metadata about the address and port
// the check should connert to.
type node struct {
	ID       string      `json:"id"`
	Address  string      `json:"address"`
	Port     uint16      `json:"port"`
	Metadata interface{} `json:"metadata,omitempty"`
}

// annotation holds the annotations for a pod.
type annotation struct {
	Name      string        `json:"name,omitempty"`
	Version   string        `json:"version,omitempty"`
	Metadata  interface{}   `json:"metadata,omitempty"`
	Endpoints []interface{} `json:"endpoints,omitempty"`
	Nodes     []node        `json:"nodes"`
}

// podSock is the pod's name and the socket address for dialing.
type podSock struct {
	name    string
	addr    string
	port    uint16
	timeout time.Duration
}

// status holds the check status.
type status struct {
	pod  string
	err  error
	pass bool
}
