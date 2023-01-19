#!/bin/bash

set -e

if ! command -v pnpm &> /dev/null
then
    apt-get update && apt-get install --no-install-recommends -y -q curl ca-certificates

    # Install volta
    bash -c "$(curl -fsSL https://get.volta.sh)" -- --skip-setup
    export VOLTA_HOME="${HOME}/.volta"
    export PATH="${VOLTA_HOME}/bin:${PATH}"

    volta install node
    volta install pnpm
fi

pnpm --dir ./ui install --frozen-lockfile
pnpm --dir ./ui build:all

## dist
DIST_DIR="./public/dist"
rm -rf ${DIST_DIR} && mkdir -p ${DIST_DIR}
cp -r ./ui/packages/artalk/dist/{Artalk.css,Artalk.js} ${DIST_DIR}

## sidebar
SIDEBAR_DIR="./public/sidebar"
rm -rf ${SIDEBAR_DIR} && mkdir -p ${SIDEBAR_DIR}
cp -r ./ui/packages/artalk-sidebar/dist/* ${SIDEBAR_DIR}
