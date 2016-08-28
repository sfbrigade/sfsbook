#!/usr/bin/env sh

set -eu

rm -rf *.cc internal/*
curl -sL https://github.com/google/snappy/archive/1.1.3.tar.gz | tar zxf - -C internal --strip-components=1
(cd internal && ./autogen.sh && ./configure)

# copy (not symlink thanks to windows) so cgo compiles them
for source_file in $(make sources); do
  cp $source_file .
done

# restore the repo to what it would look like when first cloned.
# comment this line out while updating upstream.
git clean -dxf
