# sfsbook

STOP VIOLENCE!

## Structure

We use [Gin http framework](https://github.com/gin-gonic/gin) for routing and
rendering of paths.

main.go - contains the Gin initializer and sets up various routes and runs the
server.

routes.go - contains various handlers; will be moved to own package once this
file becomes too big.

templates/ - contains various Go templates.

static/ - contains js and css files; this folder is server to the public, so
BEWARE!

refguides/ - contains the resources in various formats; nested by date pdf was
added.

vendor/ - contains vendored dependencies;

## Development

Make sure you've Go and GOPATH setup. Follow instructions here:
https://golang.org/doc/install.

Next, install this repo with:

    go get github.com/sfbrigade/sfsbook

It should be installed at $GOPATH/src/github.com/sfbrigade.

Next, cd into the repo and install [gin](https://github.com/codegangsta/gin)
which is a live reload
utility.

    cd $GOPATH/src/github.com/sfbrigade
    go get github.com/codegangsta/gin

Run server with:

    gin run main.go

You can access the server at: http://localhost:3000.

Optionally, if you add a new dep, be sure to add it to vendor/ with:

    go get -u github.com/kardianos/govendor
    govendor add +external

Be sure to commit the new contents of vendor/ folder along with your changes.

## Credits

This repo's frontend and javascript code is original from
https://github.com/sfbrigade/sfsbook. All credit goes to @rjkroege
and @cehsu and many other contributors seen
[here](https://github.com/sfbrigade/sfsbook/graphs/contributors).
