package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/sfbrigade/sfsbook/generator"
)

// Reads all .html files in the current folder
// and encodes them as strings literals in textfiles.go
// TODO(rjk): Extend this with the rich selection of thints that I need to support.

// nested directories... 
// i needs a hash. slash is not a valid symbol. the path w.r.t. site needs is a key.

// some args.
var (
	outputfile = flag.String("output", "", "Write generated source here.")
	prefix = flag.String("prefix", "", "The prefix of our website contents to intern.")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "	%s <flags listed below>\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	var output io.Writer

	if *outputfile != "" {
		f, err := os.Create(*outputfile)
		if err != nil {
			log.Fatalln("Can't create specified output file:", outputfile, "for writing because:", err)
		}
		defer f.Close()
		output = f
	} else {
		output = os.Stdout
	}


	// TODO(rjk): Don't trash the previous file unless we succeed?
	// nah. we have version control for a reason.

	// write prefix
	generator.WritePrefix(output)
	defer generator.WriteSuffix(output)

	for _, pth := range flag.Args() {
		_, err := os.Stat(pth)
		if  err != nil {
			log.Println("Skipping un-stat-able argument:", pth, "because:", err)
			continue
		}

		// file walk...
		if err := filepath.Walk(pth, func(path string, info os.FileInfo, err error) error {
			log.Println("filewalk visiting", path)

			if err != nil && info.IsDir() {
				return filepath.SkipDir
			} else if err != nil {
				return nil
			} else if info.IsDir() {
				return nil
			}

			// This is not very sophisticated. I can imagine that an advanced JS
			// minification scheme will need to be smarter. In particular, I might
			// want to group together all the files of similar 
			return generator.EmbedOneFile(path, *prefix, output)
		}); err != nil {
			log.Println("file walking had an error on path:", pth, "because", err)
		}
	}
}
