#!/usr/bin/env bash

set -euo pipefail

tag="${1:-}"
if [[ ! "$tag" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "This script requires a stable release tag (vMAJOR.MINOR.PATCH)." >&2
    exit 1
fi

name=artalk
version=${tag#v}
echo "Using version: $version"

# Clone the repository
rm -rf "$name-$version"
git -c advice.detachedHead=false clone --branch "$tag" --depth 1 https://github.com/artalkjs/artalk/ "$name-$version"

# Vendor dependencies
pushd "$name-$version" > /dev/null
GOPROXY='https://proxy.golang.org,direct' go mod vendor
popd > /dev/null

# More reproducible!
TARFLAGS=(
  --exclude .git
  --sort=name
  --format=posix
  --pax-option=delete=atime,delete=ctime
  --clamp-mtime
  --mtime='1970-01-01 00:00:00 UTC'
  --numeric-owner
  --owner=0
  --group=0
  --mode=go+u,go-w
)

tar "${TARFLAGS[@]}" -czf "$name-$version-vendored.tar.gz" "$name-$version"

# Clean up the temporary directory
rm -rf "$name-$version"
