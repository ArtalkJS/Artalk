#!/bin/bash

set -e

if ! command -v pnpm &> /dev/null
then
    apt-get update && apt-get install --no-install-recommends -y -q curl ca-certificates

    # Install volta
    bash -c "$(curl -fsSL https://get.volta.sh)" -- --skip-setup
    export VOLTA_HOME="${HOME}/.volta"
    export PATH="${VOLTA_HOME}/bin:${PATH}"

    volta install node@20.13.1
    volta install pnpm@9.1.2
fi

pnpm install --frozen-lockfile
pnpm build:all

## dist folders
DIST_DIR="./public/dist"
SIDEBAR_DIR="./public/sidebar"

## dist
rm -rf ${DIST_DIR} && mkdir -p ${DIST_DIR}
cp -r ./ui/artalk/dist/{Artalk.css,Artalk.js} ${DIST_DIR}
cp -r ./ui/artalk/dist/{ArtalkLite.css,ArtalkLite.js} ${DIST_DIR}

I18N_DIR="${DIST_DIR}/i18n"
mkdir -p ${I18N_DIR}
cp -r ./ui/artalk/dist/i18n/*.js ${I18N_DIR}

## sidebar
rm -rf ${SIDEBAR_DIR} && mkdir -p ${SIDEBAR_DIR}
cp -r ./ui/artalk-sidebar/dist/* ${SIDEBAR_DIR}

## plugins
PLUGIN_DIR="${DIST_DIR}/plugins"
mkdir -p ${PLUGIN_DIR}

pnpm build:auth
cp -r ./ui/plugin-auth/dist/artalk-plugin-auth.js ${PLUGIN_DIR}

## create tarball for release
mkdir -p ./local/artalk_ui
cp -r ${DIST_DIR} ${SIDEBAR_DIR} ./local/artalk_ui

mkdir -p ./local/release_includes
tar -czf ./local/release_includes/artalk_ui.tar.gz -C ./local artalk_ui
