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

git -c advice.detachedHead=false clone --branch $tag --depth 1 https://github.com/artalkjs/artalk/ $name-$version
pushd $name-$version
GOPROXY='https://proxy.golang.org,direct' go mod vendor
popd
tar --exclude .git -czf $name-$version-vendored.tar.gz $name-$version
rm -rf $name-$version
