package klint

import (
	"fmt"
	"log"

	"github.com/go-clix/cli"
	"github.com/grafana/tanka/pkg/kubernetes/manifest"
)

// IS_YAEGI controls whether Main runs. We don't want that in interpreted mode,
// but do what it in "native Go" mode
var IS_YAEGI = false

// Main provides klint-single-rule, used for easy execution of single rules
// during development, without the need for an interpreter.
// It allows to `go run` a .klint.go file
func Main(f func(manifest.Manifest) (Findings, error)) {
	log.SetFlags(0)

	if IS_YAEGI {
		return
	}

	r := Rule(f)

	cmd := cli.Command{
		Use:   "klint-single-rule [.yaml, ...]",
		Short: "klint-single-rule runs a single rule",
	}

	cmd.Run = func(cmd *cli.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("Please specify at least one .yaml file containing Kubernetes resources")
		}

		k := Klint{
			Rules: Rules{"": r},
		}

		return k.LintFiles(args...)
	}

	if err := cmd.Execute(); err != nil {
		log.Fatalln(err)
	}

	log.Println("All resources pass.")
}
