//
// snapshot.go
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

// LocalSnapshot creates a new local snapshot.
func LocalSnapshot() (string, error) {
	output, err := run("localsnapshot")
	if err != nil {
		return "", err
	}
	if output == "" {
		return "Local snapshot created.", nil
	}
	return output, nil
}

// ListLocalSnapshots lists local snapshots for a mount point.
func ListLocalSnapshots(args []string) (string, error) {
	mountPoint := "/"
	if len(args) > 0 && args[0] != "" {
		mountPoint = args[0]
	}
	return run("listlocalsnapshots", mountPoint)
}

// ListLocalSnapshotDates lists local snapshot dates for a mount point.
func ListLocalSnapshotDates(args []string) (string, error) {
	mountPoint := "/"
	if len(args) > 0 && args[0] != "" {
		mountPoint = args[0]
	}
	return run("listlocalsnapshotdates", mountPoint)
}

// DeleteLocalSnapshots deletes local snapshots for a mount point or date.
func DeleteLocalSnapshots(args []string) (string, error) {
	if len(args) == 0 || args[0] == "" {
		return "", fmt.Errorf("mount point or snapshot date is required")
	}
	output, err := run("deletelocalsnapshots", args[0])
	if err != nil {
		return "", err
	}
	if output == "" {
		return fmt.Sprintf("Local snapshots deleted for %s.", args[0]), nil
	}
	return output, nil
}

// ThinLocalSnapshots thins local snapshots for a mount point.
func ThinLocalSnapshots(args []string) (string, error) {
	if len(args) == 0 || args[0] == "" {
		return "", fmt.Errorf("mount point is required")
	}
	cmdArgs := []string{"thinlocalsnapshots", args[0]}
	if len(args) > 1 && args[1] != "" {
		cmdArgs = append(cmdArgs, args[1]) // purge amount
	}
	if len(args) > 2 && args[2] != "" {
		cmdArgs = append(cmdArgs, args[2]) // urgency
	}
	output, err := run(cmdArgs...)
	if err != nil {
		return "", err
	}
	if output == "" {
		return fmt.Sprintf("Local snapshots thinned for %s.", args[0]), nil
	}
	return output, nil
}
