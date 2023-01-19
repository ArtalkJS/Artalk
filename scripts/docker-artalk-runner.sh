#!/bin/bash
if [ -e /data/artalk.yml ]; then
    /artalk -w / -c /data/artalk.yml "$@"
else
    /artalk -w / -c /data/artalk-go.yml "$@"
fi
