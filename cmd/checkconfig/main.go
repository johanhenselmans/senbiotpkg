// Copyright 2017 The senbiot authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/johanhenselmans/senbiotpkg"
	"go.bug.st/serial.v1"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var (
	portID          = flag.String("portID", "", "serial port to communicate")
	device          = flag.String("device", "", "Device name to use for command strings, eq ublox01b, ublox02b, quicktel")
	provider        = flag.String("provider", "", "Provider to connect to, eg, t-mobilenl, vodafone")
	cfgFile         = flag.String("config", "config.yml", "config-file for the API-settings")
	defaultName     = "ublox01b"
	defaultProvider = "t-mobilenl"
)

type command []string

func (c *command) String() string {
	return fmt.Sprint(*c)
}

// Set is the method to set the flag value, part of the flag.Value interface.
// Set's argument is a string to be parsed to set the flag.
// It's a comma-separated list, so we split it.
func (c *command) Set(value string) error {
	// If we wanted to allow the flag to be set multiple times,
	// accumulating values, we would delete this if statement.
	// That would permit usages such as
	//-deltaT 10s -deltaT 15s
	// and other combinations.
	if len(*c) > 0 {
		return errors.New("command flag already set")
	}
	for _, dt := range strings.Split(value, ",") {
		commandsstring := dt
		*c = append(*c, commandsstring)
	}
	return nil
}

var commands command

func main() {
	flag.Var(&commands, "command", "comma-separated list of commands (ConfigInfo, NetworkInfo, ScanPorts) to use ")
	flag.Parse()

	var Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  cat yourfile.txt | %s:\n", os.Args[0])

		flag.PrintDefaults()
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
	if len(commands) > 0 {
		for _, aCommand := range commands {
			fmt.Printf("command: %s \n", aCommand)
			switch aCommand {
			case "ConfigInfo":
				senbiotpkg.ConfigInfo(port, currentSetup)
			case "NetworkInfo":
				senbiotpkg.NetworkInfo(port, currentSetup)
			case "ScanPorts":
				senbiotpkg.ScanPorts()
			}
		}
	} else {
		// we assume the device has already been setup
		senbiotpkg.ScanPorts()
		senbiotpkg.ConfigInfo(port, currentSetup)
		senbiotpkg.NetworkInfo(port, currentSetup)
	}
}
