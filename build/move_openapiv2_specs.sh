#!/bin/bash

PATH_FROM="$1"
PATH_TO="$2"

if [ "$PATH_FROM" == "" ]; then
    echo "Please specify path to folder with openapiv2 specifications"
    exit 1
fi

if [ "$PATH_TO" == "" ]; then
    echo "Please specify path to folder where to move openapiv2 specifications"
    exit 1
fi

if [ -d "$PATH_TO" ]; then
    rm -rf "{$PATH_TO}/*"
else
    mkdir "$PATH_TO"
fi

find "$PATH_FROM" -name '*.swagger.json' -exec mv {} "$PATH_TO" \;