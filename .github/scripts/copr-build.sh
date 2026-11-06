#!/usr/bin/env bash
set -euo pipefail

if [[ ! "${VERSION:-}" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo 'Expected a stable release version (vMAJOR.MINOR.PATCH).' >&2
  exit 1
fi
if [[ "${DRY_RUN:-}" == true ]]; then
  echo "Would build Artalk $VERSION in @artalk/artalk."
  exit 0
fi
if [[ "${DRY_RUN:-}" != false ]]; then
  echo 'DRY_RUN must be true or false.' >&2
  exit 1
fi

# Only send release metadata to the existing Artalk custom webhook.
webhook="${COPR_BUILD_WEBHOOK_URL:-}"
if [[ ! "$webhook" =~ ^https://copr\.fedorainfracloud\.org/webhooks/custom/147364/[a-zA-Z0-9-]+/?$ ]]; then
  echo 'Configure COPR_BUILD_WEBHOOK_URL with the Artalk custom webhook base URL.' >&2
  exit 1
fi

curl_args=(--fail --silent --show-error --proto '=https' --connect-timeout 15 --max-time 60)
release="https://github.com/ArtalkJS/Artalk/releases/download/$VERSION"
curl "${curl_args[@]}" --location --proto-redir '=https' --head "$release/artalk_ui.tar.gz" > /dev/null

# Do not retry POST: a timed-out request might already have created a build.
payload="$(jq -cn --arg version "$VERSION" '{version: $version}')"
response="$(curl "${curl_args[@]}" --write-out '\n%{http_code}' --header 'Content-Type: application/json' \
  --data "$payload" "${webhook%/}/artalk/")"
build_id="${response%$'\n'*}"
build_id="${build_id%$'\n'}"
if [[ "${response##*$'\n'}" != 200 || ! "$build_id" =~ ^[0-9]+$ ]]; then
  echo 'COPR did not return a build ID; check webhook history before retrying.' >&2
  exit 1
fi

build_url="https://copr.fedorainfracloud.org/coprs/g/artalk/artalk/build/$build_id/"
echo "COPR build: $build_url"
if [[ -n "${GITHUB_STEP_SUMMARY:-}" ]]; then
  echo "COPR build for $VERSION: $build_url" >> "$GITHUB_STEP_SUMMARY"
fi

deadline=$((SECONDS + 6600))
while (( SECONDS < deadline )); do
  build="$(curl "${curl_args[@]}" --retry 3 \
    "https://copr.fedorainfracloud.org/api_3/build/$build_id")"
  state="$(jq -er '.state' <<< "$build")"
  echo "COPR build $build_id: $state"
  case "$state" in
    succeeded)
      jq -e --arg version "${VERSION#v}-" '
        .ownername == "@artalk" and .projectname == "artalk" and
        .source_package.name == "artalk" and
        (.source_package.version | startswith($version))
      ' <<< "$build" > /dev/null
      exit 0
      ;;
    pending|starting|running|importing|waiting) sleep 30 ;;
    *) echo "COPR build did not succeed: $build_url" >&2; exit 1 ;;
  esac
done
echo "Timed out waiting for COPR; build may still be running: $build_url" >&2
exit 1
