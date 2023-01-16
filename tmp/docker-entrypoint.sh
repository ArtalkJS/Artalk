#!/bin/bash
set -e

if [ ! -e /data/artalk-go.yml ] && [ "$1" != 'gen' ]; then
    if [ ! -e /conf.yml ]; then
        # Generate new config
        artalk-go gen conf /data/artalk-go.yml
        echo "$(date) [info] Generate new config file to '/data/artalk-go.yml'"
    else
        # Move original config to `/data/` for upgrade (<= v2.1.8)
        cp /conf.yml /data/artalk-go.yml
        upMsg=""
        upMsg+=$'# [v2.1.9+ Updated]\n'
        upMsg+=$'# The new version of the ArtalkGo container recommends mounting\n'
        upMsg+=$'# an entire folder instead of a single file to avoid some issues.\n'
        upMsg+=$'#\n'
        upMsg+=$'# The original config file has been moved to the "/data/" folder,\n'
        upMsg+=$'# please unmount the config file volume from your container\n'
        upMsg+=$'# and edit "/data/artalk-go.yml" for configuration.'
        echo "$upMsg" > /conf.yml
        echo "$(date) [info] Copy config file from '/conf.yml' to '/data/artalk-go.yml' for upgrade"
    fi
fi

# Run ArtalkGo
artalk-go "$@"
