#!/bin/bash

if [ -z ${1} ]; then
	echo "No path giving. Using current path: $(pwd)"
	path=$(pwd)
else
	path=${1}
fi

docker build -t tradfri-go-blind-server ${path}
