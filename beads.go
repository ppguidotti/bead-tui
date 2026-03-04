package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

// Issue represents a Beads task as returned by `bd list --json`
type Issue struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Status      string   `json:"status"`
	Priority    int      `json:"priority"`
	Type        string   `json:"type"`
	Assignee    string   `json:"assignee,omitempty"`
	Labels      []string `json:"labels,omitempty"`
	Description string   `json:"description,omitempty"`
	Notes       string   `json:"notes,omitempty"`
	CreatedAt   string   `json:"created_at,omitempty"`
	UpdatedAt   string   `json:"updated_at,omitempty"`
}

// IssueDetail is returned by `bd show <id> --json` and includes dependencies
type IssueDetail struct {
	Issue
	Dependencies []Dependency `json:"dependencies,omitempty"`
	Blockers     []string     `json:"blockers,omitempty"`
}

type Dependency struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// ListIssues runs `bd list --json` and returns all issues.
// If filter is non-empty (e.g. "open", "in_progress"), it passes --status.
func ListIssues(filter string) ([]Issue, error) {
	args := []string{"list", "--json"}
	if filter != "" && filter != "all" {
		args = append(args, "--status", filter)
	}

	out, err := runBD(args...)
	if err != nil {
		return nil, err
	}

	var issues []Issue
	if err := json.Unmarshal(out, &issues); err != nil {
		return nil, fmt.Errorf("failed to parse issues: %w", err)
	}
	return issues, nil
}

// ShowIssue runs `bd show <id> --json` and returns full detail.
func ShowIssue(id string) (*IssueDetail, error) {
	out, err := runBD("show", id, "--json")
	if err != nil {
		return nil, err
	}

	var detail IssueDetail
	if err := json.Unmarshal(out, &detail); err != nil {
		return nil, fmt.Errorf("failed to parse issue detail: %w", err)
	}
	return &detail, nil
}

// ReadyIssues runs `bd ready --json` — tasks with no open blockers.
func ReadyIssues() ([]Issue, error) {
	out, err := runBD("ready", "--json")
	if err != nil {
		return nil, err
	}

	var issues []Issue
	if err := json.Unmarshal(out, &issues); err != nil {
		return nil, fmt.Errorf("failed to parse ready issues: %w", err)
	}
	return issues, nil
}

// runBD executes the `bd` binary with the given arguments.
func runBD(args ...string) ([]byte, error) {
	cmd := exec.Command("bd", args...)
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("bd %v failed: %s", args, exitErr.Stderr)
		}
		return nil, fmt.Errorf("bd not found — is it installed? (run: npm install -g @beads/bd): %w", err)
	}
	return out, nil
}

// PriorityLabel maps numeric priority to a human-readable label.
func PriorityLabel(p int) string {
	switch p {
	case 0:
		return "P0"
	case 1:
		return "P1"
	case 2:
		return "P2"
	case 3:
		return "P3"
	default:
		return fmt.Sprintf("P%d", p)
	}
}
