#!/usr/bin/env sh

set -euo pipefail

outdir=$1
echo "Outdir: $outdir"

tag=$(git tag -l 'v*' --sort=-v:refname | head -n 1)
echo "Using tag: $tag"

pushd .copr

# Download source
spectool -g artalk.spec

rm -rf "artalk-${tag#v}-vendored.tar.gz"
./vendor-tarball.sh $tag

# Parse %autorelease and %autochangelog
git stash
cp changelog artalk.spec ../
git add ../changelog ../artalk.spec
git commit -m "for rpmautospec [skip changelog]"
rpmautospec process-distgit ../artalk.spec artalk.spec
if git stash list | grep -q 'stash@'; then
  git stash pop
fi

# Generate srpm
mkdir -p $outdir
rpkg srpm --spec artalk.spec --outdir=$outdir

popd
