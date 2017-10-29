// Copyright 2017 The senbiot authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"github.com/johanhenselmans/senbiotpkg"
	"go.bug.st/serial.v1"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

var (
	portID          = flag.String("portID", "", "serial port to communicate")
	device          = flag.String("device", "", "Device name to use for command strings, eq ublox01b, ublox02b, quicktel")
	provider        = flag.String("provider", "", "Provider to connect to, eg, t-mobilenl, vodafone")
	message         = flag.String("message", "", "Data to send")
	cfgFile         = flag.String("config", "config.yml", "config-file for the API-settings")
	defaultName     = "ublox01b"
	defaultProvider = "t-mobilenl"
)

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
		//reader := bufio.NewReader(os.Stdin)
		//messagebyte, err = reader.ReadBytes()
		messagebyte, _ = ioutil.ReadAll(os.Stdin)
	} else if len(*message) != 0 {
		//fmt.Printf("message: %s\n", *message)
		//convert string to hex, calculate length
		var messageString string
		messageString = *message
		messagebyte = []byte(messageString)
	}

	d, err := ioutil.ReadFile(*cfgFile)
	if err != nil {
		log.Fatal("error reading config-file: %v", err)
	}

	var c senbiotpkg.Setups

	if err := yaml.Unmarshal(d, &c); err != nil {
		log.Fatal("reading config-file failed: %v", err)
	}

	//log.Print(c)
	var ChosenDevice string
	var ChosenPort string
	var ChosenProvider string

	if len(c.Device) == 0 && len(*device) == 0 {
		fmt.Println("no device name present, please set device see config.yml for devicenames eg ublox01b, ublox02b, quicktel\n")
		Usage()
		return
	} else {
		if len(*device) != 0 {
			ChosenDevice = *device
		} else {
			ChosenDevice = c.Device
		}
	}

	if len(c.PortID) == 0 && len(*portID) == 0 {
		fmt.Println("no port name present, these are the available ports:\n")
		senbiotpkg.ScanPorts()
		Usage()
		return
	} else {
		if len(*portID) != 0 {
			ChosenPort = *portID
		} else {
			ChosenPort = c.PortID
		}
	}

	if len(c.Provider) == 0 && len(*provider) == 0 {
		fmt.Println("no provider present please set provider: eg t-mobilenl, vodafone, see config.yml for provider names\n")
		Usage()
		return
	} else {
		if len(*provider) != 0 {
			ChosenProvider = *provider
		} else {
			ChosenProvider = c.Provider
		}
	}

	mode := &serial.Mode{
		BaudRate: 9600,
	}

	port, err := serial.Open(ChosenPort, mode)
	if err != nil {
		Usage()
		senbiotpkg.ScanPorts()
		log.Fatal("serial port [", ChosenPort, "] can not be opened: ", err)
	}

	var currentSetup senbiotpkg.Setup
	for _, v := range c.Stps {
		//fmt.Printf("%d = %s\n", i, v.Provider)
		if v.Provider == ChosenProvider && v.Setup == ChosenDevice {
			currentSetup = v
			break
		}
	}
	if len(currentSetup.Setup) == 0 {
		log.Fatal("could not find setup for device ", ChosenDevice, " for provider ", ChosenProvider)
	}
	// we assume the device has already been setup and a connection has been made
	SendMsgs(port, currentSetup, messagebyte)
}

//the messages section is run, with answers to be expected
func SendMsgs(port serial.Port, c senbiotpkg.Setup, messagebyte []byte) {
	dst := senbiotpkg.EncodeMessageByte(messagebyte)
	sendString := fmt.Sprintf("%s%d,%s\r\n", c.SendMessageString, len(dst), dst)
	fmt.Println(sendString)
	n, err := port.Write([]byte(sendString))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sent %v bytes\n", n)
	senbiotpkg.ReadResponse(port)

}
