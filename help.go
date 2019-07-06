package main

import (
	"fmt"
	"os"
	"strings"
)

// usage displays the command line usage.
func usage() {
	fmt.Fprintf(
		os.Stderr,
		"%s: Send pod connection failure events to DataDog.\n"+
			"\n"+
			"usage:\n"+
			"  %s -name {C} -namespace {N} -ports {X,Y,Z...}\n"+
			"  %s -dogstatsd {D} -interval {I} -timeout {T}\n"+
			"  %s -protocol {P} -debug\n"+
			"\n"+
			"arguments:\n"+
			" -name        A name for the check; used the aggregation key.\n"+
			" -namespace   The `kubernetes' namespace to check for pods.\n"+
			" -ports       Comma separated list of ports to check.\n"+
			" -dogstatsd   Address to the `dogstatsd' server.\n"+
			" -interval    Check run interval, in seconds.\n"+
			" -timeout     Dial timeout for pods, in seconds.\n"+
			" -protocol    Protocol to use for the check.\n"+
			" -debug       Print debug output.\n"+
			"\n"+
			"defaults:\n"+
			" -name        %s\n"+
			" -interval    %d\n"+
			" -timeout     %d\n"+
			" -protocol    %s\n"+
			" -debug       %v\n",
		program, program,
		strings.Repeat(" ", len(program)), strings.Repeat(" ", len(program)),
		ddAggKey, checkInterval, checkPodTimeout, checkProtocol, false,
	)
}
