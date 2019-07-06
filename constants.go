package main

const (
	program = "pod-network-check"

	// Default check parameters.
	checkProtocol   = "tcp"
	checkInterval   = 30 // In seconds.
	checkPodTimeout = 5  // In seconds.

	// For `dogstatsd'.
	ddEvtTitle    = "Pod Connection Failures"
	ddEvtMkdnPre  = "%%% \n"
	ddEvtMkdnPost = "\n %%%"
	ddEvtMessage  = "Error reaching one or more pods: %d/%d.\nFailures:\n%s"
	ddAggKey      = "pod-network-check"
	ddSrcType     = "Kubernetes"
)
