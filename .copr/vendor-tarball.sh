#!/usr/bin/env sh
tag=$1
name=artalk

if [[ -z $tag ]]; then
    echo "This script requires the tag as an argument."
    exit 1
fi

set -euo pipefail

version=${tag#v}
echo "Using version: $version"

rm -rf "$name-$version"

# Clone the repository
git -c advice.detachedHead=false clone --branch "$tag" --depth 1 https://github.com/artalkjs/artalk/ "$name-$version"
pushd "$name-$version"

# Vendor dependencies
GOPROXY='https://proxy.golang.org,direct' go mod vendor
popd

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
