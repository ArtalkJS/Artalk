#!/usr/bin/env sh

set -euo pipefail

outdir=$1
echo "Outdir: $outdir"

tag=$(git tag -l 'v*' --sort=-v:refname | head -n 1)
echo "Using tag: $tag"

# Download source
spectool -g artalk.spec

rm -rf "artalk-${tag#v}-vendored.tar.gz"
./vendor-tarball.sh $tag

# Parse %autorelease and %autochangelog
rpmautospec process-distgit ./artalk.spec ./artalk.spec

# Generate srpm
mkdir -p result_srpm $outdir
rpkg srpm --spec artalk.spec --outdir=result_srpm
cp -rv result_srpm/artalk*.src.rpm $outdir
