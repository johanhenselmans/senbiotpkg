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
	"time"
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
	flag.Var(&commands, "command", "comma-separated list of commands (Reboot, Init, SetupNetwork, WaitForNetwork, ConfigInfo, NetworkInfo, SendMessage, ScanPorts) to use ")
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
		fmt.Printf("%s", message)
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
		log.Fatal("serial port can not be opened: ", err, *portID)
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
				NetworkInfo(port, currentSetup)
			case "Init":
				SetupInit(port, currentSetup)
			case "Reboot":
				RebootDevice(port, currentSetup)
			case "SetupNetwork":
				SetupNetwork(port, currentSetup)
			case "SendMessage":
				SendMsgs(port, currentSetup, messagebyte)
			case "WaitForNetwork":
				WaitForNetwork(port, currentSetup)
			case "ScanPorts":
				senbiotpkg.ScanPorts()
			}
		}
	} else {
		// we assume the device has already been setup
		SetupNetwork(port, currentSetup)
		WaitForNetwork(port, currentSetup)
		SendMsgs(port, currentSetup, messagebyte)
	}
}

// rebootDevice reboots the device and waits 7 seconds for it to come up
func RebootDevice(port serial.Port, c senbiotpkg.Setup) {
	for _, v := range c.Reboot {
		senbiotpkg.ReadWritePort(port, v)
	}
	const delay = 7000 * time.Millisecond
	time.Sleep(delay)
}

func NetworkInfo(port serial.Port, c senbiotpkg.Setup) {
	for _, v := range c.NetworkInfo {
		senbiotpkg.ReadWritePort(port, v)
	}
}

//the init section of the yaml page of the device with answers is run, stored in nv memory, has to be run only once
func SetupInit(port serial.Port, c senbiotpkg.Setup) {
	for _, v := range c.Init {
		senbiotpkg.ReadWritePort(port, v)
	}
}

func SetupNetwork(port serial.Port, c senbiotpkg.Setup) {
	for _, v := range c.SetupNetwork {
		senbiotpkg.ReadWritePort(port, v)
	}
}

func WaitForNetwork(port serial.Port, c senbiotpkg.Setup) string {
	var result string
	for _, v := range c.WaitForNetwork {
		var i int
		for i < 10 {
			result = senbiotpkg.ReadWritePort(port, v)
			fmt.Printf("result = %s neg response: %s\n", result, v.NegativeResponse)
			if !strings.Contains(result, v.NegativeResponse) {
				i = 0
				break
			} else {
				const delay = 1000 * time.Millisecond
				time.Sleep(delay)
				i++
			}
		}
		if i != 0 {
			log.Fatal("could not get connection: ", result)
		}
	}
	return result
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
