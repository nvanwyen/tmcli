//
// browse.go
// ~~~~~~~~~~~~~~~~~~~~~
//
// Copyright (c) 2004-2026 Metasystems Technologies Inc. (MTI)
//
// Licensed under the MIT License. See LICENSE file in the project root
// for full license text.
//

package tmutil

import "fmt"

// LatestBackup returns the path to the most recent backup.
func LatestBackup() (string, error) {
	return run("latestbackup")
}

// ListBackups lists all completed backups.
func ListBackups() (string, error) {
	return run("listbackups")
}

// MachineDirectory returns the machine backup directory path.
func MachineDirectory() (string, error) {
	return run("machinedirectory")
}

// Compare compares the current system to a backup or two paths.
func Compare(args []string) (string, error) {
	cmdArgs := []string{"compare"}
	cmdArgs = append(cmdArgs, args...)
	return run(cmdArgs...)
}

// UniqueSize calculates the unique size of a path in backups.
func UniqueSize(args []string) (string, error) {
	if len(args) == 0 || args[0] == "" {
		return "", fmt.Errorf("path is required")
	}
	cmdArgs := append([]string{"uniquesize"}, args...)
	return run(cmdArgs...)
}

// VerifyChecksums verifies checksums for a path in backups.
func VerifyChecksums(args []string) (string, error) {
	if len(args) == 0 || args[0] == "" {
		return "", fmt.Errorf("path is required")
	}
	cmdArgs := append([]string{"verifychecksums"}, args...)
	return run(cmdArgs...)
}
