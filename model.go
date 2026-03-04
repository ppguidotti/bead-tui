package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// View modes
type viewMode int

const (
	viewList   viewMode = iota // main list
	viewDetail                 // detail panel for selected issue
)

// Filter tabs
var filterOptions = []string{"all", "open", "in_progress", "blocked", "closed"}

// Messages for async operations
type issuesLoadedMsg struct{ issues []Issue }
type detailLoadedMsg struct{ detail *IssueDetail }
type errMsg struct{ err error }

// Model holds all TUI state
type Model struct {
	mode         viewMode
	issues       []Issue
	cursor       int
	filterIdx    int // index into filterOptions
	detail       *IssueDetail
	loading      bool
	spinner      spinner.Model
	err          error
	windowWidth  int
	windowHeight int
}

func initialModel() Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#5C6BC0"))
	return Model{
		loading: true,
		spinner: s,
	}
}

// Init kicks off the initial data load.
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		loadIssues("all"),
	)
}

// --- Commands (async IO) ---

func loadIssues(filter string) tea.Cmd {
	return func() tea.Msg {
		issues, err := ListIssues(filter)
		if err != nil {
			return errMsg{err}
		}
		return issuesLoadedMsg{issues}
	}
}

func loadDetail(id string) tea.Cmd {
	return func() tea.Msg {
		detail, err := ShowIssue(id)
		if err != nil {
			return errMsg{err}
		}
		return detailLoadedMsg{detail}
	}
}

// --- Update ---

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height

	case tea.KeyMsg:
		return m.handleKey(msg)

	case issuesLoadedMsg:
		m.issues = msg.issues
		m.loading = false
		m.cursor = 0
		m.err = nil

	case detailLoadedMsg:
		m.detail = msg.detail
		m.loading = false

	case errMsg:
		m.err = msg.err
		m.loading = false

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {

	case "ctrl+c", "q":
		if m.mode == viewDetail {
			m.mode = viewList
			m.detail = nil
			return m, nil
		}
		return m, tea.Quit

	case "esc":
		if m.mode == viewDetail {
			m.mode = viewList
			m.detail = nil
		}

	case "up", "k":
		if m.mode == viewList && m.cursor > 0 {
			m.cursor--
		}

	case "down", "j":
		if m.mode == viewList && m.cursor < len(m.issues)-1 {
			m.cursor++
		}

	case "enter":
		if m.mode == viewList && len(m.issues) > 0 {
			m.mode = viewDetail
			m.loading = true
			return m, tea.Batch(m.spinner.Tick, loadDetail(m.issues[m.cursor].ID))
		}

	case "tab", "right", "l":
		if m.mode == viewList {
			m.filterIdx = (m.filterIdx + 1) % len(filterOptions)
			m.loading = true
			m.cursor = 0
			return m, tea.Batch(m.spinner.Tick, loadIssues(filterOptions[m.filterIdx]))
		}

	case "shift+tab", "left", "h":
		if m.mode == viewList {
			m.filterIdx = (m.filterIdx - 1 + len(filterOptions)) % len(filterOptions)
			m.loading = true
			m.cursor = 0
			return m, tea.Batch(m.spinner.Tick, loadIssues(filterOptions[m.filterIdx]))
		}

	case "r":
		// Refresh
		m.loading = true
		return m, tea.Batch(m.spinner.Tick, loadIssues(filterOptions[m.filterIdx]))
	}

	return m, nil
}

// --- View ---

func (m Model) View() string {
	if m.err != nil {
		return appStyle.Render(
			titleBarStyle.Render("⚡ bd-tui — Beads Task Viewer") + "\n\n" +
				errorStyle.Render("Error: "+m.err.Error()) + "\n\n" +
				helpStyle.Render("Make sure `bd` is installed and you're in a Beads project directory.\nPress q to quit."),
		)
	}

	header := titleBarStyle.Width(m.windowWidth - 4).Render("⚡ bd-tui — Beads Task Viewer")

	if m.loading {
		return appStyle.Render(header + "\n\n" + m.spinner.View() + " Loading...")
	}

	if m.mode == viewDetail && m.detail != nil {
		return appStyle.Render(header + "\n" + m.renderDetail())
	}

	return appStyle.Render(header + "\n" + m.renderList())
}

