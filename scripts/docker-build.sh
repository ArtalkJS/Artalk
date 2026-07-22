#!/bin/bash

set -e

IMAGE_NAME="artalk/artalk-go"
VERSION=$(git describe --tags --abbrev=0 --match 'v*')

if [[ $* == *--push* ]]
then
    # tag and push image
    for tag in {${VERSION},latest}; do
        docker image tag "${IMAGE_NAME}" "${IMAGE_NAME}:${tag}"
        docker push "${IMAGE_NAME}:${tag}"
    done
else
    # build
    docker image build \
        --build-arg "APP_VERSION=${VERSION}" \
        --build-arg "APP_COMMIT_HASH=$(git rev-parse HEAD)" \
        -t "${IMAGE_NAME}" \
        .
fi
