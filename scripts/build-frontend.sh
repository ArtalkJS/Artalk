#!/bin/bash

set -e

if ! command -v yarn &> /dev/null
then
    apt-get update && apt-get install --no-install-recommends -y -q curl ca-certificates
    
    # Install volta
    export VOLTA_HOME="${HOME}/.volta"
    export PATH="${VOLTA_HOME}/bin:${PATH}"
    bash -c "$(curl -fsSL https://get.volta.sh)" -- --skip-setup

    volta install node
    volta install yarn
fi

mkdir -p ./local/

cd ./local/
rm -rf ./Artalk

git clone https://github.com/ArtalkJS/Artalk.git Artalk
cd Artalk

# using latest tag soruce code
git fetch --tags
git checkout $(git describe --tags --abbrev=0)

yarn install --network-timeout 1000000
yarn build:all

cd ../../

rm -rf ./frontend/dist
rm -rf ./frontend/sidebar

mkdir -p ./frontend/dist
cp -r ./local/Artalk/packages/artalk/dist/{Artalk.css,Artalk.js} ./frontend/dist

mkdir -p ./frontend/sidebar
cp -r ./local/Artalk/packages/artalk-sidebar/dist/* ./frontend/sidebar
