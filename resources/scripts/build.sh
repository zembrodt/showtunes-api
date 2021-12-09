#!/bin/bash
# Generate and build the Showtunes API server code and executable

# Set executable name
execname=showtunesapi

# Set absolute project path
projectpath="$(dirname "$( cd ../"$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )")"
echo Project path: $projectpath

# Set executable build directory
buildpath=$projectpath/build

# Make the build directory
mkdir -p $buildpath
echo Build path: $buildpath

# Build the project file
if go build ./cmd/main.go ; then
    echo main.go built successfully
else
    echo Failed to compile main package
    exit 1
fi

# Remove existing binaries
rm -f $buildpath/*
# Move the new binary
mv main $buildpath/$execname
echo Project binary: $buildpath/$execname
