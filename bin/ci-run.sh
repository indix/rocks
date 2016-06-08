#!/bin/bash

set -e -x

BASE=`dirname $0`/../

ROCKSB_VERSION="4.1"

if [ "${SNAP_CI}x" == "truex" ]; then
  cd ${BASE}
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
fi

CGO_CFLAGS="-I${BASE}/build/rocksdb-${ROCKSB_VERSION}/include" CGO_LDFLAGS="-L${BASE}/build/rocksdb-${ROCKSB_VERSION} -lrocksdb -lstdc++ -lm -lz -lbz2 -lsnappy" go build -o ${APPNAME} .
