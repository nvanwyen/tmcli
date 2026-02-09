//
// main.go
// ~~~~~~~~~~~~~~~~~~~~~
//
// Copyright (c) 2004-2026 Metasystems Technologies Inc. (MTI)
//
// Licensed under the MIT License. See LICENSE file in the project root
// for full license text.
//

package main

import (
	"fmt"
	"os"
	"strings"

	"tmcli/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if len(os.Args) < 2 {
		runTUI()
		return
	}

	verb := os.Args[1]
	args := os.Args[2:]

	switch verb {
	case "--version", "-version", "-v", "version":
		fmt.Printf("Time Machine CLI %s\n\n", Version)
		fmt.Println("Copyright (c) 2004-2026 Metasystems Technologies Inc. (MTI)")
		fmt.Println("Licensed under the MIT License.")
		return
	case "--help", "-help", "-h", "help":
		printUsage()
	case "tui":
		runTUI()
	case "monitor":
		runMonitor()
	default:
		cmd := ui.FindCommand(verb)
		if cmd == nil {
			fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", verb)
			printUsage()
			os.Exit(1)
		}
		if cmd.IsMonitor {
			runMonitor()
			return
		}
		runCLI(cmd.Execute, args)
	}
}

func runCLI(fn func([]string) (string, error), args []string) {
	output, err := fn(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(output)
}

func runMonitor() {
	p := tea.NewProgram(ui.NewMonitorModel(Version, false))
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runTUI() {
	p := tea.NewProgram(ui.NewModel(Version), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "tmcli %s\n\n", Version)
	fmt.Fprintf(os.Stderr, "Usage: %s [command] [arguments]\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Running with no arguments launches the interactive TUI.\n\n")
	fmt.Fprintf(os.Stderr, "Commands:\n")
	fmt.Fprintf(os.Stderr, "  %-26s %s\n", "tui", "Launch the interactive TUI (default)")
	fmt.Fprintf(os.Stderr, "\n")

	for _, cat := range ui.Categories() {
		fmt.Fprintf(os.Stderr, "  %s:\n", strings.ToUpper(cat.Title))
		for _, cmd := range cat.Commands {
			desc := cmd.Title
			if len(cmd.Inputs) > 0 {
				var params []string
				for _, inp := range cmd.Inputs {
					if inp.Required {
						params = append(params, fmt.Sprintf("<%s>", inp.Label))
					} else {
						params = append(params, fmt.Sprintf("[%s]", inp.Label))
					}
				}
				desc = fmt.Sprintf("%s %s", cmd.ID, strings.Join(params, " "))
			}
			fmt.Fprintf(os.Stderr, "    %-24s %s\n", cmd.ID, desc)
		}
		fmt.Fprintf(os.Stderr, "\n")
	}
}
