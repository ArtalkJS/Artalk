#!/usr/bin/env bash
set -euo pipefail

# This script is also the COPR package's Custom source script.
version="$(jq -er '.version | strings' hook_payload)"
git clone --depth 1 --branch dist/fedora https://github.com/ArtalkJS/Artalk.git packaging
packaging/.copr/copr_script.sh "$COPR_RESULTDIR" "$version"
