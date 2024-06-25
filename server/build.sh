#!/bin/bash

# Define the name of your output binary
BINARY_NAME="kuneiform-lsp"

# Define the list of OS and architecture combinations you want to build for
OS=("linux" "darwin" "windows")
ARCH=("amd64" "arm64")

# Create a build directory
mkdir -p .build

# Loop through each target and build the binary
for os in "${OS[@]}"; do
    for arch in "${ARCH[@]}"; do
        TARGET="${os}/${arch}"
        OUTPUT="./.build/${BINARY_NAME}-${os}-${arch}"

        # For Windows, add .exe extension
        if [ "$os" == "windows" ]; then
            OUTPUT="${OUTPUT}.exe"
        fi

        echo "Building for ${os}/${arch}..."
        env GOOS=$os GOARCH=$arch go build -o $OUTPUT ./...

        if [ $? -ne 0 ]; then
            echo "Failed to build for ${os}/${arch}"
            exit 1
        fi
    done
done

echo "Build completed. Binaries are in the 'server' directory."
