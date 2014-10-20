package etcdtar

import (
	"archive/tar"
	"bufio"
	"fmt"
	"github.com/coreos/go-etcd/etcd"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func ExportFromEtcdToTar(client *etcd.Client, host, path string) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	dirs := make(map[string]bool)
	dirs[path] = true

outer:
	for {

		for dir, _ := range dirs {
			r, err := client.Get(dir, false, true)

			if nil != err {
				panic(err)
			}

			for _, n := range r.Node.Nodes {
				if !n.Dir {
					continue
				}

				if !dirs[n.Key] {
					dirs[n.Key] = true
					continue outer
				}
			}
		}

		break
	}

	content := make(map[string]string)

	for dir, _ := range dirs {
		r, err := client.Get(dir, false, true)

		if nil != err {
			panic(err)
		}

		for _, n := range r.Node.Nodes {
			if n.Dir {
				continue
			}

			content[n.Key[1:]] = n.Value
		}
	}

	buf := bufio.NewWriter(os.Stdout)

	defer buf.Flush()

	tw := tar.NewWriter(buf)

	for path, content := range content {
		hdr := &tar.Header{
			Name: path,
			Mode: 0640,
			Size: int64(len(content)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			panic(err)
		}
		if _, err := tw.Write([]byte(content)); err != nil {
			panic(err)
		}
	}

	if err := tw.Close(); err != nil {
		panic(err)
	}

}

func ExportFromTarToEtcd(client *etcd.Client, host string) {
	buf := bufio.NewReader(os.Stdin)

	tr := tar.NewReader(buf)

	for {
		hdr, err := tr.Next()

		if err == io.EOF {
			// end of tar archive
			break
		}
		if nil != err {
			panic(err)
		}

		byteValue, err := ioutil.ReadAll(tr)

		if nil != err {
			panic(err)
		}

		strValue := string(byteValue)

		fmt.Printf("k: %s; v: %s\n", hdr.Name, strValue)

		r, err := client.Set(hdr.Name, strValue, 0)

		fmt.Println(r, err)
	}
}
