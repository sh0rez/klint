package main

import (
	"fmt"

	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/sh0rez/klint/pkg/klint"
)

func lint(m manifest.Manifest) (klint.Findings, error) {
	got, err := klint.GetField(m, "spec", "template", "spec", "containers")
	if err != nil {
		return nil, err
	}

	for i, v := range got.Slice() {
		ct := v.(map[string]interface{})

		if _, ok := ct["resources"]; !ok {
			return klint.Found(klint.Finding{
				Level:   klint.Warning,
				Field:   fmt.Sprintf("spec.template.spec.containers[%v].resources", i),
				Message: "Resource requests / limits should be set",
			}), nil
		}
	}

	return nil, nil
}

func main() {
	klint.Main(lint)
}
