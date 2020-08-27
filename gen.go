// +build ignore

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"path"

	"github.com/containous/yaegi/extract"
)

func main() {
	genPkg("github.com/grafana/tanka/pkg/kubernetes/manifest")
	genPkg("github.com/sh0rez/klint/pkg/klint")
}

var ex = extract.Extractor{
	Dest: "dynamic",
}

func genPkg(name string) {
	buf := bytes.Buffer{}
	if _, err := ex.Extract(name, "", &buf); err != nil {
		log.Fatalln(err)
	}

	filename := fmt.Sprintf("pkg/dynamic/yaegi_pkg_%s.go", path.Base(name))
	fmt.Println(filename)
	if err := ioutil.WriteFile(filename, buf.Bytes(), 0644); err != nil {
		log.Fatalln(err)
	}
}
