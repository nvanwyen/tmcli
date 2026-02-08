#
# CMakeLists.txt
# ~~~~~~~~~~~~~~~~~~~~~
#
# Copyright (c) 2004-2026 Metasystems Technologies Inc. (MTI)
# All rights reserved
#
# Distributed under the MTI Software License, Version 0.1.
#
# as defined by accompanying file MTI-LICENSE-0.1.info or
# at http://www.mtihq.com/license/MTI-LICENSE-0.1.info
#

cmake_minimum_required(VERSION 3.20)
project(tmcli LANGUAGES NONE)

list(APPEND CMAKE_MODULE_PATH "${CMAKE_SOURCE_DIR}/cmake")

# ensure out-of-source builds (error when not out-of-source)
include(MacroEnsureOutOfSourceBuild)
EnsureOutOfSourceBuild( "
    *** IMPORTANT *** 
    Builds can *ONLY* be out-of-source!
    Make sure to delete the ./CMakeCache.txt file *before* retrying
    You can use the supplied ./configure.sh script, to build the
    project with release, debug and in verbose modes use:
        ./configure.sh --help
    for a full list of options
    *** IMPORTANT *** 
    " )

find_package(Go REQUIRED)

message(STATUS "Go compiler: ${GO_EXECUTABLE}")
message(STATUS "Go version:  ${GO_VERSION}")

set(BINARY_NAME tmcli)
set(BINARY_OUTPUT_DIR "${CMAKE_SOURCE_DIR}/bin")
set(SOURCE_DIR "${CMAKE_SOURCE_DIR}/source")

# Read version from VERSION file
file(STRINGS "${SOURCE_DIR}/VERSION" TMCLI_VERSION)
message(STATUS "Version:     ${TMCLI_VERSION}")

# Generate version.go from template
configure_file(
    "${SOURCE_DIR}/version.go.in"
    "${SOURCE_DIR}/version.go"
    @ONLY
)

# Set Go build flags based on build type
if(CMAKE_BUILD_TYPE STREQUAL "Debug")
    set(GO_BUILD_FLAGS -gcflags "all=-N -l")
    message(STATUS "Build type:  Debug (symbols enabled, no optimization)")
else()
    set(GO_BUILD_FLAGS -ldflags "-s -w")
    message(STATUS "Build type:  Release (stripped, optimized)")
endif()

file(MAKE_DIRECTORY ${BINARY_OUTPUT_DIR})

# Collect source files for dependency tracking
file(GLOB_RECURSE GO_SOURCES "${SOURCE_DIR}/*.go")

add_custom_target(${BINARY_NAME} ALL
    COMMAND ${CMAKE_COMMAND} -E env
        GOPATH=${CMAKE_BINARY_DIR}/go
        ${GO_EXECUTABLE} build ${GO_BUILD_FLAGS} -o ${BINARY_OUTPUT_DIR}/${BINARY_NAME} .
    WORKING_DIRECTORY ${SOURCE_DIR}
    SOURCES ${GO_SOURCES}
    COMMENT "Building ${BINARY_NAME} (${CMAKE_BUILD_TYPE})"
)

add_custom_target(tidy
    COMMAND ${GO_EXECUTABLE} mod tidy
    WORKING_DIRECTORY ${SOURCE_DIR}
    COMMENT "Running go mod tidy"
)

add_custom_target(clean-bin
    COMMAND ${CMAKE_COMMAND} -E remove -f ${BINARY_OUTPUT_DIR}/${BINARY_NAME}
    COMMENT "Cleaning binary output"
)

# Install target
set(INSTALL_PREFIX "/usr/local/bin" CACHE PATH "Installation directory")

add_custom_target(install-bin
    COMMAND ${CMAKE_COMMAND} -E copy ${BINARY_OUTPUT_DIR}/${BINARY_NAME} ${INSTALL_PREFIX}/${BINARY_NAME}
    DEPENDS ${BINARY_NAME}
    COMMENT "Installing ${BINARY_NAME} to ${INSTALL_PREFIX}"
)

add_custom_target(uninstall-bin
    COMMAND ${CMAKE_COMMAND} -E remove -f ${INSTALL_PREFIX}/${BINARY_NAME}
    COMMENT "Uninstalling ${BINARY_NAME} from ${INSTALL_PREFIX}"
)
