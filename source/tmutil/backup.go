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
		return formatIdleStatus(raw)
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

func formatIdleStatus(raw string) string {
	fields := parseFields(raw)
	var b strings.Builder

	b.WriteString("Time Machine Status\n")
	b.WriteString(strings.Repeat("─", 40) + "\n\n")
	b.WriteString("  State:         Idle\n")

	// Read preferences plist for rich data (available even when disk is unmounted).
	prefs, prefsErr := GetBackupPrefs()

	// Auto-backup: prefer plist, fall back to tmutil status fields.
	if prefsErr == nil && prefs.AutoBackupSet {
		if prefs.AutoBackup {
			b.WriteString("  Auto Backup:   Enabled\n")
		} else {
			b.WriteString("  Auto Backup:   Disabled\n")
		}
	} else if v, ok := fields["AutoBackup"]; ok {
		if v == "1" {
			b.WriteString("  Auto Backup:   Enabled\n")
		} else {
			b.WriteString("  Auto Backup:   Disabled\n")
		}
	}

	// Destination info (best-effort).
	if dest, err := GetDestinationInfo(); err == nil && dest.Name != "" {
		label := dest.Name
		if dest.MountPoint != "" && dest.MountPoint != dest.Name {
			label += " (" + dest.MountPoint
			if dest.Kind != "" {
				label += ", " + dest.Kind
			}
			label += ")"
		} else if dest.Kind != "" {
			label += " (" + dest.Kind + ")"
		}
		b.WriteString(fmt.Sprintf("  Destination:   %s\n", label))
	}

	if prefsErr == nil && prefs.Encryption != "" {
		b.WriteString(fmt.Sprintf("  Encryption:    %s\n", prefs.Encryption))
	}

	// Last backup: prefer plist SnapshotDates (works without disk), fall back to tmutil latestbackup.
	lastShown := false
	if prefsErr == nil && !prefs.LastSnapshot().IsZero() {
		t := prefs.LastSnapshot()
		b.WriteString("\n  Last Backup\n")
		b.WriteString(fmt.Sprintf("    Completed:   %s\n", t.Local().Format("2006-01-02 15:04:05")))
		if d := prefs.LastBackupDuration(); d > 0 {
			b.WriteString(fmt.Sprintf("    Elapsed:     %s\n", FormatDuration(d)))
		}
		lastShown = true
	}
	if !lastShown {
		if latest, err := LatestBackup(); err == nil && latest != "" {
			if t, parseErr := parseBackupDate(latest); parseErr == nil {
				b.WriteString("\n  Last Backup\n")
				b.WriteString(fmt.Sprintf("    Completed:   %s\n", t.Local().Format("2006-01-02 15:04:05")))
				lastShown = true
			}
		}
	}

	// Backup history: prefer plist SnapshotDates, fall back to tmutil listbackups.
	historyShown := false
	if prefsErr == nil && len(prefs.SnapshotDates) > 0 {
		b.WriteString("\n  Backup History\n")
		b.WriteString(fmt.Sprintf("    Total:       %d snapshot(s)\n", len(prefs.SnapshotDates)))
		b.WriteString(fmt.Sprintf("    Oldest:      %s\n", prefs.FirstSnapshot().Local().Format("2006-01-02 15:04:05")))
		b.WriteString(fmt.Sprintf("    Newest:      %s\n", prefs.LastSnapshot().Local().Format("2006-01-02 15:04:05")))
		historyShown = true
	}
	if !historyShown {
		if paths, err := listBackupPaths(); err == nil && len(paths) > 0 {
			b.WriteString("\n  Backup History\n")
			b.WriteString(fmt.Sprintf("    Total:       %d snapshot(s)\n", len(paths)))
			if oldest, err := parseBackupDate(paths[0]); err == nil {
				b.WriteString(fmt.Sprintf("    Oldest:      %s\n", oldest.Local().Format("2006-01-02 15:04:05")))
			}
			if newest, err := parseBackupDate(paths[len(paths)-1]); err == nil {
				b.WriteString(fmt.Sprintf("    Newest:      %s\n", newest.Local().Format("2006-01-02 15:04:05")))
			}
		}
	}

	// Disk usage from plist.
	if prefsErr == nil && (prefs.BytesUsed > 0 || prefs.BytesAvailable > 0) {
		b.WriteString("\n  Disk Usage\n")
		if prefs.BytesUsed > 0 {
			b.WriteString(fmt.Sprintf("    Used:        %s\n", FormatBytesInt64(prefs.BytesUsed)))
		}
		if prefs.BytesAvailable > 0 {
			b.WriteString(fmt.Sprintf("    Available:   %s\n", FormatBytesInt64(prefs.BytesAvailable)))
		}
		if prefs.BytesUsed > 0 && prefs.BytesAvailable > 0 {
			total := prefs.BytesUsed + prefs.BytesAvailable
			b.WriteString(fmt.Sprintf("    Total:       %s\n", FormatBytesInt64(total)))
		}
	}

	return b.String()
}
