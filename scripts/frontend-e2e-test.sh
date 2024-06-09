#!/usr/bin/env bash

set -e

# Check if `data/ui_test.db` exits
if [ ! -f data/ui_test.tpl.db ]; then
    echo "data/ui_test.tpl.db not found!"
    exit 1
fi

# Copy `data/ui_test.tpl.db` to `data/ui_test.db`
cp data/ui_test.tpl.db data/ui_test.db

# Check if yq is installed
if ! [ -x "$(command -v yq)" ]; then
    echo -e 'Error: yq is not installed. \nSee: https://github.com/mikefarah/yq#install' >&2
    exit 1
fi

# Modify `artalk.yml` db.file to `data/ui_test.db` use yq
yq -i '.db.file = "data/ui_test.db"' artalk.yml

# Run frontend e2e test
pnpm -F artalk run test:e2e

# if args contains `--show-report` then open report
if [[ $* == *--show-report* ]]; then
    pnpm -F artalk run test:e2e-report
fi
