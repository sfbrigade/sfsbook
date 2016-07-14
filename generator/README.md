# Scheme
The general idea is as follows:

*  There are some number of static processing passes.
* They can be used from within the web server while running if
operating in debug mode.
* They can be used by the generate tool to construct resources for
deployment.

The tool binary is in `./tool`. Libraries for generation (e.g. JS minification) go
here.

