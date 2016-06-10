#!/bin/bash

set -e -x

BASE="$( cd "$( dirname "$0" )"/.. && pwd )"

ROCKSB_VERSION="4.1"

if [ "${SNAP_CI}x" == "truex" ]; then
  pushd ${SNAP_CACHE_DIR}
  # Download and extract rocksdb source file
  mkdir -p build/
  pushd build
  if [ ! -d rocksdb-${ROCKSB_VERSION} ]; then
    wget --continue https://github.com/facebook/rocksdb/archive/v${ROCKSB_VERSION}.tar.gz
    tar xzvf v${ROCKSB_VERSION}.tar.gz
    pushd rocksdb-${ROCKSB_VERSION}
    make static_lib
    popd
  fi
  popd
  popd
fi

sudo apt-get install -y libsnappy-dev
export CGO_CFLAGS="-I${SNAP_CACHE_DIR}/build/rocksdb-${ROCKSB_VERSION}/include"
export CGO_LDFLAGS="-L${SNAP_CACHE_DIR}/build/rocksdb-${ROCKSB_VERSION} -lrocksdb -lstdc++ -lm -lz -lbz2 -lsnappy"

make build
make test
