#!/bin/bash

rm -rf ./local/artalk_ui && mkdir -p ./local/artalk_ui
cp -r ./public/* ./local/artalk_ui

mkdir -p ./local/release_includes
tar -czf ./local/release_includes/artalk_ui.tar.gz -C ./local artalk_ui

echo $(realpath ./local/release_includes/artalk_ui.tar.gz)
