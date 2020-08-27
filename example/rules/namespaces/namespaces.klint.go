package main

import (
	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/sh0rez/klint/pkg/klint"
)

func lint(m manifest.Manifest) (klint.Findings, error) {
	if m.Metadata().Namespace() == "" {
		f := klint.Finding{
			Level:   klint.Warning,
			Field:   "metadata.namespace",
			Message: "Namespace is required for all resources",
		}

		return klint.Found(f), nil
	}

	return nil, nil
}

func main() {
	klint.Main(lint)
}
