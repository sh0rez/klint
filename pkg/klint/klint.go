package klint

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"gopkg.in/yaml.v3"
)

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

// Klint lints Kubernetes manifests using runtime-loaded rules.
type Klint struct {
	Rules Rules
}

// LintFiles lints all Manifests in the given files using the rules inside the
// Klint instance
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

// Lint all given manifests using the rules in the Klint instance
func (k Klint) Lint(list manifest.List) error {
	result := make(Result)

	for _, m := range list {
		f, err := k.Rules.run(m)
		if err != nil {
			return err
		}

		if f == nil {
			continue
		}

		result[m.KindName()] = append(result[m.KindName()], f...)
	}

	if len(result) == 0 {
		return nil
	}

	return result
}

// Result is returned by Lint. It maps resource names to their respective
// findings.
type Result map[string]Findings

func (r Result) Error() string {
	s := ""
	for kindName, findings := range r {
		if len(findings) == 0 {
			continue
		}

		s += kindName + ":\n"
		for _, f := range findings {
			s += " " + f.String() + "\n"
		}
	}

	return strings.TrimSpace(s)
}
