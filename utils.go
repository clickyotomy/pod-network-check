package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/davecgh/go-spew/spew"
)

// contains is a helper function to check if a slice has a member.
func contains(port uint16, ports []uint16) bool {
	for _, p := range ports {
		if p == port {
			return true
		}
	}
	return false
}

// alive checks if the pod is reachable.
func alive(sock podSock, proto string, ch chan status) {
	var (
		conn net.Conn
		err  error
	)

	// IP address check.
	if net.ParseIP(sock.addr) == nil {
		ch <- status{
			sock.name,
			fmt.Errorf("ip-parse: bad address: %s", sock.addr),
			false,
		}
	} else {
		conn, err = net.DialTimeout(
			proto,
			fmt.Sprintf("%s:%d", sock.addr, sock.port),
			sock.timeout,
		)
		if err != nil {
			ch <- status{sock.name, fmt.Errorf("fail: %s", err), false}
		} else {
			ch <- status{sock.name, nil, true}
			defer conn.Close()
		}
	}
}

// event sends an event to `dogstatsd'.
func event(client *statsd.Client, checks []status, ns, name string) error {
	var (
		snd uint8
		bad []status
		err error
		evt *statsd.Event
		tmp string
	)

	for _, chk := range checks {
		if !chk.pass {
			snd++
			bad = append(bad, chk)
		}
	}

	if printDbg {
		log.Printf("pods: %s\n", spew.Sdump(checks))
	}

	if snd > 0 {
		log.Printf(
			"fail: Unable to reach one or more pod(s): %d/%d.",
			snd, len(checks),
		)

		for _, s := range bad {
			tmp += fmt.Sprintf("* pod: %s\n    %s\n", s.pod, s.err)
		}

		tmp = "```\n" + tmp + "\n```"
		evt = &statsd.Event{
			Title: ddEvtTitle,
			Text: strings.Join(
				[]string{
					ddEvtMkdnPre,
					fmt.Sprintf(ddEvtMessage, snd, len(checks), tmp),
					ddEvtMkdnPost,
				}, " ",
			),
			AggregationKey: ddAggKey,
			Priority:       statsd.Normal,
			AlertType:      statsd.Error,
			SourceTypeName: ddSrcType,
			Tags: []string{
				fmt.Sprintf("check:%s", name),
				fmt.Sprintf("check-namespace:%s", ns),
			},
		}

		log.Printf("send: Pushing event to DataDog.")
		err = client.Event(evt)

		if printDbg {
			log.Printf("send: %s\n", spew.Sdump(*evt))
		}
	} else {
		log.Printf(
			"pass: All pods reachable: %d/%d.",
			(len(checks) - int(snd)), len(checks),
		)
	}

	return err
}
