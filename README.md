# Time Machine CLI (tmcli)

Version: v1.0.2

A macOS Time Machine CLI and interactive TUI built in Go. Wraps the native
`tmutil` utility with a friendlier interface, hotkey-driven menus, a live
backup monitor, and direct command-line access to every operation.

## Personal Notes
_This project represents my first dive into the Go language. Since it is my
first, I do not expect it to be perfect or follow all the "industry standards"
accepted by Professional Go Developers. So, if you have comments about how it's
written, organized and formatted please be kind with your comments. I am more
than happy to listen and take constructive criticism, but I will not listen to
rude comments particularly those that add nothing constructive._

### Motive
This project was written because, I wanted to learn Go, but more becasue I
wanted a command line TUI that wrapped up the Mac Time Machine. I wanted a
statusbar monitor, a way to start, stop, restore etc... _(those features
provided by the Mac tmutil)_ that I could call from the terminal, because that's
where I primarly live in my day-to-day job. I made this project because I
couldn't find another wrapper that would provide a simple TUI interface or a
monitoring system that was easy to use and undertand. I made it a public project
with the idea that I might not be the only one with this request.

### Bugs, reporting and requests
Feel free to add bugs and comments in the project, and if you have any feature
requests I'll be happy to consider them. If you clone the project and have merge
requests, I'll also consider them. However, please note that this is not a
primary project for me, and it may take awhile for me to get around to your
requests. I'm not ignoring you, I'm just busy with other things, so please be
patient with me.


## Features

- **Interactive TUI** — category menus, hotkey navigation, input forms, scrollable output
- **Direct CLI** — run any command as a subcommand (`tmcli status`, `tmcli listbackups`, etc.)
- **Live monitor** — real-time progress bar with bytes/files copied, ETA, and elapsed time
- **Built-in help** — browse detailed descriptions, parameters, and CLI usage for every command
- **Root awareness** — commands that require `sudo` are clearly marked in the TUI

## Requirements

- macOS (Time Machine / `tmutil` must be available)
- Go 1.25.7+
- CMake 3.20+

## Building

Use the provided `configure` script:

```bash
# Release build (default)
./configure

# Debug build
./configure --debug

# Build and install to /usr/local/bin
./configure --install

# Build and install to a custom prefix
./configure --install --prefix /opt/bin

# Clean build artifacts
./configure --clean

# Uninstall
./configure --uninstall
```

The compiled binary is placed in `bin/tmcli`.

## Usage

### Interactive TUI

Launch with no arguments to enter the interactive terminal UI:

```bash
tmcli
```

Navigate with arrow keys or hotkeys, press `enter` to select, `esc` to go
back, and `q` to quit.

### CLI Mode

Run any command directly from the shell:

```bash
tmcli <command> [arguments]
```

Use `--help` to see all available commands:

```bash
tmcli --help
```

## Commands

### General

| Flag              | Description                          | Example            |
|-------------------|--------------------------------------|--------------------|
| `-v`, `--version` | Print the version and exit           | `tmcli --version`  |
| `-h`, `--help`    | Print usage information and exit     | `tmcli --help`     |

### Backup

| Command   | Description                          | Root | Example                 |
|-----------|--------------------------------------|------|-------------------------|
| `start`   | Start a Time Machine backup          | yes  | `sudo tmcli start`      |
| `stop`    | Stop a running backup                | yes  | `sudo tmcli stop`       |
| `status`  | Show current backup status           | no   | `tmcli status`          |
| `monitor` | Live progress monitor                | no   | `tmcli monitor`         |
| `enable`  | Enable automatic backups             | yes  | `sudo tmcli enable`     |
| `disable` | Disable automatic backups            | yes  | `sudo tmcli disable`    |
| `version` | Show tmutil version                  | no   | `tmcli version`         |

### Destinations

| Command             | Description                        | Root | Example                                                        |
|---------------------|------------------------------------|------|----------------------------------------------------------------|
| `destinationinfo`   | Show destination details           | no   | `tmcli destinationinfo`                                        |
| `setdestination`    | Set backup destination             | yes  | `sudo tmcli setdestination /Volumes/Backup`                    |
| `removedestination` | Remove a destination by ID         | yes  | `sudo tmcli removedestination XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX` |
| `setquota`          | Set storage quota (GB)             | yes  | `sudo tmcli setquota XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX 500` |

