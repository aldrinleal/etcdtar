package main

import (
  "flag"
  "github.com/coreos/go-etcd/etcd"
  "github.com/aldrinleal/etcdtar"
)

func main() {
        var mode, host string

        flag.StringVar(&mode, "mode", "c", "c (etcd -> tar) | x (tar -> etcd)")
        flag.StringVar(&host, "host", "http://127.0.0.1:4001", "host to use")

        flag.Parse()

        client := etcd.NewClient([]string{host})

        if mode == "c" {
                etcdtar.ExportFromEtcdToTar(client, host)
        } else if mode == "x" {
                etcdtar.ExportFromTarToEtcd(client, host)
        }
}

