//
// restore.go
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
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// backupPathDateLayout is the date format embedded in backup path names.
const backupPathDateLayout = "2006-01-02-150405"

// defaultFindLimit is the default number of backup snapshots to search.
const defaultFindLimit = 5

// Restore restores files from a backup.
func Restore(args []string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("source and destination paths are required")
	}
	cmdArgs := append([]string{"restore", "-v"}, args...)
	return run(cmdArgs...)
}

// FindFile searches for a filename/pattern across recent backup snapshots.
// args[0] = filename or glob pattern (required)
// args[1] = max number of backups to search (optional, default 5)
func FindFile(args []string) (string, error) {
	if len(args) == 0 || args[0] == "" {
		return "", fmt.Errorf("filename or pattern is required")
	}
	pattern := args[0]
	limit := defaultFindLimit
	if len(args) > 1 && args[1] != "" {
		n, err := strconv.Atoi(args[1])
		if err == nil && n > 0 {
			limit = n
		}
	}

	backups, err := listBackupPaths()
	if err != nil {
		return "", err
	}
	reverseStrings(backups)
	if len(backups) > limit {
		backups = backups[:limit]
	}

	var results []string
	for _, bp := range backups {
		matches, walkErr := findInBackup(bp, pattern)
		if walkErr != nil {
			results = append(results, fmt.Sprintf("# Error scanning %s: %v", bp, walkErr))
			continue
		}
		results = append(results, matches...)
	}

	if len(results) == 0 {
		return fmt.Sprintf("No matches for %q in the last %d backup(s).", pattern, len(backups)), nil
	}

	matchCount := 0
	for _, r := range results {
		if !strings.HasPrefix(r, "#") {
			matchCount++
		}
	}
	header := fmt.Sprintf("Found %d match(es) for %q across %d backup(s):\n",
		matchCount, pattern, len(backups))
	return header + strings.Join(results, "\n"), nil
}

// FindByDate lists backup snapshots within a date range.
// args[0] = start date YYYY-MM-DD (required)
// args[1] = end date YYYY-MM-DD (optional; defaults to today)
func FindByDate(args []string) (string, error) {
	if len(args) == 0 || args[0] == "" {
		return "", fmt.Errorf("start date (YYYY-MM-DD) is required")
	}
	startDate, err := time.Parse("2006-01-02", args[0])
	if err != nil {
		return "", fmt.Errorf("invalid start date %q: expected YYYY-MM-DD", args[0])
	}
	endDate := time.Now()
	if len(args) > 1 && args[1] != "" {
		endDate, err = time.Parse("2006-01-02", args[1])
		if err != nil {
			return "", fmt.Errorf("invalid end date %q: expected YYYY-MM-DD", args[1])
		}
		endDate = endDate.Add(24*time.Hour - time.Second)
	}

	backups, err := listBackupPaths()
	if err != nil {
		return "", err
	}

	var matches []string
	for _, bp := range backups {
		t, parseErr := parseBackupDate(bp)
		if parseErr != nil {
			continue
		}
		if (t.Equal(startDate) || t.After(startDate)) && (t.Equal(endDate) || t.Before(endDate)) {
			matches = append(matches, bp)
		}
	}

	if len(matches) == 0 {
		return fmt.Sprintf("No backups found between %s and %s.",
			startDate.Format("2006-01-02"),
			endDate.Format("2006-01-02")), nil
	}
	header := fmt.Sprintf("Found %d backup(s) between %s and %s:\n",
		len(matches),
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"))
	return header + strings.Join(matches, "\n"), nil
}

// BrowseBackup lists the contents of a backup snapshot directory.
// args[0] = backup path (required)
// args[1] = subdirectory within the backup (optional)
func BrowseBackup(args []string) (string, error) {
	if len(args) == 0 || args[0] == "" {
		return "", fmt.Errorf("backup path is required")
	}
	dir := args[0]
	if len(args) > 1 && args[1] != "" {
		dir = filepath.Join(args[0], args[1])
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("cannot read %s: %w", dir, err)
	}

	var b strings.Builder
	fmt.Fprintf(&b, "Contents of %s\n", dir)
	b.WriteString(strings.Repeat("─", 60) + "\n\n")
	for _, entry := range entries {
		info, infoErr := entry.Info()
		if entry.IsDir() {
			fmt.Fprintf(&b, "  %-40s  %s\n", entry.Name()+"/", "<dir>")
		} else if infoErr == nil {
			fmt.Fprintf(&b, "  %-40s  %s\n", entry.Name(), FormatBytesInt64(info.Size()))
		} else {
			fmt.Fprintf(&b, "  %s\n", entry.Name())
		}
	}
	fmt.Fprintf(&b, "\n%d item(s)", len(entries))
	return b.String(), nil
}

// --- internal helpers ---

// listBackupPaths calls tmutil listbackups and returns the paths as a slice.
func listBackupPaths() ([]string, error) {
	output, err := run("listbackups")
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(output), "\n")
	var paths []string
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l != "" {
			paths = append(paths, l)
		}
	}
	if len(paths) == 0 {
		return nil, fmt.Errorf("no backups found")
	}
	return paths, nil
}

// parseBackupDate extracts the date from the last path component of a backup path.
func parseBackupDate(backupPath string) (time.Time, error) {
	base := filepath.Base(backupPath)
	return time.Parse(backupPathDateLayout, base)
}

// findInBackup walks a backup snapshot looking for entries matching a glob pattern.
func findInBackup(backupPath, pattern string) ([]string, error) {
	var matches []string
	err := filepath.WalkDir(backupPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if d != nil && d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		matched, matchErr := filepath.Match(pattern, d.Name())
		if matchErr != nil {
			return matchErr
		}
		if matched {
			matches = append(matches, path)
		}
		return nil
	})
	return matches, err
}

// reverseStrings reverses a string slice in place.
func reverseStrings(s []string) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
