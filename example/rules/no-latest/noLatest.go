package main

import (
	"fmt"
	"strings"

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

		image := ct["image"].(string)

		if strings.HasSuffix(image, ":latest") {
			return klint.Found(klint.Finding{
				Level:   klint.Warning,
				Field:   fmt.Sprintf("spec.template.spec.containers[%v].image", i),
				Message: "Latest tag should not be used",
			}), nil
		}
	}

	return nil, nil
}

func main() {
	klint.Main(lint)
}
