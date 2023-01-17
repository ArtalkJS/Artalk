#!/bin/bash
set -e

if [ ! -e /data/artalk.yml ] && [ "$1" != 'gen' ]; then
    if [ ! -e /conf.yml ]; then
        # Generate new config
        artalk gen conf /data/artalk.yml
        echo "$(date) [info] Generate new config file to '/data/artalk.yml'"
    else
        # Move original config to `/data/` for upgrade (<= v2.1.8)
        cp /conf.yml /data/artalk.yml
        upMsg=""
        upMsg+=$'# [v2.1.9+ Updated]\n'
        upMsg+=$'# The new version of the ArtalkGo container recommends mounting\n'
        upMsg+=$'# an entire folder instead of a single file to avoid some issues.\n'
        upMsg+=$'#\n'
        upMsg+=$'# The original config file has been moved to the "/data/" folder,\n'
        upMsg+=$'# please unmount the config file volume from your container\n'
        upMsg+=$'# and edit "/data/artalk.yml" for configuration.'
        echo "$upMsg" > /conf.yml
        echo "$(date) [info] Copy config file from '/conf.yml' to '/data/artalk.yml' for upgrade"
    fi
fi

# Run ArtalkGo
artalk "$@"
