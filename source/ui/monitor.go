//
// monitor.go
// ~~~~~~~~~~~~~~~~~~~~~
//
// Copyright (c) 2004-2026 Metasystems Technologies Inc. (MTI)
// All rights reserved
//
// Distributed under the MTI Software License, Version 0.1.
//
// as defined by accompanying file MTI-LICENSE-0.1.info or
// at http://www.mtihq.com/license/MTI-LICENSE-0.1.info
//

package ui

import (
	"fmt"
	"strings"
	"time"

	"tmcli/tmutil"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const pollInterval = 1 * time.Second

type statusTickMsg struct{}
type statusUpdateMsg struct {
	info tmutil.StatusInfo
	err  error
}

func pollStatus() tea.Msg {
	info, err := tmutil.GetStatus()
	return statusUpdateMsg{info: info, err: err}
}

func tickCmd() tea.Cmd {
	return tea.Tick(pollInterval, func(time.Time) tea.Msg {
		return statusTickMsg{}
	})
}

// MonitorModel is a Bubbletea model for monitoring backup progress.
// It can be used standalone (CLI --monitor) or embedded in the TUI.
type MonitorModel struct {
	version  string
	info     tmutil.StatusInfo
	err      error
	width    int
	height   int
	done     bool // backup finished while monitoring
	altScreen bool // true when running as full TUI
}

// NewMonitorModel creates a monitor model.
func NewMonitorModel(version string, altScreen bool) MonitorModel {
	return MonitorModel{version: version, altScreen: altScreen}
}

// Init starts the first poll immediately.
func (m MonitorModel) Init() tea.Cmd {
	return pollStatus
}

// Update handles messages.
func (m MonitorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		}

	case statusUpdateMsg:
		m.err = msg.err
		if msg.err == nil {
			m.info = msg.info
			if !m.info.Running {
				m.done = true
			}
		}
		return m, tickCmd()

	case statusTickMsg:
		return m, pollStatus
	}

	return m, nil
}

// View renders the monitor.
func (m MonitorModel) View() string {
	body := m.renderBody()

	if m.altScreen {
		var b strings.Builder
		content := lipgloss.JoinVertical(lipgloss.Center, "Backup Monitor", m.version)
		b.WriteString(titleStyle.Render(content))
		b.WriteString("\n\n")
		b.WriteString(outputStyle.Render(body))
		b.WriteString("\n\n")
		b.WriteString(helpStyle.Render("b/esc: back • q: quit • updates every 1s"))
		return lipgloss.Place(m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			b.String())
	}

	return body
}

// renderBody builds the monitor content as plain text so alignment is
// consistent regardless of whether it is later wrapped by lipgloss.Place.
func (m MonitorModel) renderBody() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}

	if m.done {
		return fmt.Sprintf("Backup complete.\n\n%s  100.0%%", renderProgressBar(1.0))
	}

	if !m.info.Running {
		return fmt.Sprintf("No backup in progress. Waiting...\n\n%s    0.0%%", renderProgressBar(0))
	}

	var b strings.Builder

	if m.info.Phase != "" {
		fmt.Fprintf(&b, "Phase:       %s\n", m.info.Phase)
	}
	if m.info.Destination != "" {
		fmt.Fprintf(&b, "Destination: %s\n", m.info.Destination)
	}
	b.WriteString("\n")

	pct := m.info.Percent
	fmt.Fprintf(&b, "%s  %.1f%%\n\n", renderProgressBar(pct), pct*100)

	if m.info.TotalBytes > 0 {
		fmt.Fprintf(&b, "Bytes:       %s / %s\n",
			tmutil.FormatBytesInt64(m.info.BytesCopied),
			tmutil.FormatBytesInt64(m.info.TotalBytes))
	}
	if m.info.TotalFiles > 0 {
		fmt.Fprintf(&b, "Files:       %d / %d\n",
			m.info.FilesCopied, m.info.TotalFiles)
	}
	if m.info.TimeRemaining > 0 {
		mins := int(m.info.TimeRemaining) / 60
		hrs := mins / 60
		mins = mins % 60
		estimate := time.Now().Add(time.Duration(m.info.TimeRemaining) * time.Second)
		if hrs > 0 {
			fmt.Fprintf(&b, "Remaining:   %dh %dm [%s]\n", hrs, mins, estimate.Format("2006-01-02 15:04:05"))
		} else {
			fmt.Fprintf(&b, "Remaining:   %dm [%s]\n", mins, estimate.Format("2006-01-02 15:04:05"))
		}
	} else {
		fmt.Fprintf(&b, "Remaining:   Calculating...\n")
	}

	if !m.info.StartedAt.IsZero() {
		now := time.Now()
		elapsed := now.Sub(m.info.StartedAt)
		b.WriteString("\n")
		fmt.Fprintf(&b, "Started:     %s\n", m.info.StartedAt.Local().Format("2006-01-02 15:04:05"))
		fmt.Fprintf(&b, "Current:     %s\n", now.Format("2006-01-02 15:04:05"))
		fmt.Fprintf(&b, "Elapsed:     %s\n", tmutil.FormatDuration(elapsed))
	}

	if !m.altScreen {
		b.WriteString("\nq: quit • updates every 1s")
	}

	return b.String()
}

func renderProgressBar(percent float64) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 1 {
		percent = 1
	}
	filled := int(float64(progressBarWidth) * percent)
	empty := progressBarWidth - filled

	bar := progressFullStyle.Render(strings.Repeat("█", filled)) +
		progressEmptyStyle.Render(strings.Repeat("░", empty))

	return fmt.Sprintf("[%s]", bar)
}
