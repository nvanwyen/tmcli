//
// command.go
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

package ui

import "tmcli/tmutil"

// InputField describes a text input field for a parameterized command.
type InputField struct {
	Label       string
	Placeholder string
	Required    bool
}

// Command describes a single tmutil command exposed in the TUI and CLI.
type Command struct {
	ID          string                              // CLI subcommand name
	Title       string                              // TUI display title
	Description string                              // detailed help text
	Hotkey      string                              // TUI hotkey
	Execute     func(args []string) (string, error) // run the command
	Inputs       []InputField                        // nil = no args needed
	IsMonitor    bool                                // special monitor mode
	RequiresRoot bool                                // needs root/sudo
}

// Category groups related commands for the TUI submenu.
type Category struct {
	Title    string
	Hotkey   string
	Commands []Command
}

// noArgs wraps a zero-argument function into the standard args signature.
func noArgs(fn func() (string, error)) func([]string) (string, error) {
	return func([]string) (string, error) { return fn() }
}

// Categories returns all command categories for the TUI.
func Categories() []Category {
	return []Category{
		{
			Title:  "Backup",
			Hotkey: "b",
			Commands: []Command{
				{ID: "start", Title: "Start", Hotkey: "s", Execute: noArgs(tmutil.StartBackup), RequiresRoot: true,
					Description: "Begin a new Time Machine backup. The backup runs in the background and backs up all volumes that are configured for backup to the default or specified destination. Requires root privileges."},
				{ID: "stop", Title: "Stop", Hotkey: "t", Execute: noArgs(tmutil.StopBackup), RequiresRoot: true,
					Description: "Stop a currently running Time Machine backup. If no backup is in progress, this command has no effect. Requires root privileges."},
				{ID: "status", Title: "Status", Hotkey: "a", Execute: noArgs(tmutil.Status),
					Description: "Display the current status of Time Machine. Shows whether a backup is running, the current phase, percent complete, bytes and files copied, time remaining, and the destination volume."},
				{ID: "monitor", Title: "Monitor", Hotkey: "m", IsMonitor: true,
					Description: "Open a live progress monitor that polls Time Machine status every second. Displays a progress bar, bytes/files copied, time remaining, and elapsed time. Updates in real time until the backup completes or you exit."},
				{ID: "enable", Title: "Enable", Hotkey: "e", Execute: noArgs(tmutil.Enable), RequiresRoot: true,
					Description: "Enable automatic Time Machine backups. When enabled, macOS will automatically perform periodic backups according to the system schedule. Requires root privileges."},
				{ID: "disable", Title: "Disable", Hotkey: "d", Execute: noArgs(tmutil.Disable), RequiresRoot: true,
					Description: "Disable automatic Time Machine backups. Prevents macOS from performing scheduled backups. Manual backups can still be started with the start command. Requires root privileges."},
				{ID: "version", Title: "Version", Hotkey: "v", Execute: noArgs(tmutil.Version),
					Description: "Display the version of the tmutil command-line utility installed on this system."},
			},
		},
		{
			Title:  "Destinations",
			Hotkey: "d",
			Commands: []Command{
				{ID: "destinationinfo", Title: "Destination Info", Hotkey: "i", Execute: noArgs(tmutil.DestinationInfo),
					Description: "Display detailed information about configured backup destinations, including the destination name, kind (local/network), mount point, and unique destination ID."},
				{ID: "setdestination", Title: "Set Destination", Hotkey: "s", Execute: tmutil.SetDestination, RequiresRoot: true, Inputs: []InputField{
					{Label: "Mount Point", Placeholder: "/Volumes/backup", Required: true},
				}, Description: "Set the backup destination to the specified mount point. Use the -a flag to add a destination rather than replacing the current one. For network destinations, use an AFP URL. Requires root privileges."},
				{ID: "removedestination", Title: "Remove Destination", Hotkey: "r", Execute: tmutil.RemoveDestination, RequiresRoot: true, Inputs: []InputField{
					{Label: "Destination ID", Placeholder: "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX", Required: true},
				}, Description: "Remove a backup destination by its unique ID. Use 'destinationinfo' to find the ID of the destination you want to remove. Requires root privileges."},
				{ID: "setquota", Title: "Set Quota", Hotkey: "q", Execute: tmutil.SetQuota, RequiresRoot: true, Inputs: []InputField{
					{Label: "Destination ID", Placeholder: "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX", Required: true},
					{Label: "Quota (GB)", Placeholder: "500", Required: true},
				}, Description: "Set a storage quota in gigabytes for a specific backup destination. This limits how much space Time Machine will use on that destination. Use 'destinationinfo' to find the destination ID."},
			},
		},
		{
			Title:  "Snapshots",
			Hotkey: "s",
			Commands: []Command{
				{ID: "localsnapshot", Title: "Create Snapshot", Hotkey: "c", Execute: noArgs(tmutil.LocalSnapshot),
					Description: "Create a new local APFS snapshot on the boot volume. Local snapshots are lightweight, point-in-time copies of your data stored on the same disk. They provide quick recovery without needing a backup destination."},
				{ID: "listlocalsnapshots", Title: "List Snapshots", Hotkey: "l", Execute: tmutil.ListLocalSnapshots, Inputs: []InputField{
					{Label: "Mount Point", Placeholder: "/ (default)"},
				}, Description: "List all local Time Machine snapshots for a given mount point. Defaults to the root volume (/) if no mount point is specified. Shows snapshot identifiers with timestamps."},
				{ID: "listlocalsnapshotdates", Title: "List Snapshot Dates", Hotkey: "d", Execute: tmutil.ListLocalSnapshotDates, Inputs: []InputField{
					{Label: "Mount Point", Placeholder: "/ (default)"},
				}, Description: "List the dates of all local Time Machine snapshots for a given mount point. Provides a concise date-only view of available snapshots. Defaults to root volume (/) if not specified."},
				{ID: "deletelocalsnapshots", Title: "Delete Snapshots", Hotkey: "x", Execute: tmutil.DeleteLocalSnapshots, RequiresRoot: true, Inputs: []InputField{
					{Label: "Mount Point or Date", Placeholder: "/ or 2026-02-07", Required: true},
				}, Description: "Delete local Time Machine snapshots. Specify either a mount point to delete all snapshots on that volume, or a specific snapshot date to delete a single snapshot. Useful for reclaiming disk space. Requires root privileges."},
				{ID: "thinlocalsnapshots", Title: "Thin Snapshots", Hotkey: "t", Execute: tmutil.ThinLocalSnapshots, RequiresRoot: true, Inputs: []InputField{
					{Label: "Mount Point", Placeholder: "/", Required: true},
					{Label: "Purge Amount (bytes)", Placeholder: "(optional)"},
					{Label: "Urgency (1-4)", Placeholder: "(optional)"},
				}, Description: "Thin (reduce) local snapshots for a mount point to free disk space. Optionally specify a purge amount in bytes and an urgency level (1=low to 4=high). Higher urgency levels delete more aggressively. Requires root privileges."},
			},
		},
		{
			Title:  "Exclusions",
			Hotkey: "e",
			Commands: []Command{
				{ID: "addexclusion", Title: "Add Exclusion", Hotkey: "a", Execute: tmutil.AddExclusion, Inputs: []InputField{
					{Label: "Path", Placeholder: "/path/to/exclude", Required: true},
				}, Description: "Add a fixed-path exclusion so Time Machine will skip the specified file or directory during backups. The exclusion is tied to the exact path and persists across backups. Useful for excluding large build artifacts, caches, or temporary files."},
				{ID: "removeexclusion", Title: "Remove Exclusion", Hotkey: "r", Execute: tmutil.RemoveExclusion, Inputs: []InputField{
					{Label: "Path", Placeholder: "/path/to/include", Required: true},
				}, Description: "Remove a previously added exclusion, allowing Time Machine to back up the specified path again. The path must match the one used when the exclusion was added."},
				{ID: "isexcluded", Title: "Check Exclusion", Hotkey: "e", Execute: tmutil.IsExcluded, Inputs: []InputField{
					{Label: "Path", Placeholder: "/path/to/check", Required: true},
				}, Description: "Check whether a file or directory is excluded from Time Machine backups. Reports whether the item is included or excluded, and whether the exclusion is fixed-path or volume-based."},
			},
		},
		{
			Title:  "Browse",
			Hotkey: "r",
			Commands: []Command{
				{ID: "latestbackup", Title: "Latest Backup", Hotkey: "l", Execute: noArgs(tmutil.LatestBackup),
					Description: "Display the path to the most recent completed Time Machine backup. This is the newest backup snapshot available for restoration."},
				{ID: "listbackups", Title: "List Backups", Hotkey: "b", Execute: noArgs(tmutil.ListBackups),
					Description: "List the paths of all completed Time Machine backup snapshots on all mounted backup destinations. Each entry represents a point-in-time backup that can be browsed or restored from."},
				{ID: "machinedirectory", Title: "Machine Directory", Hotkey: "m", Execute: noArgs(tmutil.MachineDirectory),
					Description: "Display the path to the machine-specific backup directory on the backup destination. This is the top-level directory that contains all backups for this computer."},
				{ID: "compare", Title: "Compare", Hotkey: "c", Execute: tmutil.Compare, Inputs: []InputField{
					{Label: "Path 1", Placeholder: "/path/one (optional)"},
					{Label: "Path 2", Placeholder: "/path/two (optional)"},
				}, Description: "Compare the current system state to a backup, or compare two paths. With no arguments, compares the live system to the latest backup. With one path, compares to that backup snapshot. With two paths, compares them directly. Reports added, removed, and changed files."},
				{ID: "uniquesize", Title: "Unique Size", Hotkey: "u", Execute: tmutil.UniqueSize, Inputs: []InputField{
					{Label: "Path", Placeholder: "/path/to/check", Required: true},
				}, Description: "Calculate the unique disk space consumed by a specific backup path. Shows how much space would be freed if that backup were deleted, accounting for data shared with other backups via hard links."},
				{ID: "verifychecksums", Title: "Verify Checksums", Hotkey: "v", Execute: tmutil.VerifyChecksums, Inputs: []InputField{
					{Label: "Path", Placeholder: "/path/to/verify", Required: true},
				}, Description: "Verify the integrity of backed-up files by checking their stored checksums. Reads each file in the specified backup path and compares its checksum to the stored value. Reports any corrupted files."},
			},
		},
		{
			Title:  "Restore",
			Hotkey: "t",
			Commands: []Command{
				{ID: "findfile", Title: "Find File", Hotkey: "f", Execute: tmutil.FindFile, Inputs: []InputField{
					{Label: "Filename / Pattern", Placeholder: "*.txt or myfile.doc", Required: true},
					{Label: "Max Backups to Search", Placeholder: "5 (default)"},
				}, Description: "Search for a file or directory by name across recent Time Machine backup snapshots. Uses glob pattern matching against file basenames. Searches from the most recent backup backward, limited to a configurable number of snapshots (default 5) for performance. Results show full paths that can be used with the Restore command."},
				{ID: "findbydate", Title: "Find by Date", Hotkey: "d", Execute: tmutil.FindByDate, Inputs: []InputField{
					{Label: "Start Date", Placeholder: "2026-01-01", Required: true},
					{Label: "End Date", Placeholder: "2026-02-07 (default: today)"},
				}, Description: "List available Time Machine backup snapshots within a date range. Dates are parsed from backup path names. Provide a start date in YYYY-MM-DD format. The end date is optional and defaults to today. Useful for finding which backups cover a specific time period before restoring."},
				{ID: "browsebackup", Title: "Browse Backup", Hotkey: "w", Execute: tmutil.BrowseBackup, Inputs: []InputField{
					{Label: "Backup Path", Placeholder: "/Volumes/Backup/Backups.backupdb/Mac/2026-02-07-143022", Required: true},
					{Label: "Subdirectory", Placeholder: "Users/name/Documents (optional)"},
				}, Description: "List the contents of a specific Time Machine backup snapshot directory. Provide the full backup path (from 'List Backups' or 'Find by Date') and optionally a subdirectory within it. Shows files and directories with sizes, useful for identifying what to restore."},
				{ID: "restore", Title: "Restore File", Hotkey: "r", Execute: tmutil.Restore, RequiresRoot: true, Inputs: []InputField{
					{Label: "Source Path", Placeholder: "/backup/path/file", Required: true},
					{Label: "Destination Path", Placeholder: "/restore/to/here", Required: true},
				}, Description: "Restore files or directories from a Time Machine backup to a specified destination. Copies files from the backup source path to the destination with verbose output. The source should be a path within a backup snapshot. Requires root privileges."},
			},
		},
		{
			Title:  "Advanced",
			Hotkey: "a",
			Commands: []Command{
				{ID: "delete", Title: "Delete Backup", Hotkey: "d", Execute: tmutil.Delete, RequiresRoot: true, Inputs: []InputField{
					{Label: "Arguments", Placeholder: "-d mount_point -t timestamp  or  -p path", Required: true},
				}, Description: "Delete a specific backup snapshot. Use '-d mount_point -t timestamp' to delete by destination and time, or '-p path' to delete by path. This permanently removes the backup data and cannot be undone. Requires root privileges."},
				{ID: "associatedisk", Title: "Associate Disk", Hotkey: "a", Execute: tmutil.AssociateDisk, RequiresRoot: true, Inputs: []InputField{
					{Label: "Mount Point", Placeholder: "/Volumes/disk", Required: true},
					{Label: "Volume Backup Dir", Placeholder: "/path/to/backup/dir", Required: true},
				}, Description: "Associate a volume with a backup directory when a disk has been reformatted or replaced. This tells Time Machine that the specified volume corresponds to the given backup directory, allowing backups to continue without starting from scratch. Requires root privileges."},
				{ID: "inheritbackup", Title: "Inherit Backup", Hotkey: "i", Execute: tmutil.InheritBackup, RequiresRoot: true, Inputs: []InputField{
					{Label: "Machine Dir or Sparse Bundle", Placeholder: "/path/to/machine_dir", Required: true},
				}, Description: "Claim ownership of a machine directory or sparse bundle from another computer. Allows this machine to continue backing up to an existing backup set, useful when migrating to new hardware. Requires root privileges."},
				{ID: "calculatedrift", Title: "Calculate Drift", Hotkey: "c", Execute: tmutil.CalculateDrift, Inputs: []InputField{
					{Label: "Machine Directory", Placeholder: "/path/to/machine_dir", Required: true},
				}, Description: "Analyze a machine backup directory and calculate the drift (differences) between backup snapshots. Useful for diagnosing backup performance issues or understanding what changed between backups."},
				{ID: "deleteinprogress", Title: "Delete In Progress", Hotkey: "p", Execute: tmutil.DeleteInProgress, RequiresRoot: true, Inputs: []InputField{
					{Label: "Machine Directory", Placeholder: "/path/to/machine_dir", Required: true},
				}, Description: "Delete an incomplete or failed in-progress backup from the specified machine directory. Use this to clean up after a backup that was interrupted or failed partway through. Requires root privileges."},
			},
		},
	}
}

// AllCommands returns a flat list of all commands across all categories.
func AllCommands() []Command {
	var cmds []Command
	for _, cat := range Categories() {
		cmds = append(cmds, cat.Commands...)
	}
	return cmds
}

// FindCommand looks up a command by its CLI ID.
func FindCommand(id string) *Command {
	for _, cat := range Categories() {
		for i := range cat.Commands {
			if cat.Commands[i].ID == id {
				return &cat.Commands[i]
			}
		}
	}
	return nil
}
