package klint

import "fmt"

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

// Found constructs Findings from a set of Finding. This is required because the
// Yaegi interpreter does not support the Findings{} slice literal
func Found(f ...Finding) Findings {
	return f
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

// Levels provides human readable values for the Level enu,
var Levels = map[Level]string{
	Info:    "INFO",
	Warning: "WARN",
	Error:   "ERROR",
}
