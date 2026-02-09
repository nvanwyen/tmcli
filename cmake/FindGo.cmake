#
# FindGo.cmake
# ~~~~~~~~~~~~~~~~~~~~~
#
# Copyright (c) 2004-2026 Metasystems Technologies Inc. (MTI)
#
# Licensed under the MIT License. See LICENSE file in the project root
# for full license text.
#

# FindGo.cmake — Locate the Go compiler
#
# This module defines:
#   GO_FOUND       - True if Go was found
#   GO_EXECUTABLE  - Path to the go binary
#   GO_VERSION     - Version of Go found

find_program(GO_EXECUTABLE
    NAMES go
    PATHS
        /usr/local/go/bin
        /usr/local/bin
        /opt/homebrew/bin
        $ENV{GOROOT}/bin
        $ENV{HOME}/go/bin
)

if(GO_EXECUTABLE)
    execute_process(
        COMMAND ${GO_EXECUTABLE} version
        OUTPUT_VARIABLE _GO_VERSION_OUTPUT
        OUTPUT_STRIP_TRAILING_WHITESPACE
    )
    string(REGEX MATCH "go([0-9]+\\.[0-9]+\\.?[0-9]*)" _GO_VERSION_MATCH "${_GO_VERSION_OUTPUT}")
    set(GO_VERSION "${CMAKE_MATCH_1}")
endif()

include(FindPackageHandleStandardArgs)
find_package_handle_standard_args(Go
    REQUIRED_VARS GO_EXECUTABLE
    VERSION_VAR GO_VERSION
)

mark_as_advanced(GO_EXECUTABLE)
