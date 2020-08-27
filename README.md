# Klint

Klint lints Kubernetes resources, using an arbitrary ruleset.

> **Pre-Alpha**: The development of Klint has just started and is far from
> finished. It probably does not work.

It does not ship with rules itself (yet), however these are very easy to write, in Golang:

```go
// filename: rule.klint
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
```

The key part here is the `lint` function, which must follow exactly above
function signature. It receives a single Kubernetes manifest
(`map[string]interface{}` in a helper type) and may return any number of
findings on that resource.

Rules are dynamically loaded at runtime, and produce such a result:

```bash
$ klint deployment.yaml

Deployment/nginx-deployment:
 WARN – rule.klint – metadata.namespace: Namespace is required for all resources
```

## Example

For an example, check the
[`example`](https://github.com/sh0rez/klint/tree/master/example) directory.
