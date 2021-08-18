#!/bin/bash

# ./push.sh will only run when a new tag is 
# pushed to github and github actions runs the 
# push-release.yml action

set -eu

S3_TARBALL="traverse.tar.gz"
S3_BUCKET="lucassha-traverse-releases"
OS="darwin"
ARCH="amd64"
BINARY="traverse"

LATEST_RELEASE=$(git describe --tags --abbrev=0)

GOOS=${OS} GOARCH=${ARCH} go build -o ${BINARY}

tar -czvf ${S3_TARBALL}-${LATEST_RELEASE} ./${BINARY}

aws s3 cp ${S3_TARBALL}-${LATEST_RELEASE} s3://${S3_BUCKET}/releases/${LATEST_RELEASE}