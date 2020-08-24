package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/containous/yaegi/interp"
	"github.com/containous/yaegi/stdlib"
	"github.com/fatih/color"
	"github.com/go-clix/cli"
	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/pkg/errors"
	"github.com/sh0rez/klint/pkg/klint"
	"gopkg.in/yaml.v3"
)

var Symbols = make(interp.Exports)

func main() {
	log.SetFlags(0)

	cmd := cli.Command{
		Use: "klint [.yaml, ...]",
	}

	cmd.Run = func(cmd *cli.Command, args []string) error {
		list, err := loadManifests(args)
		if err != nil {
			return err
		}

		if len(list) == 0 {
			return fmt.Errorf("No resources found. Please speficy one or more YAML files that contain Kubernetes resources")
		}

		i, err := setupInterp()
		if err != nil {
			return err
		}

		// load rule.klint
		// TODO: dynamic rule loading
		ruleSrc, err := ioutil.ReadFile("rule.klint")
		if err != nil {
			return err
		}

		if _, err := i.Eval(string(ruleSrc)); err != nil {
			return errors.Wrap(err, "Parsing source")
		}

		rulePtr, err := i.Eval("lint")
		if err != nil {
			return errors.Wrap(err, "Getting `lint` function from rule source")
		}

		rule, ok := rulePtr.Interface().(func(manifest.Manifest) (klint.Findings, error))
		if !ok {
			return fmt.Errorf("`lint` function is not of `klint.Rule` type but `%T`", rulePtr.Interface())
		}

		// run rule
		for _, m := range list {
			findings, err := rule(m)
			if err != nil {
				return err
			}

			if len(findings) == 0 {
				continue
			}

			log.Println(color.YellowString(m.KindName() + ":"))
			for _, f := range findings {
				log.Printf(" %s – %s – %s: %s", f.Level, "rule.klint", f.Field, f.Message)
			}
		}

		return nil
	}

	if err := cmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func loadManifests(files []string) (manifest.List, error) {
	var list manifest.List

	for _, f := range files {
		file, err := os.Open(f)
		if err != nil {
			return nil, err
		}

		d := yaml.NewDecoder(file)
		for {
			var m manifest.Manifest
			err := d.Decode(&m)
			if err == io.EOF {
				break
			} else if err != nil {
				return nil, err
			}

			list = append(list, m)
		}
	}

	return list, nil
}

func setupInterp() (*interp.Interpreter, error) {
	i := interp.New(interp.Options{
		GoPath: "/tmp",
	})
	i.Use(stdlib.Symbols)
	i.Use(Symbols)

	// 	const base = `
	// package rule

	// import (
	//   "fmt"
	//   "github.com/grafana/tanka/pkg/kubernetes/manifest"
	// )
	// `
	// if _, err := i.Eval(base); err != nil {
	// 	return nil, err
	// }

	return i, nil
}