func (m Model) renderFilterTabs() string {
	tabs := make([]string, len(filterOptions))
	for i, f := range filterOptions {
		if i == m.filterIdx {
			tabs[i] = activeFilterStyle.Render(f)
		} else {
			tabs[i] = inactiveFilterStyle.Render(f)
		}
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
}

func (m Model) renderList() string {
	tabs := m.renderFilterTabs()
	count := fmt.Sprintf("  %d issue(s)", len(m.issues))
	tabLine := lipgloss.JoinHorizontal(lipgloss.Center, tabs,
		lipgloss.NewStyle().Foreground(lipgloss.Color("#616161")).Render(count))

	if len(m.issues) == 0 {
		return tabLine + "\n\n" +
			lipgloss.NewStyle().Foreground(lipgloss.Color("#9E9E9E")).Render("  No issues found for this filter.") +
			"\n" + m.renderHelp()
	}

	rows := make([]string, len(m.issues))
	for i, issue := range m.issues {
		assignee := ""
		if issue.Assignee != "" {
			assignee = lipgloss.NewStyle().Foreground(lipgloss.Color("#80CBC4")).Render("@" + issue.Assignee)
		}

		row := fmt.Sprintf("  %-14s %s  %s  %s  %s  %s",
			lipgloss.NewStyle().Foreground(lipgloss.Color("#90A4AE")).Render(issue.ID),
			PriorityBadge(issue.Priority),
			StatusBadge(issue.Status),
			TypeBadge(issue.Type),
			truncate(issue.Title, 40),
			assignee,
		)

		if i == m.cursor {
			rows[i] = selectedRowStyle.Render(row)
		} else {
			rows[i] = normalRowStyle.Render(row)
		}
	}

	return tabLine + "\n\n" +
		strings.Join(rows, "\n") +
		"\n" + m.renderHelp()
}

func (m Model) renderDetail() string {
	d := m.detail

	lines := []string{
		detailTitleStyle.Render("📋 " + d.Title),
		"",
		detailLabelStyle.Render("ID:       ") + d.ID,
		detailLabelStyle.Render("Status:   ") + StatusBadge(d.Status),
		detailLabelStyle.Render("Priority: ") + PriorityBadge(d.Priority),
		detailLabelStyle.Render("Type:     ") + TypeBadge(d.Type),
	}

	if d.Assignee != "" {
		lines = append(lines, detailLabelStyle.Render("Assignee: ")+"@"+d.Assignee)
	}
	if len(d.Labels) > 0 {
		lines = append(lines, detailLabelStyle.Render("Labels:   ")+strings.Join(d.Labels, ", "))
	}
	if d.Description != "" {
		lines = append(lines, "", detailLabelStyle.Render("Description:"), "  "+wordWrap(d.Description, 60))
	}
	if d.Notes != "" {
		lines = append(lines, "", detailLabelStyle.Render("Notes:"), "  "+wordWrap(d.Notes, 60))
	}
	if len(d.Dependencies) > 0 {
		lines = append(lines, "", detailLabelStyle.Render("Dependencies:"))
		for _, dep := range d.Dependencies {
			lines = append(lines, fmt.Sprintf("  • %s [%s]", dep.ID, dep.Type))
		}
	}
	if len(d.Blockers) > 0 {
		lines = append(lines, "", detailLabelStyle.Render("Blocked by:"))
		for _, b := range d.Blockers {
			lines = append(lines, "  ⛔ "+b)
		}
	}
	if d.CreatedAt != "" {
		lines = append(lines, "", detailLabelStyle.Render("Created:  ")+d.CreatedAt)
	}

	content := strings.Join(lines, "\n")
	panel := detailPanelStyle.Render(content)

	help := helpStyle.Render("esc/q  back to list")
	return panel + "\n" + help
}

func (m Model) renderHelp() string {
	return helpStyle.Render("\n  ↑/↓  navigate  •  tab  filter  •  enter  detail  •  r  refresh  •  q  quit")
}

// truncate shortens a string to maxLen, adding "…" if needed.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-1] + "…"
}

// wordWrap naively wraps text at maxWidth.
func wordWrap(s string, maxWidth int) string {
	if len(s) <= maxWidth {
		return s
	}
	var sb strings.Builder
	words := strings.Fields(s)
	line := 0
	for _, w := range words {
		if line+len(w)+1 > maxWidth {
			sb.WriteString("\n  ")
			line = 0
		}
		if line > 0 {
			sb.WriteString(" ")
			line++
		}
		sb.WriteString(w)
		line += len(w)
	}
	return sb.String()
}
