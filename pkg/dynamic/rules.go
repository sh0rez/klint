package dynamic

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/containous/yaegi/interp"
	"github.com/containous/yaegi/stdlib"
	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/sh0rez/klint/pkg/klint"
)

// Symbolds hold pre-compiled libraries that are exposed to the rules (see
// yaegi_pkg.* files)
// - github.com/grafana/tanka/pkg/manifest
// - github.com/sh0rez/klint/pkg/klint
var Symbols = make(interp.Exports)

// LoadFiles loads each file as a rule
func LoadFiles(files []string) (klint.Rules, error) {
	rules := make(klint.Rules)

	for _, f := range files {
		data, err := ioutil.ReadFile(f)
		if err != nil {
			return nil, err
		}

		rule, err := Load(data)
		if err != nil {
			return nil, err
		}

		name := filepath.Base(f)
		name = strings.TrimSuffix(name, ".klint.go")
		name = strings.TrimSuffix(name, ".go")

		rules[name] = rule
	}

	return rules, nil
}

// Load uses Yaegi to interpret a .klint.go file at runtime. The `lint` function
// is accessed from the interpreted result and returned as a regular
// `klint.Rule`
func Load(src []byte) (klint.Rule, error) {
	klint.IS_YAEGI = true
	defer func() {
		klint.IS_YAEGI = false
	}()

	i := interp.New(interp.Options{
		GoPath:    "/tmp/klint",
		BuildTags: []string{"klint"},
	})
	i.Use(stdlib.Symbols)
	i.Use(Symbols)

	if _, err := i.Eval(string(src)); err != nil {
		return nil, fmt.Errorf("Interpreting Go source: %w", err)
	}

	rulePtr, err := i.Eval("lint")
	if err != nil {
		return nil, fmt.Errorf("Accessing 'lint' function of interpreted rule: %w", err)
	}

	rule, ok := rulePtr.Interface().(func(manifest.Manifest) (klint.Findings, error))
	if !ok {
		return nil, fmt.Errorf("`lint` function is not of `klint.Rule` type but `%T`", rulePtr.Interface())
	}

	return rule, nil
}
