package conflict

import (
	"fmt"
	"strings"
)

// ConflictType identifies the type of conflict.
type ConflictType int

const (
	ConflictPort ConflictType = iota
	ConflictContainerName
	ConflictNetwork
	ConflictVolume
	ConflictDisabled
)

// String returns the string representation of ConflictType.
func (t ConflictType) String() string {
	switch t {
	case ConflictPort:
		return "port"
	case ConflictContainerName:
		return "container_name"
	case ConflictNetwork:
		return "network"
	case ConflictVolume:
		return "volume"
	case ConflictDisabled:
		return "disabled"
	default:
		return "unknown"
	}
}

// ConflictSeverity indicates the severity of a conflict.
type ConflictSeverity int

const (
	SeverityError ConflictSeverity = iota
	SeverityWarning
)

// String returns the string representation of ConflictSeverity.
func (s ConflictSeverity) String() string {
	switch s {
	case SeverityError:
		return "error"
	case SeverityWarning:
		return "warning"
	default:
		return "unknown"
	}
}

// Conflict represents a detected conflict between services.
type Conflict struct {
	Type      ConflictType
	Severity  ConflictSeverity
	ServiceA  string // first service involved
	ServiceB  string // second service involved (empty for host conflicts)
	Resource  string // the conflicting resource (port number, container name, etc.)
	Message   string // human-readable description
	Suggested string // suggested fix (e.g., alternative port)
}

// Error implements the error interface.
func (c Conflict) Error() string {
	return c.Message
}

// HasSuggestion returns true if a suggestion is available.
func (c Conflict) HasSuggestion() bool {
	return c.Suggested != ""
}

// ConflictList is a slice of conflicts with helper methods.
type ConflictList []Conflict

// HasErrors returns true if any conflict is severity error.
func (cl ConflictList) HasErrors() bool {
	for _, c := range cl {
		if c.Severity == SeverityError {
			return true
		}
	}
	return false
}

// HasWarnings returns true if any conflict is severity warning.
func (cl ConflictList) HasWarnings() bool {
	for _, c := range cl {
		if c.Severity == SeverityWarning {
			return true
		}
	}
	return false
}

// Errors returns only error-severity conflicts.
func (cl ConflictList) Errors() ConflictList {
	var result ConflictList
	for _, c := range cl {
		if c.Severity == SeverityError {
			result = append(result, c)
		}
	}
	return result
}

// Warnings returns only warning-severity conflicts.
func (cl ConflictList) Warnings() ConflictList {
	var result ConflictList
	for _, c := range cl {
		if c.Severity == SeverityWarning {
			result = append(result, c)
		}
	}
	return result
}

// ByType groups conflicts by type.
func (cl ConflictList) ByType() map[ConflictType]ConflictList {
	result := make(map[ConflictType]ConflictList)
	for _, c := range cl {
		result[c.Type] = append(result[c.Type], c)
	}
	return result
}

// Messages returns all conflict messages as a formatted string.
func (cl ConflictList) Messages() string {
	if len(cl) == 0 {
		return ""
	}
	msgs := make([]string, len(cl))
	for i, c := range cl {
		msgs[i] = fmt.Sprintf("[%s] %s", c.Severity.String(), c.Message)
	}
	return strings.Join(msgs, "\n")
}
