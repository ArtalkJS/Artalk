#!/bin/bash
set -e

if [ "$1" != 'gen' ] && ( [ ! -e /data/artalk.yml ] && [ ! -e /data/artalk-go.yml ] ); then
    if [ -e /conf.yml ]; then
        # Move original config to `/data/` for upgrade (<= v2.1.8)
        cp /conf.yml /data/artalk.yml
        upMsg=""
        upMsg+=$'# [v2.1.9+ Updated]\n'
        upMsg+=$'# The new version of the Artalk container recommends mounting\n'
        upMsg+=$'# an entire folder instead of a single file to avoid some issues.\n'
        upMsg+=$'#\n'
        upMsg+=$'# The original config file has been moved to the "/data/" folder,\n'
        upMsg+=$'# please unmount the config file volume from your container\n'
        upMsg+=$'# and edit "/data/artalk.yml" for configuration.'
        echo "$upMsg" > /conf.yml
        echo "$(date) [info][docker] Copy config file from '/conf.yml' to '/data/artalk.yml' for upgrade"
    else
        # Generate new config
        artalk gen conf /data/artalk.yml
        echo "$(date) [info][docker] Generate new config file to '/data/artalk.yml'"
    fi
fi

# Run Artalk
/usr/bin/artalk "$@"
