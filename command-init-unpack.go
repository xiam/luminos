// Copyright (c) 2012-2014 José Carlos Nieto, https://menteslibres.net/xiam
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package main

import (
	"archive/tar"
	"bytes"
	"compress/bzip2"
	"fmt"
	"io"
	"os"
)

// Attempts to extract a compressed project into the given destination.
func unpackExampleProject(root string) (err error) {
	var stat os.FileInfo

	// Validating destination.
	if stat, err = os.Stat(root); err != nil {
		return err
	}

	if stat.IsDir() == false {
		return fmt.Errorf("Expecting a directory.")
	}

	// Creating a tarbz2 reader.
	tbz := tar.NewReader(bzip2.NewReader(bytes.NewBuffer(compressedProject)))

	// Extracting tarred files.
	for {

		hdr, err := tbz.Next()

		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err.Error())
		}

		// See http://en.wikipedia.org/wiki/Tar_(computing)
		filePath := root + pathSeparator + hdr.Name

		switch hdr.Typeflag {
		case '0':
			// Normal file
			fp, err := os.Create(filePath)

			if err != nil {
				return err
			}

			io.Copy(fp, tbz)

			err = os.Chmod(filePath, os.FileMode(hdr.Mode))

			if err != nil {
				return err
			}

			fp.Close()
		case '5':
			// Directory
			os.MkdirAll(filePath, os.FileMode(hdr.Mode))
		default:
			// fmt.Printf("--> %s, %d, %c\n", hdr.Name, hdr.Mode, hdr.Typeflag)
			panic(fmt.Sprintf("Unhandled tar type: %c in file: %s", hdr.Typeflag, hdr.Name))
		}
	}

	return nil

}
