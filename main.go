package main

//go:generate go run gen.go

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

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
	createNewRule := cmd.Flags().String("new", "", "Create a .klint.go file for a new rule")

	cmd.Run = func(cmd *cli.Command, args []string) error {
		// flag actions
		switch {
		case *createNewRule != "":
			return writeNewRule(*createNewRule)
		}

		// main action
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

func writeNewRule(name string) error {
	name = filepath.Clean(name + ".klint.go")

	if _, err := os.Stat(name); err == nil {
		return fmt.Errorf("File at '%s' already exists. Aborting", name)
	}

	const content = `// +build klint

package main

import (
	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/sh0rez/klint/pkg/klint"
)

func lint(m manifest.Manifest) (klint.Findings, error) {
	// Example: check if metadata.namespace is set
	//
	// if m.Metadata().Namespace() == "" {
	// 	f := klint.Finding{
	// 		Level:   klint.Warning,
	// 		Field:   "metadata.namespace",
	// 		Message: "Namespace is required for all resources",
	// 	}

	// 	return klint.Found(f), nil
	// }

	return nil, nil
}

func main() {
	klint.Main(lint)
}
`

	return ioutil.WriteFile(name, []byte(content), 0644)
}
