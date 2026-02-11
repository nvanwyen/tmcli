//
// destination.go
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
	"strings"
)

// DestInfo holds structured destination information.
type DestInfo struct {
	Name       string
	Kind       string
	MountPoint string
	ID         string
}

// DestinationInfo returns backup destination details.
func DestinationInfo() (string, error) {
	return run("destinationinfo")
}

// GetDestinationInfo returns structured destination information.
func GetDestinationInfo() (DestInfo, error) {
	raw, err := run("destinationinfo")
	if err != nil {
		return DestInfo{}, err
	}
	return parseDestinationInfo(raw), nil
}

func parseDestinationInfo(raw string) DestInfo {
	var info DestInfo
	for _, line := range strings.Split(raw, "\n") {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		switch key {
		case "Name":
			info.Name = val
		case "Kind":
			info.Kind = val
		case "Mount Point":
			info.MountPoint = val
		case "ID":
			info.ID = val
		}
	}
	return info
}

// SetDestination sets a backup destination mount point.
func SetDestination(args []string) (string, error) {
	if len(args) == 0 || args[0] == "" {
		return "", fmt.Errorf("mount point is required")
	}
	cmdArgs := append([]string{"setdestination"}, args...)
	output, err := run(cmdArgs...)
	if err != nil {
		return "", err
	}
	if output == "" {
		return fmt.Sprintf("Destination set to %s.", args[len(args)-1]), nil
	}
	return output, nil
}

// RemoveDestination removes a backup destination by ID.
func RemoveDestination(args []string) (string, error) {
	if len(args) == 0 || args[0] == "" {
		return "", fmt.Errorf("destination ID is required")
	}
	output, err := run("removedestination", args[0])
	if err != nil {
		return "", err
	}
	if output == "" {
		return fmt.Sprintf("Destination %s removed.", args[0]), nil
	}
	return output, nil
}

// SetQuota sets the quota for a destination in gigabytes.
func SetQuota(args []string) (string, error) {
	if len(args) < 2 || args[0] == "" || args[1] == "" {
		return "", fmt.Errorf("destination ID and quota (GB) are required")
	}
	output, err := run("setquota", args[0], args[1])
	if err != nil {
		return "", err
	}
	if output == "" {
		return fmt.Sprintf("Quota set to %s GB for destination %s.", args[1], args[0]), nil
	}
	return output, nil
}
