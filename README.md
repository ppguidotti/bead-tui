# bd-tui

A terminal user interface (TUI) for [Beads](https://github.com/steveyegge/beads) — the distributed, git-backed task tracker for AI agents.

## Demo

```
⚡ bd-tui — Beads Task Viewer

 all   open   in_progress   blocked   closed     5 issue(s)

  bd-a1b2   P0  in_progress  [feature]  Add OAuth support          @alice
  bd-c3d4   P1  open         [bug]      Fix token refresh
  bd-e5f6   P1  blocked      [task]     Deploy to staging           @bob
  bd-g7h8   P2  open         [chore]    Update dependencies
  bd-i9j0   P3  deferred     [task]     Refactor config loader

  ↑/↓ navigate  •  tab filter  •  enter detail  •  r refresh  •  q quit
```

## Features

- **Live data**: Reads directly from `bd list --json` and `bd show --json`
- **Filter tabs**: Quickly switch between all / open / in_progress / blocked / closed
- **Colour coding**: Priority (P0→P3) and status each have distinct colours
- **Detail view**: Press Enter to see full task info including dependencies and blockers
- **Async loading**: Non-blocking IO with a spinner — the UI never freezes
- **Refresh**: Press `r` to reload from the database

## Prerequisites

- Go 1.22+
- `bd` CLI installed: `npm install -g @beads/bd` or `brew install beads`
- A Beads-initialized project: `bd init`

## Installation

```bash
git clone https://github.com/ppguidotti/bead-tui
cd bead-tui
go mod tidy
go build -o bead-tui .

# Run from inside a Beads project:
cd your-project
/path/to/bead-tui
```

Or install globally:
```bash
go install github.com/ppguidotti/bead-tui@latest
```

## Architecture

```
bd-tui/
├── main.go      Entry point — starts bubbletea program
├── model.go     All UI state and event handling (bubbletea Model/Update/View)
├── beads.go     Thin wrapper around `bd` CLI: runs commands, parses JSON
├── styles.go    All lipgloss styles in one place (colours, borders, badges)
└── go.mod
```

The project follows the **Elm architecture** used by bubbletea:

- **Model**: pure data struct (no side effects)
- **Update**: pure function mapping `(Model, Msg) → (Model, Cmd)`
- **View**: pure function mapping `Model → string`
- **Cmd**: async IO (shell commands) run outside the model

This makes the TUI fully testable: you can unit-test `Update` without touching the terminal or the `bd` binary.

## Design decisions

**Why wrap `bd` CLI instead of using the Go library directly?**
The `beads` Go package API is minimal and intended for orchestration. The recommended integration pattern from the Beads docs is `bd --json` CLI output. This also means bd-tui works with any Beads backend (SQLite or Dolt server mode) without changes.

**Why bubbletea?**
It's the de-facto standard for Go TUIs, used widely in DevOps tooling (e.g. Charm's own ecosystem). The Elm-inspired architecture makes state management explicit and side-effect free.

## Keyboard shortcuts

| Key | Action |
|-----|--------|
| `↑` / `k` | Move up |
| `↓` / `j` | Move down |
| `Tab` / `→` | Next filter |
| `Shift+Tab` / `←` | Previous filter |
| `Enter` | Open detail view |
| `Esc` / `q` | Back / Quit |
| `r` | Refresh issues |