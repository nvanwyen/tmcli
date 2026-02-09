//
// tmutil.go
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
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const tmutilTimeLayout = "2006-01-02 15:04:05 -0700"

func run(args ...string) (string, error) {
	cmd := exec.Command("tmutil", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s: %w", strings.TrimSpace(string(output)), err)
	}
	return strings.TrimSpace(string(output)), nil
}

// StatusInfo holds structured status data from tmutil.
type StatusInfo struct {
	Running       bool
	Phase         string
	Destination   string
	StartedAt     time.Time
	Percent       float64
	TimeRemaining float64 // seconds
	BytesCopied   int64
	TotalBytes    int64
	FilesCopied   int64
	TotalFiles    int64
}

// GetStatus returns structured backup status information.
func GetStatus() (StatusInfo, error) {
	output, err := run("status")
	if err != nil {
		return StatusInfo{}, err
	}
	return parseStatusInfo(output), nil
}

func parseStatusInfo(raw string) StatusInfo {
	fields := parseFields(raw)
	progress := parseProgress(raw)
	var info StatusInfo

	info.Running = fields["Running"] == "1"
	info.Phase = fields["BackupPhase"]
	info.Destination = fields["DestinationMountPoint"]
	if v, ok := fields["DateOfStateChange"]; ok {
		info.StartedAt, _ = time.Parse(tmutilTimeLayout, v)
	}

	if v, ok := progress["Percent"]; ok {
		info.Percent, _ = strconv.ParseFloat(v, 64)
	}
	if v, ok := progress["TimeRemaining"]; ok {
		info.TimeRemaining, _ = strconv.ParseFloat(v, 64)
	}
	if v, ok := progress["bytes"]; ok {
		f, _ := strconv.ParseFloat(v, 64)
		info.BytesCopied = int64(f)
	}
	if v, ok := progress["totalBytes"]; ok {
		f, _ := strconv.ParseFloat(v, 64)
		info.TotalBytes = int64(f)
	}
	if v, ok := progress["files"]; ok {
		f, _ := strconv.ParseFloat(v, 64)
		info.FilesCopied = int64(f)
	}
	if v, ok := progress["totalFiles"]; ok {
		f, _ := strconv.ParseFloat(v, 64)
		info.TotalFiles = int64(f)
	}

	return info
}

func parseFields(raw string) map[string]string {
	fields := make(map[string]string)
	lines := strings.Split(raw, "\n")
	inProgress := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "Progress =") {
			inProgress = true
			continue
		}
		if inProgress {
			if trimmed == "};" {
				inProgress = false
			}
			continue
		}
		if trimmed == "{" || trimmed == "}" || trimmed == "" {
			continue
		}
		trimmed = strings.TrimSuffix(trimmed, ";")
		parts := strings.SplitN(trimmed, " = ", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			value = strings.Trim(value, "\"")
			fields[key] = value
		}
	}
	return fields
}

func parseProgress(raw string) map[string]string {
	progress := make(map[string]string)
	lines := strings.Split(raw, "\n")
	inProgress := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "Progress =") {
			inProgress = true
			continue
		}
		if inProgress {
			if trimmed == "};" {
				break
			}
			if trimmed == "{" {
				continue
			}
			trimmed = strings.TrimSuffix(trimmed, ";")
			parts := strings.SplitN(trimmed, " = ", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				value = strings.Trim(value, "\"")
				progress[key] = value
			}
		}
	}
	return progress
}

// FormatDuration formats a time.Duration in human-readable uptime style.
func FormatDuration(d time.Duration) string {
	if d < time.Second {
		return "0 seconds"
	}

	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	var parts []string
	if days > 0 {
		if days == 1 {
			parts = append(parts, "1 day")
		} else {
			parts = append(parts, fmt.Sprintf("%d days", days))
		}
	}
	if hours > 0 {
		if hours == 1 {
			parts = append(parts, "1 hour")
		} else {
			parts = append(parts, fmt.Sprintf("%d hours", hours))
		}
	}
	if minutes > 0 {
		if minutes == 1 {
			parts = append(parts, "1 minute")
		} else {
			parts = append(parts, fmt.Sprintf("%d minutes", minutes))
		}
	}
	if seconds > 0 {
		if seconds == 1 {
			parts = append(parts, "1 second")
		} else {
			parts = append(parts, fmt.Sprintf("%d seconds", seconds))
		}
	}
	return strings.Join(parts, ", ")
}

// FormatBytesInt64 formats a byte count as a human-readable string.
func FormatBytesInt64(n int64) string {
	f := float64(n)
	switch {
	case f >= 1e12:
		return fmt.Sprintf("%.1f TB", f/1e12)
	case f >= 1e9:
		return fmt.Sprintf("%.1f GB", f/1e9)
	case f >= 1e6:
		return fmt.Sprintf("%.1f MB", f/1e6)
	case f >= 1e3:
		return fmt.Sprintf("%.1f KB", f/1e3)
	default:
		return fmt.Sprintf("%d B", n)
	}
}

func formatBytes(s string) string {
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return s
	}
	switch {
	case n >= 1e12:
		return fmt.Sprintf("%.1f TB", n/1e12)
	case n >= 1e9:
		return fmt.Sprintf("%.1f GB", n/1e9)
	case n >= 1e6:
		return fmt.Sprintf("%.1f MB", n/1e6)
	case n >= 1e3:
		return fmt.Sprintf("%.1f KB", n/1e3)
	default:
		return fmt.Sprintf("%.0f B", n)
	}
}
