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
	"flag"
	"fmt"
	"os"

	"menteslibres.net/gosexy/cli"
)

// initCommand is the structure that provides instructions for the "luminos
// init" subcommand.
type initCommand struct {
}

// Execute creates a new luminos site scaffold in the given PATH.
func (c *initCommand) Execute() (err error) {
	var stat os.FileInfo

	// Default PATH if the current working directory.
	dest := "."

	// If a PATH was given, use it instead of the default PATH.
	if flag.NArg() > 1 {
		dest = flag.Arg(1)
	}

	// Verifying PATH.
	if stat, _ = os.Stat(dest); stat == nil {
		// Directory does not exists, try to create it.
		// TODO: Use system's default mask.
		if err = os.MkdirAll(dest, os.ModeDir|0755); err != nil {
			return err
		}
	} else {
		// Path exists, is it a directory?
		if stat.IsDir() == false {
			// Nope, the we can't use it.
			return fmt.Errorf("Cannot create directory, file %s already exists!", dest)
		}
	}

	// If the PATH was already initialized, then it must contain a LOCKFILE.
	lockFile := dest + pathSeparator + ".luminos"

	if stat, _ = os.Stat(lockFile); stat != nil {
		// If the LOCKFILE exists we cannot continue.
		if dest == "." {
			// Use "the current directory" to avoid confusing non-technical users
			// with a dot instead of a full PATH.
			return fmt.Errorf("A Luminos project already exists at the current directory.")
		}
		return fmt.Errorf("A luminos project already exists in %s.", dest)
	}

	// We may extract the example project now.
	if err = unpackExampleProject(dest); err != nil {
		return err
	}

	// And then create the LOCKFILE.
	var lfp *os.File
	if lfp, err = os.Create(lockFile); err != nil {
		return err
	}
	lfp.Close()

	// All done! Let's tell the user we've finished.
	if dest == "." {
		fmt.Printf("New project created at the current directory.\n")
	} else {
		fmt.Printf("New project created at %s.\n", dest)
	}

	return nil
}

func init() {
	// Describing the "init" subcommand.
	cli.Register("init", cli.Entry{
		Description: "Creates a new Luminos site scaffold in the given PATH.",
		Usage:       "init [PATH]",
		Command:     &initCommand{},
	})
}
