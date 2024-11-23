#!/bin/bash

set -e

if ! command -v pnpm &>/dev/null && command -v apt-get &>/dev/null;
then
    apt-get update && apt-get install --no-install-recommends -y -q curl ca-certificates

    # Install volta
    bash -c "$(curl -fsSL https://get.volta.sh)" -- --skip-setup
    export VOLTA_HOME="${HOME}/.volta"
    export PATH="${VOLTA_HOME}/bin:${PATH}"

    volta install node@20.17.0
    volta install pnpm@9.10.0
fi

# build
pnpm install --frozen-lockfile
pnpm build:all
pnpm build:auth

## dist folders
DIST_DIR="./public/dist"
I18N_DIR="./public/dist/i18n"
SIDEBAR_DIR="./public/sidebar"
PLUGIN_DIR="./public/dist/plugins"

# clean
rm -rf ${DIST_DIR} && mkdir -p ${DIST_DIR}
rm -rf ${I18N_DIR} && mkdir -p ${I18N_DIR}
rm -rf ${SIDEBAR_DIR} && mkdir -p ${SIDEBAR_DIR}
rm -rf ${PLUGIN_DIR} && mkdir -p ${PLUGIN_DIR}

## copy
cp ./ui/artalk/dist/{Artalk.css,Artalk.js} ${DIST_DIR}
cp ./ui/artalk/dist/{ArtalkLite.css,ArtalkLite.js} ${DIST_DIR}
cp ./ui/artalk/dist/i18n/*.js ${I18N_DIR}
cp -r ./ui/artalk-sidebar/dist/* ${SIDEBAR_DIR}
cp ./ui/plugin-*/dist/*.js ${PLUGIN_DIR}
