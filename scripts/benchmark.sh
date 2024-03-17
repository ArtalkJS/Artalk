#!/bin/bash

bombardier -c 125 -d 10s -m GET "http://127.0.0.1:23366/api/v2/comments?limit=10&offset=0&flat_mode=false&page_key=https%3A%2F%2Fartalk.js.org%2Fguide%2Fintro.html&site_name=ArtalkDocs"
