//
// advanced.go
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

package tmutil

import "fmt"

// AssociateDisk associates a volume with a backup.
func AssociateDisk(args []string) (string, error) {
	if len(args) < 2 || args[0] == "" || args[1] == "" {
		return "", fmt.Errorf("mount point and volume backup directory are required")
	}
	output, err := run("associatedisk", args[0], args[1])
	if err != nil {
		return "", err
	}
	if output == "" {
		return fmt.Sprintf("Disk associated: %s -> %s.", args[0], args[1]), nil
	}
	return output, nil
}

// InheritBackup inherits a machine directory or sparse bundle.
func InheritBackup(args []string) (string, error) {
	if len(args) == 0 || args[0] == "" {
		return "", fmt.Errorf("machine directory or sparse bundle path is required")
	}
	output, err := run("inheritbackup", args[0])
	if err != nil {
		return "", err
	}
	if output == "" {
		return fmt.Sprintf("Inherited backup from %s.", args[0]), nil
	}
	return output, nil
}

// CalculateDrift calculates drift for a machine directory.
func CalculateDrift(args []string) (string, error) {
	if len(args) == 0 || args[0] == "" {
		return "", fmt.Errorf("machine directory is required")
	}
	return run("calculatedrift", args[0])
}

// DeleteInProgress deletes an in-progress backup.
func DeleteInProgress(args []string) (string, error) {
	if len(args) == 0 || args[0] == "" {
		return "", fmt.Errorf("machine directory is required")
	}
	output, err := run("deleteinprogress", args[0])
	if err != nil {
		return "", err
	}
	if output == "" {
		return fmt.Sprintf("In-progress backup deleted for %s.", args[0]), nil
	}
	return output, nil
}

// Delete deletes a backup snapshot.
func Delete(args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("arguments are required (e.g. -d mount_point -t timestamp or -p path)")
	}
	cmdArgs := append([]string{"delete"}, args...)
	output, err := run(cmdArgs...)
	if err != nil {
		return "", err
	}
	if output == "" {
		return "Backup deleted.", nil
	}
	return output, nil
}
