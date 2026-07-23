#!/usr/bin/env bash
set -euo pipefail

tag="${2:-}"
if [[ ! "$tag" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo 'Usage: copr_script.sh OUTPUT_DIRECTORY vMAJOR.MINOR.PATCH' >&2
    exit 1
fi
mkdir -p -- "${1:?missing output directory}"
outdir="$(realpath -- "$1")"
cd -- "$(dirname -- "$0")/.."
version="${tag#v}"
echo "Preparing Artalk $tag"

# Both sources come from the requested release. Never fall back to latest.
./vendor-tarball.sh "$tag"
mv -- "artalk-$version-vendored.tar.gz" "$outdir/"
cp -- artalk.spec artalk.service artalk.sysusers vendor-tarball.sh \
    0002-remove-upgrade-command.patch "$outdir/"
# Keep packaging revisions for the same version; start other versions at 1.
if [[ "$version" != "$(awk '/^Version:/ {print $2}' artalk.spec)" ]]; then
    sed -i -E 's/^Release:[[:space:]]+.*/Release:        1%{?dist}/' "$outdir/artalk.spec"
fi
sed -i -E "s/^Version:[[:space:]]+.*/Version:        $version/" "$outdir/artalk.spec"

release="https://github.com/ArtalkJS/Artalk/releases/download/$tag"
curl_args=(--fail --silent --show-error --location --proto '=https' --proto-redir '=https'
    --connect-timeout 15 --max-time 300 --retry 3)
curl "${curl_args[@]}" "$release/artalk_ui.tar.gz" -o "$outdir/artalk_ui.tar.gz"
metadata="$(curl "${curl_args[@]}" "https://api.github.com/repos/ArtalkJS/Artalk/releases/tags/$tag")"
cd -- "$outdir"
checksum="$(jq -er '.assets[] | select(.name == "artalk_ui.tar.gz") | .digest | strings | sub("^sha256:"; "")' <<< "$metadata")"
if [[ ! "$checksum" =~ ^[0-9a-fA-F]{64}$ ]]; then
    echo 'Expected exactly one SHA-256 checksum for artalk_ui.tar.gz.' >&2
    exit 1
fi
printf '%s  artalk_ui.tar.gz\n' "$checksum" | sha256sum --check -
