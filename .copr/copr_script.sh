#!/usr/bin/env bash

set -euo pipefail

# Copr needs the SRPM in the requested output directory.
outdir="${1:?missing output directory}"
mkdir -p -- "$outdir"
outdir="$(realpath -- "$outdir")"
echo "Outdir: $outdir"

# Prefer the exact release supplied by the custom webhook. Fall back to the
# latest release tag so manual COPR rebuilds remain supported.
tag=""
if [[ -f hook_payload ]]; then
    tag="$(jq -r '.version // empty' hook_payload)"
fi
if [[ -z "$tag" ]]; then
    tag="$(git tag --list 'v*' --sort=-version:refname | sed -n '1p')"
fi
if [[ ! "$tag" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "Invalid release tag: $tag" >&2
    exit 1
fi
git rev-parse --quiet --verify "refs/tags/$tag^{commit}" > /dev/null
echo "Using tag: $tag"

version="${tag#v}"
sed -i -E "s/^Version:[[:space:]]+.*/Version:        $version/" artalk.spec
grep --fixed-strings "Version:        $version" artalk.spec

# Create the vendored source and download the release UI archive.
rm -f -- "artalk-$version-vendored.tar.gz"
./vendor-tarball.sh "$tag"
spectool -g artalk.spec

# Parse %autorelease and %autochangelog.
rpmautospec process-distgit ./artalk.spec ./artalk.spec

# Generate the SRPM.
rm -rf -- result_srpm
mkdir -p -- result_srpm "$outdir"
rpkg srpm --spec artalk.spec --outdir=result_srpm
cp -v -- ./result_srpm/artalk*.src.rpm "$outdir/"
