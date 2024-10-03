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

# fix for Please tell me who you are.
git config --global user.email "nonexist@artalk.com"
git config --global user.name "artalk"

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
mkdir -p result_srpm
rpkg srpm --spec artalk.spec --outdir=result_srpm
cp -r result_srpm/artalk*.src.rpm $outdir
popd
