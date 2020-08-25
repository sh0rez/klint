// +build ignore

package main

import (
	"bytes"
	"io/ioutil"
	"log"

	"github.com/containous/yaegi/extract"
)

func main() {
	ex := extract.Extractor{
		Dest: "dynamic",
	}

	// github.com/grafana/tanka/pkg/kubernetes/manifest
	var buf bytes.Buffer
	if _, err := ex.Extract("github.com/grafana/tanka/pkg/kubernetes/manifest", "", &buf); err != nil {
		log.Fatalln(err)
	}
	if err := ioutil.WriteFile("pkg/dynamic/yaegi_pkg_manifest.go", buf.Bytes(), 0644); err != nil {
		log.Fatalln(err)
	}

	// github.com/sh0rez/klint/pkg/klint
	buf = bytes.Buffer{}
	if _, err := ex.Extract("github.com/sh0rez/klint/pkg/klint", "", &buf); err != nil {
		log.Fatalln(err)
	}
	if err := ioutil.WriteFile("pkg/dynamic/yaegi_pkg_klint.go", buf.Bytes(), 0644); err != nil {
		log.Fatalln(err)
	}
}
