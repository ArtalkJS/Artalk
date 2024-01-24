#!/bin/bash

bombardier -c 125 -d 10s -b "site_name=ArtalkDocs&page_key=https%3A%2F%2Fartalk.js.org%2Fguide%2Fintro.html&limit=20&offset=0" -m GET -H "Origin: http://127.0.0.1:5173" -H "Content-Type: application/x-www-form-urlencoded" http://127.0.0.1:23366/api/v2/comments
