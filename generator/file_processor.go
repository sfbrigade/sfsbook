package generator

import (
	"fmt"
	"io/ioutil"
	"io"
	"log"
	"path/filepath"
	"os"
	"strings"
)

// prepareHashPath takes the provided path and converts it to a subpath
// that we can use a hash.
func prepareHashPath(path, prefix  string) string {
	// I need to strip the prefix. But... where does the prefix come from
	return strings.TrimPrefix(path, prefix)
}

// EmbedOneFile writes a single file to the output in a format appropriate for
// embedding in the code (adding it to the hash table.) The general structure:
// procesing layer should take the path and write it (e.g. minified) to the provided
// Writer. The default processor is a copy. The content doesn't need to be readable.
// TODO(rjk): Add parallel processing.
func EmbedOneFile(path, prefix string, output io.Writer) error {
	switch filepath.Ext(path) {
	case ".html", ".png":
		f, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("Can't open file path %s: %v", path, err)
		}
		array, err  := ioutil.ReadAll(f)
		if err != nil {
			return fmt.Errorf("Can't read file path %s: %v", path, err)
		}

		fmt.Fprintf(output, "%s: %#v,\n", prepareHashPath(path, prefix), string(array))
	// case ".js":
	// TODO(rjk): Invoke CSS, JS processing here. Aside: minification might
	// suggest that we want to combine resources. I suppose that means that
	// I need to refactor this code. But that should wait until later.
	default:
		log.Printf("unsupported extension %v, skipping\n", filepath.Ext(path))
	}
	return nil
}

const prefix = `// Machine generated. Do not edit. Go read ../generator/README.md
package server

var Resources map[string]string = map[string]string{
`

func WritePrefix(output io.Writer) error {
	_, err := io.WriteString(output, prefix)
	return err
}

const suffix = `}
`
func WriteSuffix(output io.Writer) error {
	_, err := io.WriteString(output, suffix)
	return err
}

