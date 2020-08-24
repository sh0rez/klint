package klint

import "github.com/grafana/tanka/pkg/kubernetes/manifest"

// Rule is the function interface used by the individual rules
type Rule func(manifest.Manifest) (Findings, error)

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
}

type Findings []Finding

func Found(f ...Finding) Findings {
	return f
}
