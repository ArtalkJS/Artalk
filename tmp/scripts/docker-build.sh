#!/bin/bash

set -e

IMAGE_NAME="artalk-go"
REPO_NAME="artalk/artalk-go"

VERSION=$(git describe --tags --abbrev=0)

if [[ $* == *--push* ]]
then
    # tag and push image
    for tag in {${VERSION},latest}; do
        docker image tag ${IMAGE_NAME} ${REPO_NAME}:${tag}
        docker push ${REPO_NAME}:${tag}
    done
else
    # build
    docker image build -t ${IMAGE_NAME} .
fi
