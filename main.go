package main

//go:generate go run gen.go

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/containous/yaegi/interp"
	"github.com/containous/yaegi/stdlib"
	"github.com/go-clix/cli"
	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/sh0rez/klint/pkg/klint"
)

var Symbols = make(interp.Exports)

func main() {
	log.SetFlags(0)

	cmd := cli.Command{
		Use: "klint [.yaml, ...]",
	}

	cmd.Run = func(cmd *cli.Command, args []string) error {
		i := interp.New(interp.Options{
			GoPath:    "/tmp",
			BuildTags: []string{"klint"},
		})
		i.Use(stdlib.Symbols)
		i.Use(Symbols)

		// load rule.klint
		// TODO: dynamic rule loading
		ruleSrc, err := ioutil.ReadFile("rule.klint.go")
		if err != nil {
			return err
		}

		if _, err := i.Eval(string(ruleSrc)); err != nil {
			return err
		}

		rulePtr, err := i.Eval("lint")
		if err != nil {
			return err
		}

		rule, ok := rulePtr.Interface().(func(manifest.Manifest) (klint.Findings, error))
		if !ok {
			return fmt.Errorf("`lint` function is not of `klint.Rule` type but `%T`", rulePtr.Interface())
		}

		k := klint.Klint{
			Rules: klint.Rules{
				"rule.klint": rule,
			},
		}

		return k.LintFiles(args...)
	}

	if err := cmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
