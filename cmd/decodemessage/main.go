// Copyright 2017 The senbiot authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"github.com/johanhenselmans/senbiotpkg"
	"io/ioutil"
	"os"
)

var (
	message = flag.String("message", "", "Data to send")
)

type command []string

func (c *command) String() string {
	return fmt.Sprint(*c)
}

func main() {
	flag.Parse()

	var Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  cat yourfile.txt | %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	pipemessage, _ := os.Stdin.Stat()
	var messagebyte []byte

	if (pipemessage.Mode()&os.ModeCharDevice) == os.ModeCharDevice && len(*message) == 0 && flag.NFlag() == 0 {
		Usage()
		return
	} else if pipemessage.Size() > 0 {
		messagebyte, _ = ioutil.ReadAll(os.Stdin)
	} else if len(*message) != 0 {
		//fmt.Printf("%s", *message)
		//convert string to []byte
		var messageString string
		messageString = *message
		messagebyte = []byte(messageString)
	}
	fmt.Printf("%s", senbiotpkg.DecodeMessageByte(messagebyte))

}
