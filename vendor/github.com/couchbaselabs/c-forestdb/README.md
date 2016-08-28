# ForestDB

This is a go-gettable version of the ForestDB C++ library for use in Go code
that needs to link against the C++ ForestDB library but wants to integrate with
go get and go build.

To use in your project you need to import the package and set appropriate
cgo flag directives:

import _ "github.com/couchbaselabs/c-forestdb"

    // #cgo CXXFLAGS: -std=c++11
    // #cgo CPPFLAGS: -I ../../couchbaselabs/c-forestdb/internal/include
    // #cgo darwin LDFLAGS: -Wl,-undefined -Wl,dynamic_lookup
    // #cgo !darwin LDFLAGS: -Wl,-unresolved-symbols=ignore-all
    // #include <libforestdb/forestdb.h>
    import "C"

To update the upstream version of ForestDB you'll want to update ./import.sh
to point to the new version (just change the URL), and then run it.

### Credits

This package was inspired by the [cockroachdb](https://github.com/cockroachdb) teams
[c-rocksdb](https://github.com/cockroachdb/c-rocksdb) project wrapping the RocksDB library.
