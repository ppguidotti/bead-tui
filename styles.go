package main

import "github.com/charmbracelet/lipgloss"

var (
	// Layout
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleBarStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#5C6BC0")).
			Padding(0, 2).
			MarginBottom(1)

	// Status colors
	statusStyles = map[string]lipgloss.Style{
		"open":        lipgloss.NewStyle().Foreground(lipgloss.Color("#64B5F6")),
		"in_progress": lipgloss.NewStyle().Foreground(lipgloss.Color("#81C784")).Bold(true),
		"blocked":     lipgloss.NewStyle().Foreground(lipgloss.Color("#EF9A9A")).Bold(true),
		"deferred":    lipgloss.NewStyle().Foreground(lipgloss.Color("#B0BEC5")),
		"closed":      lipgloss.NewStyle().Foreground(lipgloss.Color("#78909C")).Strikethrough(true),
	}

	// Priority colors
	priorityStyles = map[int]lipgloss.Style{
		0: lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5252")).Bold(true), // P0 red
		1: lipgloss.NewStyle().Foreground(lipgloss.Color("#FFB74D")),            // P1 orange
		2: lipgloss.NewStyle().Foreground(lipgloss.Color("#FFF176")),            // P2 yellow
		3: lipgloss.NewStyle().Foreground(lipgloss.Color("#A5D6A7")),            // P3 green
	}

	// Issue type badge
	typeBadgeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#CE93D8")).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#CE93D8")).
			Padding(0, 1)

	// Selected row
	selectedRowStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#303F9F")).
				Foreground(lipgloss.Color("#FFFFFF"))

	normalRowStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E0E0E0"))

	// Detail panel
	detailPanelStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#5C6BC0")).
				Padding(1, 2).
				MarginLeft(2)

	detailTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FAFAFA")).
				MarginBottom(1)

	detailLabelStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#90CAF9")).
				Bold(true)

	// Filter tab
	activeFilterStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(lipgloss.Color("#5C6BC0")).
				Padding(0, 2).
				Bold(true)

	inactiveFilterStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#9E9E9E")).
				Padding(0, 2)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#616161")).
			MarginTop(1)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5252")).
			Bold(true)
)

// StatusBadge renders a colored status string.
func StatusBadge(status string) string {
	s, ok := statusStyles[status]
	if !ok {
		return status
	}
	return s.Render(status)
}

// PriorityBadge renders a colored priority label.
func PriorityBadge(p int) string {
	label := PriorityLabel(p)
	s, ok := priorityStyles[p]
	if !ok {
		return label
	}
	return s.Render(label)
}

// TypeBadge renders an issue type badge.
func TypeBadge(t string) string {
	if t == "" {
		t = "task"
	}
	return typeBadgeStyle.Render(t)
}
