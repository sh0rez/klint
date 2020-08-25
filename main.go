package main

//go:generate go run gen.go

import (
	"fmt"
	"log"

	"github.com/go-clix/cli"
	"github.com/sh0rez/klint/pkg/dynamic"
	"github.com/sh0rez/klint/pkg/klint"
)

func main() {
	log.SetFlags(0)

	cmd := cli.Command{
		Use:   "klint [.yaml, ...]",
		Short: "klint validates Kubernetes configurations using custom rules written in Golang",
	}

	ruleFiles := cmd.Flags().StringSliceP("rules", "R", nil, ".klint.go files containing rules to load")

	cmd.Run = func(cmd *cli.Command, args []string) error {
		if len(*ruleFiles) == 0 {
			return fmt.Errorf("Please pass at least one rule using --rules / -R")
		}

		rules, err := dynamic.LoadFiles(*ruleFiles)
		if err != nil {
			return err
		}

		k := klint.Klint{
			Rules: rules,
		}

		return k.LintFiles(args...)
	}

	if err := cmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
