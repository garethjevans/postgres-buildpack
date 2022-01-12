#!/bin/bash

CONFIG_FILE=$1
PACKAGE=$2
OS=$3

cat <<- EOF > ${CONFIG_FILE}
[buildpack]
uri = "${PACKAGE}"

[platform]
os = "${OS}"
EOF
