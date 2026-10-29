package backlog

import (
	"fmt"
	"regexp"
)

// taskIDPattern is deliberately narrow. Beans generates ids like "task-a1b2" and the
// workshop seeds ids like "task-101"; nothing else is accepted. In particular an id
// can never begin with '-', so it can never be mistaken for a CLI flag, and it can
// never contain a path separator or shell metacharacter.
var taskIDPattern = regexp.MustCompile(`^[a-z0-9]{1,24}(-[a-z0-9]{1,24}){1,3}$`)

const maxTaskIDLen = 64

// ValidateTaskID returns a typed error for anything outside the accepted syntax.
func ValidateTaskID(id string) error {
	if id == "" {
		return newError(CodeInvalidID, "task id is required", "")
	}
	if len(id) > maxTaskIDLen {
		return newError(CodeInvalidID,
			fmt.Sprintf("task id is longer than %d characters", maxTaskIDLen), "")
	}
	if !taskIDPattern.MatchString(id) {
		return newError(CodeInvalidID,
			"task id must look like task-101: lowercase letters and digits in 2 to 4 dash-separated groups", "")
	}
	return nil
}

// allowedStatuses mirrors the Beans status vocabulary. The adapter refuses to pass
// through a status it does not know, so a client cannot smuggle a flag value.
var allowedStatuses = map[string]bool{
	"todo":        true,
	"in-progress": true,
	"draft":       true,
	"completed":   true,
	"scrapped":    true,
}

var tagPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,31}$`)

const (
	// DefaultListLimit keeps an unfiltered list small enough for an agent context.
	DefaultListLimit = 20
	// MaxListLimit bounds list_tasks regardless of what the client asks for.
	MaxListLimit = 50
)

// ValidateListFilter normalises and bounds a list filter.
func ValidateListFilter(f ListFilter) (ListFilter, error) {
	out := ListFilter{Limit: f.Limit}
	if f.Status != "" {
		if !allowedStatuses[f.Status] {
			return out, newError(CodeInvalidFilter,
				"status must be one of: todo, in-progress, draft, completed, scrapped", "")
		}
		out.Status = f.Status
	}
	if f.Tag != "" {
		if !tagPattern.MatchString(f.Tag) {
			return out, newError(CodeInvalidFilter,
				"tag must be lowercase letters, digits and dashes, at most 32 characters", "")
		}
		out.Tag = f.Tag
	}
	switch {
	case out.Limit == 0:
		out.Limit = DefaultListLimit
	case out.Limit < 0:
		return out, newError(CodeInvalidFilter, "limit must be a positive integer", "")
	case out.Limit > MaxListLimit:
		return out, newError(CodeInvalidFilter,
			fmt.Sprintf("limit must be %d or less", MaxListLimit), "")
	}
	return out, nil
}
