package main

import (
	"fmt"
	"simulator/pkg/algorithm/fifo"
	"simulator/pkg/consts"
	"simulator/pkg/job"
	"simulator/pkg/node"
)

var Nodes []node.Node

func initInfra() {
	Nodes = make([]node.Node, consts.NodeCnt)
	for i := range Nodes {
		Nodes[i] = node.NewNode(i, 64, 512, 8, 2)
	}
}

var Jobs []job.Job

func initJob() {
	// build Jobs from trace
}

var schedulePolicy string

func main() {
	fmt.Println("begin")
	initInfra()
	initJob()
	schedulePolicy = "fifo"

	// main routine
	for currentTime := 0; currentTime < consts.MaxTime; currentTime += consts.SessionInterval {
		if schedulePolicy == "fifo" {
			fifo.SessionRun(currentTime ,Jobs, Nodes)
		} else if schedulePolicy == "lease" {

		}
	}
}
