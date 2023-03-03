/*
Command pod-network-check sends pod connection failure events to DataDog.

Usage:

	pod-network-check -name {C} -namespace {N} -ports {X,Y,Z...}
	                  -dogstatsd {D} -interval {I} -timeout {T}
	                  -protocol {P} -debug

Arguments:

	-name        A name for the check; used the aggregation key.
	-namespace   The `kubernetes' namespace to check for pods.
	-ports       Comma separated list of ports to check.
	-dogstatsd   Address to the `dogstatsd' server.
	-interval    Check run interval, in seconds.
	-timeout     Dial timeout for pods, in seconds.
	-protocol    Protocol to use for the check.
	-debug       Print debug output.

Defaults:

	-name        pod-network-check
	-interval    30
	-timeout     5
	-protocol    tcp
	-debug       false
*/
package main
