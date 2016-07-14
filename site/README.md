All the contents here get swallowed up into the executable for
deployment. But when the web server (`sfsbook`) is run from the parent
directory, it serves files from here inplace of the files compiled
into the executable.

To update the files compiled into the executable, run go generate at the top
level and rebuild. i.e.:

	cd ..; go generate ./... && go build 

The site layout can be reasonably arbitrary and under the control of
front-end developers. 
