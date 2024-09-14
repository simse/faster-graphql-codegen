set shell := ["bash", "-eu", "-o", "pipefail", "-c"]

build:
    #!/usr/bin/env bash

    # Determine version
    if [ "$(git rev-parse --abbrev-ref HEAD)" = "main" ]; then
        version="latest"
    elif git describe --tags --exact-match >/dev/null 2>&1; then
        version="$(git describe --tags --exact-match)"
    else
        version="$(git rev-parse --short HEAD)"
    fi

    echo "Building version: $version"

    platforms=("linux" "windows" "darwin")
    archs=("amd64" "arm64")

    for GOOS in "${platforms[@]}"; do
        for GOARCH in "${archs[@]}"; do
            output="fgc-${GOOS}-${GOARCH}-${version}"
            if [ "$GOOS" = "windows" ]; then
                output+=".exe"
            fi
            echo "Building for $GOOS/$GOARCH: $output"

            # Build the binary
            env CGO_ENABLED=0 GOOS="$GOOS" GOARCH="$GOARCH" \
                go build -ldflags="-s -w" -o "build/$output" .
        done
    done

lint:
    staticcheck ./...

test:
    go test ./...

format:
    go fmt ./...

docs-dev:
    pnpm -C docs docs:dev