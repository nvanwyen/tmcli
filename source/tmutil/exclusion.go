//
// exclusion.go
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

// AddExclusion adds a fixed-path exclusion for an item.
func AddExclusion(args []string) (string, error) {
	if len(args) == 0 || args[0] == "" {
		return "", fmt.Errorf("path is required")
	}
	cmdArgs := append([]string{"addexclusion"}, args...)
	output, err := run(cmdArgs...)
	if err != nil {
		return "", err
	}
	if output == "" {
		return fmt.Sprintf("Exclusion added for %s.", args[0]), nil
	}
	return output, nil
}

// RemoveExclusion removes an exclusion for an item.
func RemoveExclusion(args []string) (string, error) {
	if len(args) == 0 || args[0] == "" {
		return "", fmt.Errorf("path is required")
	}
	cmdArgs := append([]string{"removeexclusion"}, args...)
	output, err := run(cmdArgs...)
	if err != nil {
		return "", err
	}
	if output == "" {
		return fmt.Sprintf("Exclusion removed for %s.", args[0]), nil
	}
	return output, nil
}

// IsExcluded checks if one or more items are excluded from backup.
func IsExcluded(args []string) (string, error) {
	if len(args) == 0 || args[0] == "" {
		return "", fmt.Errorf("path is required")
	}
	cmdArgs := append([]string{"isexcluded"}, args...)
	return run(cmdArgs...)
}
