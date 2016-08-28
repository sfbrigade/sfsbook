#!/usr/bin/env sh

set -eu

rm -rf *.cc internal/
git clone https://github.com/couchbase/forestdb internal

# symlink so cgo compiles them
for source_file in $(make sources); do
    cp $source_file .
done

# rename filemgr_ops_linux.cc because Go will ignore it on
# platforms other than linux
mv filemgr_ops_linux.cc filemgr_ops_posix.cc
