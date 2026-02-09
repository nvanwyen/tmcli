//
// backup.go
// ~~~~~~~~~~~~~~~~~~~~~
//
// Copyright (c) 2004-2026 Metasystems Technologies Inc. (MTI)
//
// Licensed under the MIT License. See LICENSE file in the project root
// for full license text.
//

package tmutil

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// StartBackup starts a Time Machine backup.
func StartBackup() (string, error) {
	output, err := run("startbackup")
	if err != nil {
		return "", err
	}
	if output == "" {
		return "Backup started successfully.", nil
	}
	return output, nil
}

// StopBackup stops a running Time Machine backup.
func StopBackup() (string, error) {
	output, err := run("stopbackup")
	if err != nil {
		return "", err
	}
	if output == "" {
		return "Backup stopped successfully.", nil
	}
	return output, nil
}

// Status returns a human-readable status of the current backup.
func Status() (string, error) {
	output, err := run("status")
	if err != nil {
		return "", err
	}
	return formatStatus(output), nil
}

// Enable enables automatic Time Machine backups.
func Enable() (string, error) {
	output, err := run("enable")
	if err != nil {
		return "", err
	}
	if output == "" {
		return "Time Machine enabled.", nil
	}
	return output, nil
}

// Disable disables automatic Time Machine backups.
func Disable() (string, error) {
	output, err := run("disable")
	if err != nil {
		return "", err
	}
	if output == "" {
		return "Time Machine disabled.", nil
	}
	return output, nil
}

// Version returns the tmutil version.
func Version() (string, error) {
	return run("version")
}

func formatStatus(raw string) string {
	if strings.Contains(raw, "Running = 0") {
		return "Time Machine is idle — no backup in progress."
	}

	fields := parseFields(raw)
	var b strings.Builder

	b.WriteString("Time Machine Backup Status\n")
	b.WriteString(strings.Repeat("─", 40) + "\n\n")

	if v, ok := fields["BackupPhase"]; ok {
		b.WriteString(fmt.Sprintf("  Phase:         %s\n", v))
	}
	if v, ok := fields["Running"]; ok {
		if v == "1" {
			b.WriteString("  Running:       Yes\n")
		} else {
			b.WriteString("  Running:       No\n")
		}
	}
	if v, ok := fields["DestinationMountPoint"]; ok {
		b.WriteString(fmt.Sprintf("  Destination:   %s\n", v))
	}
	if v, ok := fields["DateOfStateChange"]; ok {
		b.WriteString(fmt.Sprintf("  Started:       %s\n", v))
		if t, err := time.Parse(tmutilTimeLayout, v); err == nil {
			elapsed := time.Since(t)
			b.WriteString(fmt.Sprintf("  Elapsed:       %s\n", FormatDuration(elapsed)))
		}
	}

	progress := parseProgress(raw)
	if len(progress) > 0 {
		b.WriteString("\n  Progress:\n")
		if v, ok := progress["Percent"]; ok {
			if pct, err := strconv.ParseFloat(v, 64); err == nil {
				b.WriteString(fmt.Sprintf("    Completed:   %.1f%%\n", pct*100))
			}
		}
		if v, ok := progress["TimeRemaining"]; ok {
			if secs, err := strconv.ParseFloat(v, 64); err == nil && secs > 0 {
				mins := int(secs) / 60
				hrs := mins / 60
				mins = mins % 60
				estimate := time.Now().Add(time.Duration(secs) * time.Second)
				if hrs > 0 {
					b.WriteString(fmt.Sprintf("    Remaining:   %dh %dm [%s]\n", hrs, mins, estimate.Local().Format("2006-01-02 15:04:05")))
				} else {
					b.WriteString(fmt.Sprintf("    Remaining:   %dm [%s]\n", mins, estimate.Local().Format("2006-01-02 15:04:05")))
				}
			} else {
				b.WriteString("    Remaining:   Calculating...\n")
			}
		} else {
			b.WriteString("    Remaining:   Calculating...\n")
		}
		if bytes, ok := progress["bytes"]; ok {
			if total, ok2 := progress["totalBytes"]; ok2 {
				b.WriteString(fmt.Sprintf("    Bytes:       %s / %s\n", formatBytes(bytes), formatBytes(total)))
			}
		}
		if files, ok := progress["files"]; ok {
			if total, ok2 := progress["totalFiles"]; ok2 {
				b.WriteString(fmt.Sprintf("    Files:       %s / %s\n", files, total))
			}
		}
	}

	return b.String()
}
