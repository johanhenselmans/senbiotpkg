// Copyright 2017 The senbiot authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/hex"
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
		flag.PrintDefaults()
	}
	if flag.NFlag() == 0 {
		Usage()
		return
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

	if len(*device) == 0 {
		fmt.Println("no device name present, please set device see config.yml for devicenames eg ublox01b, ublox02b, quicktel\n")
		Usage()
		return
	}

	if len(*portID) == 0 {
		fmt.Println("no port name present, these are the available ports:\n")
		senbiotpkg.ScanPorts()
		Usage()
		return
	}

	if len(*provider) == 0 {
		fmt.Println("no provider present please set provider: eg t-mobilenl, vodafone, see config.yml for provider names\n")
		Usage()
		return
	}

	mode := &serial.Mode{
		BaudRate: 9600,
	}

	port, err := serial.Open(*portID, mode)
	if err != nil {
		log.Fatal("serial port can not be opened: ", err, *portID)
	}
	var currentSetup senbiotpkg.Setup
	for _, v := range c.Stps {
		//fmt.Printf("%d = %s\n", i, v.Provider)
		if v.Provider == fmt.Sprintf("%s", *provider) && v.Setup == fmt.Sprintf("%s", *device) {
			currentSetup = v
			break
		}
	}
	if len(currentSetup.Setup) == 0 {
		log.Fatal("could not find setup for device ", *device, " for provider ", *provider)
	}
	if len(commands) > 0 {
		for _, aCommand := range commands {
			fmt.Printf("command: %s \n", aCommand)
			switch aCommand {
			case "ConfigInfo":
				configDevice(port, currentSetup)
			case "NetworkInfo":
				networkDevice(port, currentSetup)
			case "Init":
				setupInit(port, currentSetup)
			case "Reboot":
				rebootDevice(port, currentSetup)
			case "SetupNetwork":
				setupNetwork(port, currentSetup)
			case "SendMessage":
				sendMsgs(port, currentSetup, message)
			case "WaitForNetwork":
				waitForNetwork(port, currentSetup)
			case "ScanPorts":
				senbiotpkg.ScanPorts()
			}
		}
	} else {
		// we assume the device has already been setup
		setupNetwork(port, currentSetup)
		waitForNetwork(port, currentSetup)
		sendMsgs(port, currentSetup, message)
	}
}

// rebootDevice reboots the device and waits 7 seconds for it to come up
func rebootDevice(port serial.Port, c senbiotpkg.Setup) {
	for _, v := range c.Reboot {
		readwriteport(port, v)
	}
	const delay = 7000 * time.Millisecond
	time.Sleep(delay)
}

func configDevice(port serial.Port, c senbiotpkg.Setup) {
	for _, v := range c.ConfigInfo {
		readwriteport(port, v)
	}
}

func networkDevice(port serial.Port, c senbiotpkg.Setup) {
	for _, v := range c.NetworkInfo {
		readwriteport(port, v)
	}
}

//the init section of the yaml page of the device with answers is run, stored in nv memory, has to be run only once
func setupInit(port serial.Port, c senbiotpkg.Setup) {
	for _, v := range c.Init {
		readwriteport(port, v)
	}
}

func setupNetwork(port serial.Port, c senbiotpkg.Setup) {
	for _, v := range c.SetupNetwork {
		readwriteport(port, v)
	}
}

func waitForNetwork(port serial.Port, c senbiotpkg.Setup) string {
	var result string
	for _, v := range c.WaitForNetwork {
		var i int
		for i < 10 {
			result = readwriteport(port, v)
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

func readResponse(port serial.Port) (response string) {
	var n int
	var err error
	buff := make([]byte, 10)
	var stringbuff string
	var noofreturns int = 0

	for {
		n, err = port.Read(buff)
		if err != nil {
			log.Fatal(err)
			break
		}
		if n == 0 {
			fmt.Println("\nEOF")
			break
		}
		//fmt.Printf("%v\n", string(buff[:n]))
		switch string(buff[:n]) {
		case "\r":
			//fmt.Println("found carriage return")
		case "\n":
			//fmt.Println("found newline")
			noofreturns++
		default:
			stringbuff = fmt.Sprintf("%s%s", stringbuff, string(buff[:n]))
		}
		if noofreturns > 1 {
			fmt.Printf("result: %s\n", stringbuff)
			break
		}
	}
	return stringbuff
}

func readwriteport(port serial.Port, v senbiotpkg.RequestResponse) (response string) {
	var n int
	var err error
	fmt.Printf("%s\n", v.Request)
	n, err = port.Write([]byte(fmt.Sprintf("%s\r\n", v.Request)))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sent %v bytes\n", n)
	resultString := readResponse(port)
	if len(v.Response) != 0 && resultString != v.Response {
		log.Fatal("response was:", resultString, "expected: ", v.Response)
	} else {
		response = resultString
	}
	// Wait two seconds to have this machine stabilize a bit
	const delay = 1000 * time.Millisecond
	time.Sleep(delay)

	return response
}

//the messages section is run, with answers to be expected
func sendMsgs(port serial.Port, c senbiotpkg.Setup, message *string) {
	fmt.Printf("%s", message)
	//convert string to hex, calculate length
	var messageString string
	messageString = *message
	src := []byte(messageString)
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)
	//fmt.Printf("%s %d %s, %s\n", dst, len(dst), src, messageString)
	sendString := fmt.Sprintf("%s%d,%s\r\n", c.SendMessageString, len(dst), dst)
	fmt.Println(sendString)
	n, err := port.Write([]byte(sendString))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sent %v bytes\n", n)
	readResponse(port)

}
