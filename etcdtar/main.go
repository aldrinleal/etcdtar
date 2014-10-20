package main

import (
	"flag"
	"log"
	"os"

	"github.com/aldrinleal/etcdtar"
	"github.com/coreos/go-etcd/etcd"
)

var (
	mode = flag.String("mode", "c", "c (etcd -> tar) | x (tar -> etcd)")
	host = flag.String("host", "http://127.0.0.1:4001", "host to use")
	path = flag.String("path", "/", "path to export from etcd")

	weakConsistency = flag.Bool("weak-consistency", false, "Use weak consistency checks")
	debugEtcd       = flag.Bool("debug-etcd", false, "Enable debug output for go-etcd")
	showCurl        = flag.Bool("show-curl", false, "Show CURL commands.")

	errLogger = log.New(os.Stderr, "", log.LstdFlags)
)

func main() {
	flag.Parse()

	client := etcdClient()
	if *mode == "c" {
		etcdtar.ExportFromEtcdToTar(client, *host, *path)
	} else if *mode == "x" {
		etcdtar.ExportFromTarToEtcd(client, *host)
	} else {
		errLogger.Fatalf("Unknown mode: %s", *mode)
	}
}

func etcdClient() *etcd.Client {
	client := etcd.NewClient([]string{*host})

	if *debugEtcd {
		etcd.SetLogger(log.New(os.Stdout, "[etcd] ", log.LstdFlags))
	}

	if *showCurl {
		client.OpenCURL()
		defer client.CloseCURL()

		logger := log.New(os.Stderr, "[curl] ", log.LstdFlags)

		go func() {
			for {
				cmd := client.RecvCURL()
				logger.Printf("CURL: %v\n", cmd)
			}
		}()
	}

	if *weakConsistency {
		// Use this if you have a local test instance.
		client.SetConsistency(etcd.WEAK_CONSISTENCY)
	}
	return client
}
