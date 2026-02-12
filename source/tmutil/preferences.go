//
// preferences.go
// ~~~~~~~~~~~~~~~~~~~~~
//
// Copyright (c) 2004-2026 Metasystems Technologies Inc. (MTI)
//
// Licensed under the MIT License. See LICENSE file in the project root
// for full license text.
//

package tmutil

import (
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const tmPlistDomain = "/Library/Preferences/com.apple.TimeMachine"
const plistTimeLayout = "2006-01-02 15:04:05 +0000"

// BackupPrefs holds data read from the Time Machine preferences plist.
// This data is available even when the backup disk is not mounted.
type BackupPrefs struct {
	AutoBackup     bool
	AutoBackupSet  bool // true if the key was found
	Encryption     string
	BytesUsed      int64
	BytesAvailable int64
	SnapshotDates  []time.Time
	AttemptDates   []time.Time
}

// LastSnapshot returns the most recent snapshot date, or zero time if none.
func (p BackupPrefs) LastSnapshot() time.Time {
	if len(p.SnapshotDates) == 0 {
		return time.Time{}
	}
	return p.SnapshotDates[len(p.SnapshotDates)-1]
}

// FirstSnapshot returns the oldest snapshot date, or zero time if none.
func (p BackupPrefs) FirstSnapshot() time.Time {
	if len(p.SnapshotDates) == 0 {
		return time.Time{}
	}
	return p.SnapshotDates[0]
}

// LastBackupDuration returns the elapsed time of the most recent backup
// by finding the attempt date that corresponds to the last snapshot.
// Returns zero if the data is unavailable.
func (p BackupPrefs) LastBackupDuration() time.Duration {
	last := p.LastSnapshot()
	if last.IsZero() || len(p.AttemptDates) == 0 {
		return 0
	}
	// Find the latest attempt date that is on or before the last snapshot.
	var best time.Time
	for _, a := range p.AttemptDates {
		if !a.After(last) && a.After(best) {
			best = a
		}
	}
	if best.IsZero() {
		return 0
	}
	return last.Sub(best)
}

// GetBackupPrefs reads the Time Machine preferences plist via `defaults read`.
func GetBackupPrefs() (BackupPrefs, error) {
	cmd := exec.Command("defaults", "read", tmPlistDomain)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return BackupPrefs{}, err
	}
	return parseBackupPrefs(string(output)), nil
}

func parseBackupPrefs(raw string) BackupPrefs {
	var prefs BackupPrefs

	// Top-level keys.
	for _, line := range strings.Split(raw, "\n") {
		trimmed := strings.TrimSpace(line)
		trimmed = strings.TrimSuffix(trimmed, ";")
		parts := strings.SplitN(trimmed, " = ", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		val = strings.Trim(val, "\"")
		switch key {
		case "AutoBackup":
			prefs.AutoBackupSet = true
			prefs.AutoBackup = val == "1"
		}
	}

	// Parse within the first Destinations block.
	prefs.parseDestinationBlock(raw)

	return prefs
}

func (p *BackupPrefs) parseDestinationBlock(raw string) {
	// Find the Destinations array and parse the first entry.
	// Tracks named arrays (SnapshotDates, AttemptDates, etc.) so we only
	// collect dates from SnapshotDates.
	lines := strings.Split(raw, "\n")
	inDest := false
	inSnapshots := false
	inAttempts := false
	inArray := "" // name of the current array being parsed
	depth := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if !inDest {
			if strings.HasPrefix(trimmed, "Destinations =") {
				inDest = true
			}
			continue
		}

		// Detect named array openings like "SnapshotDates =             ("
		if strings.Contains(trimmed, " = ") {
			clean := strings.TrimSuffix(trimmed, ";")
			parts := strings.SplitN(clean, " = ", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				val := strings.TrimSpace(parts[1])
				val = strings.Trim(val, "\"")

				// Check if this opens a nested array.
				if val == "(" {
					inArray = key
					if key == "SnapshotDates" {
						inSnapshots = true
					} else if key == "AttemptDates" {
						inAttempts = true
					}
					continue
				}

				// Scalar key-value pairs.
				switch key {
				case "BytesUsed":
					if n, err := strconv.ParseInt(val, 10, 64); err == nil {
						p.BytesUsed = n
					}
				case "BytesAvailable":
					if n, err := strconv.ParseInt(val, 10, 64); err == nil {
						p.BytesAvailable = n
					}
				case "LastKnownEncryptionState":
					p.Encryption = val
				}
			}
			continue
		}

		// Track structural tokens.
		if trimmed == "(" || trimmed == "{" {
			depth++
			continue
		}
		if trimmed == ");" || trimmed == ")" {
			if inArray != "" {
				if inArray == "SnapshotDates" {
					inSnapshots = false
				} else if inArray == "AttemptDates" {
					inAttempts = false
				}
				inArray = ""
			} else {
				depth--
				if depth <= 0 {
					break
				}
			}
			continue
		}
		if trimmed == "};" || trimmed == "}" {
			depth--
			if depth <= 0 {
				break
			}
			continue
		}

		// Collect quoted date strings from SnapshotDates or AttemptDates.
		if inSnapshots || inAttempts {
			clean := strings.TrimSuffix(trimmed, ",")
			clean = strings.TrimSpace(clean)
			if strings.HasPrefix(clean, "\"") && strings.HasSuffix(clean, "\"") {
				val := strings.Trim(clean, "\"")
				if t, err := time.Parse(plistTimeLayout, val); err == nil {
					if inSnapshots {
						p.SnapshotDates = append(p.SnapshotDates, t)
					} else {
						p.AttemptDates = append(p.AttemptDates, t)
					}
				}
			}
		}
	}
}
