package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// For debug output.
var printDbg bool

func main() {
	var (
		// For flags.
		interval  = flag.Uint("interval", checkInterval, "check interval")
		timeout   = flag.Uint("timeout", checkPodTimeout, "pod dial timeout")
		debug     = flag.Bool("debug", false, "debug output")
		name      = flag.String("name", ddAggKey, "check name")
		protocol  = flag.String("protocol", checkProtocol, "check protocol")
		namespace = flag.String("namespace", "", "kubernetes namespace")
		chkPorts  = flag.String("ports", "", "comma separated list of ports")
		ddogAddr  = flag.String("dogstatsd", "", "address for dogstatsd")

		// For `kubernetes'.
		conf  *rest.Config
		cset  *kubernetes.Clientset
		pods  *corev1.PodList
		ports []uint16

		// For `dogstatsd'.
		ddog *statsd.Client

		// For tracking go-routines.
		ch = make(chan status)
		gc int

		// Misc.
		err error
		tmp int
		dbg *bool
		atn annotation
		chk []status
		tkr *time.Ticker
	)

	flag.Usage = usage
	flag.Parse()

	if *namespace == "" {
		fmt.Printf("flag: cannot run check without a namespace\n")
		os.Exit(1)
	}

	if *chkPorts == "" {
		fmt.Printf("flag: must specify at least one port\n")
		os.Exit(1)
	}

	if *ddogAddr == "" {
		fmt.Printf("flag: must specify an address for dogstatsd\n")
		os.Exit(1)
	}

	// For debug printing.
	dbg = &printDbg
	*dbg = *debug

	for _, p := range strings.Split(*chkPorts, ",") {
		if tmp, err = strconv.Atoi(p); err != nil {
			log.Panic(fmt.Errorf("flag: bad port %s", p))
		}

		ports = append(ports, uint16(tmp))
	}

	// Log arguments.
	log.Printf(
		"args: "+
			"namespace: %s, ports: %v, interval: %d, timeout: %d, protocol: %s",
		*namespace, ports, *interval, *timeout, *protocol,
	)

	// Create an in-cluster config.
	if conf, err = rest.InClusterConfig(); err != nil {
		log.Panic(err.Error())
	}

	// Create a client-set.
	if cset, err = kubernetes.NewForConfig(conf); err != nil {
		log.Panic(err.Error())
	}

	tkr = time.NewTicker(time.Second * time.Duration(*interval))

	if ddog, err = statsd.New(*ddogAddr); err != nil {
		log.Panic(err)
	}

	// Loop for every interval.
	for {
		select {
		case <-tkr.C:
			log.Printf("kube: Starting checks...")

			// Query for the pods.
			pods, err = cset.CoreV1().Pods(*namespace).List(
				metav1.ListOptions{},
			)
			if err != nil {
				log.Panic(err.Error())
			}

			for _, pod := range pods.Items {
				if pod.ObjectMeta.Annotations != nil {
					for _, value := range pod.ObjectMeta.Annotations {
						err = json.Unmarshal([]byte(value), &atn)
						if err != nil {
							log.Printf("json: %s", err)
							continue
						}

						// Run checks.
						for _, node := range atn.Nodes {
							if contains(node.Port, ports) {
								go alive(
									podSock{
										pod.ObjectMeta.Name,
										node.Address,
										node.Port,
										time.Second * time.Duration(*timeout),
									},
									*protocol,
									ch,
								)
								gc++
							}
						}
					}
				}
			}

			for gc > 0 {
				if c, ok := <-ch; ok {
					chk = append(chk, c)
					gc--
				}
			}

			// We don't want to `panic' here, because it could be a temporary
			// network issue, or something similar.
			if err = event(ddog, chk, *namespace, *name); err != nil {
				log.Printf("ddog: %s", err)
			}

			// Reset checks and counters.
			chk = []status{}
			gc = 0
		}
	}
}
