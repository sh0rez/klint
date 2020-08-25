package klint

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/go-clix/cli"
	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"gopkg.in/yaml.v3"
)

func Main(f func(manifest.Manifest) (Findings, error)) {
	r := Rule(f)

	cmd := cli.Command{
		Use:   "klint-rule [.yaml, ...]",
		Short: "klint-rule runs a single rule",
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
		log.SetFlags(0)
		log.Fatalln(err)
	}
}

// Rule is the function interface used by the individual rules
type Rule func(manifest.Manifest) (Findings, error)

type Rules map[string]Rule

func (r Rules) run(m manifest.Manifest) (Findings, error) {
	var f Findings
	for name, rule := range r {
		found, err := rule(m)
		if err != nil {
			return nil, err
		}

		for i := range found {
			found[i].rule = name
		}
		f = append(f, found...)
	}

	return f, nil
}

// Level is the severity of a Finding
type Level int

func (l Level) String() string {
	return Levels[l]
}

const (
	Info Level = iota
	Warning
	Error
)

var Levels = map[Level]string{
	Info:    "INFO",
	Warning: "WARN",
	Error:   "ERROR",
}

// Finding describes a single thing the Rule found required to be mentioned
type Finding struct {
	// Level of the finding
	Level Level

	// Field that is offending
	Field string

	// Message in human readable format what happened, and possibly how to fix it
	Message string

	// internal fields
	rule string
}

func (f Finding) String() string {
	name := ""
	if f.rule != "" {
		name = " – " + f.rule
	}

	return fmt.Sprintf("%s%s – %s: %s", f.Level, name, f.Field, f.Message)
}

type Findings []Finding

func Found(f ...Finding) Findings {
	return f
}

type Klint struct {
	Rules Rules
}

func (k Klint) LintFiles(files ...string) error {
	var list manifest.List

	for _, f := range files {
		file, err := os.Open(f)
		if err != nil {
			return err
		}

		d := yaml.NewDecoder(file)
		for {
			var m manifest.Manifest
			err := d.Decode(&m)
			if err == io.EOF {
				break
			} else if err != nil {
				return err
			}

			list = append(list, m)
		}
	}

	if len(list) == 0 {
		return fmt.Errorf("No resources found. Aborting")
	}

	return k.Lint(list)
}

func (k Klint) Lint(list manifest.List) error {
	result := make(Result)

	for _, m := range list {
		f, err := k.Rules.run(m)
		if err != nil {
			return err
		}

		result[m.KindName()] = append(result[m.KindName()], f...)
	}

	if len(result) == 0 {
		return nil
	}

	return result
}

type Result map[string]Findings

func (r Result) Error() string {
	s := ""
	for kindName, findings := range r {
		s += kindName + ":\n"
		for _, f := range findings {
			s += " " + f.String() + "\n"
		}
	}

	return strings.TrimSpace(s)
}