### Snapshots

| Command                  | Description                     | Root | Example                                   |
|--------------------------|---------------------------------|------|--------------------------------------------|
| `localsnapshot`          | Create a local APFS snapshot    | no   | `tmcli localsnapshot`                      |
| `listlocalsnapshots`     | List local snapshots            | no   | `tmcli listlocalsnapshots /`               |
| `listlocalsnapshotdates` | List snapshot dates             | no   | `tmcli listlocalsnapshotdates /`           |
| `deletelocalsnapshots`   | Delete snapshots by date/mount  | yes  | `sudo tmcli deletelocalsnapshots 2026-02-07` |
| `thinlocalsnapshots`     | Thin snapshots to free space    | yes  | `sudo tmcli thinlocalsnapshots / 1000000000 2` |

### Exclusions

| Command           | Description                          | Root | Example                                  |
|-------------------|--------------------------------------|------|------------------------------------------|
| `addexclusion`    | Exclude a path from backups          | no   | `tmcli addexclusion /path/to/exclude`    |
| `removeexclusion` | Remove an exclusion                  | no   | `tmcli removeexclusion /path/to/include` |
| `isexcluded`      | Check if a path is excluded          | no   | `tmcli isexcluded /path/to/check`        |

### Browse

| Command            | Description                         | Root | Example                              |
|--------------------|-------------------------------------|------|--------------------------------------|
| `latestbackup`     | Show most recent backup path        | no   | `tmcli latestbackup`                 |
| `listbackups`      | List all completed backups          | no   | `tmcli listbackups`                  |
| `machinedirectory` | Show machine backup directory       | no   | `tmcli machinedirectory`             |
| `compare`          | Compare system to backup            | no   | `tmcli compare`                      |
| `uniquesize`       | Calculate unique size of a backup   | no   | `tmcli uniquesize /path/to/backup`   |
| `verifychecksums`  | Verify backup file integrity        | no   | `tmcli verifychecksums /path/to/backup` |

### Restore

| Command        | Description                              | Root | Example                                                         |
|----------------|------------------------------------------|------|-----------------------------------------------------------------|
| `findfile`     | Search for a file across backups         | no   | `tmcli findfile "*.txt" 10`                                     |
| `findbydate`   | List backups within a date range         | no   | `tmcli findbydate 2026-01-01 2026-02-07`                       |
| `browsebackup` | List contents of a backup snapshot       | no   | `tmcli browsebackup /Volumes/Backup/Backups.backupdb/Mac/2026-02-07-143022` |
| `restore`      | Restore files from a backup              | yes  | `sudo tmcli restore /backup/path/file /restore/to/here`        |

### Advanced

| Command            | Description                           | Root | Example                                                   |
|--------------------|---------------------------------------|------|------------------------------------------------------------|
| `delete`           | Delete a backup snapshot              | yes  | `sudo tmcli delete -d /Volumes/Backup -t 2026-02-07-143022` |
| `associatedisk`    | Associate a volume with a backup dir  | yes  | `sudo tmcli associatedisk /Volumes/disk /path/to/backup`  |
| `inheritbackup`    | Claim a backup from another machine   | yes  | `sudo tmcli inheritbackup /path/to/machine_dir`           |
| `calculatedrift`   | Analyze drift between backups         | no   | `tmcli calculatedrift /path/to/machine_dir`               |
| `deleteinprogress` | Delete an incomplete backup           | yes  | `sudo tmcli deleteinprogress /path/to/machine_dir`        |

## TUI Navigation

| Key            | Action                        |
|----------------|-------------------------------|
| `Up` / `k`     | Move cursor up                |
| `Down` / `j`   | Move cursor down              |
| `Enter`        | Select item                   |
| `Esc` / `Backspace` | Go back                 |
| `h`            | Open help                     |
| `q`            | Quit                          |
| `Tab`          | Next input field              |
| `Shift+Tab`    | Previous input field          |
| `PgUp` / `PgDn` | Scroll output pages         |

## License

Copyright (c) 2004-2026 Metasystems Technologies Inc. (MTI). All rights reserved.

Distributed under the MTI Software License, Version 0.1.
