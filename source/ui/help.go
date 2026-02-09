//
// help.go
// ~~~~~~~~~~~~~~~~~~~~~
//
// Copyright (c) 2004-2026 Metasystems Technologies Inc. (MTI)
//
// Licensed under the MIT License. See LICENSE file in the project root
// for full license text.
//

package ui

import (
	"fmt"
	"strings"
)

// BuildCommandHelp generates detailed help for a single command.
func BuildCommandHelp(cmd Command) string {
	var b strings.Builder

	fmt.Fprintf(&b, "%s\n", cmd.Title)
	b.WriteString(strings.Repeat("─", 40) + "\n\n")

	// Description
	b.WriteString(wordWrap(cmd.Description, 48))
	b.WriteString("\n\n")

	// Hotkey
	fmt.Fprintf(&b, "Hotkey:  %s\n", cmd.Hotkey)
	if cmd.RequiresRoot {
		fmt.Fprintf(&b, "Root:    yes\n")
	}

	// CLI usage
	if cmd.IsMonitor {
		fmt.Fprintf(&b, "CLI:     tmcli %s\n", cmd.ID)
	} else if len(cmd.Inputs) > 0 {
		var params []string
		for _, inp := range cmd.Inputs {
			if inp.Required {
				params = append(params, fmt.Sprintf("<%s>", inp.Label))
			} else {
				params = append(params, fmt.Sprintf("[%s]", inp.Label))
			}
		}
		fmt.Fprintf(&b, "CLI:     tmcli %s %s\n", cmd.ID, strings.Join(params, " "))
	} else {
		fmt.Fprintf(&b, "CLI:     tmcli %s\n", cmd.ID)
	}

	// Parameters
	if len(cmd.Inputs) > 0 {
		b.WriteString("\nParameters:\n")
		for _, inp := range cmd.Inputs {
			req := "optional"
			if inp.Required {
				req = "required"
			}
			fmt.Fprintf(&b, "  %-22s %s\n", inp.Label, req)
			if inp.Placeholder != "" {
				fmt.Fprintf(&b, "  %-22s e.g. %s\n", "", inp.Placeholder)
			}
		}
	}

	return b.String()
}

// wordWrap wraps text at the given width on word boundaries.
func wordWrap(text string, width int) string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return ""
	}

	var lines []string
	line := words[0]
	for _, w := range words[1:] {
		if len(line)+1+len(w) > width {
			lines = append(lines, line)
			line = w
		} else {
			line += " " + w
		}
	}
	lines = append(lines, line)
	return strings.Join(lines, "\n")
}
