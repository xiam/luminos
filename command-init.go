/*
  Copyright (c) 2012-2013 José Carlos Nieto, http://xiam.menteslibres.org/

  Permission is hereby granted, free of charge, to any person obtaining
  a copy of this software and associated documentation files (the
  "Software"), to deal in the Software without restriction, including
  without limitation the rights to use, copy, modify, merge, publish,
  distribute, sublicense, and/or sell copies of the Software, and to
  permit persons to whom the Software is furnished to do so, subject to
  the following conditions:

  The above copyright notice and this permission notice shall be
  included in all copies or substantial portions of the Software.

  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
  EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
  MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
  NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
  LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
  OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
  WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package main

import (
	"flag"
	"fmt"
	"github.com/gosexy/cli"
	"os"
)

func init() {
	cli.Register("init", cli.Entry{
		Description: "Initializes a working directory with a Luminos base project.",
		Usage:       "init [directory]",
		Command:     &initCommand{},
	})
}

type initCommand struct {
}

func (self *initCommand) Execute() error {

	var err error

	// Default destinarion is current working directory
	dest := "."

	if flag.NArg() > 1 {
		dest = flag.Arg(1)
	}

	stat, _ := os.Stat(dest)

	if stat == nil {
		// Directory does not exists, attemping to create it.
		err = os.MkdirAll(dest, os.ModeDir|0755)
		if err != nil {
			return err
		}
	} else {
		// Path exists, is it a directory?
		if stat.IsDir() == false {
			return fmt.Errorf("Cannot create directory, file %s already exists.", dest)
		}
	}

	lockFile := dest + PS + ".luminos"

	stat, _ = os.Stat(lockFile)

	if stat == nil {

		err = unpackExampleProject(dest)

		if err != nil {
			return err
		}

		lfp, _ := os.Create(lockFile)
		lfp.Close()

		if dest == "." {
			fmt.Printf("Created an empty Luminos project in the current directory.\n")
		} else {
			fmt.Printf("Created an empty Luminos project in %s/.\n", dest)
		}

	} else {
		if dest == "." {
			return fmt.Errorf("A Luminos project already exists in the current directory.\n")
		} else {
			return fmt.Errorf("A Luminos project already exists in %s/.\n", dest)
		}
	}

	return nil
}
