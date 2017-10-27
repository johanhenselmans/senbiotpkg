package senbiotpkg

import (
	"fmt"
	"go.bug.st/serial.v1"
	"log"
)

// We can have a range of setups for different circumstances and devices
type Setups struct {
	Device   string  `yaml:"device"`
	Provider string  `yaml:"provider"`
	PortID   string  `yaml:"portID"`
	Stps     []Setup `yaml:"setups"`
}

// Setup struct has the complete sequence of commands
type Setup struct {
	Setup             string            `yaml:"setup"`
	Date              string            `yaml:"date"`
	Provider          string            `yaml:"provider"`
	Reboot            []RequestResponse `yaml:"reboot"`
	Init              []RequestResponse `yaml:"init"`
	SetupNetwork      []RequestResponse `yaml:"setupnetwork"`
	WaitForNetwork    []RequestResponse `yaml:"waitfornetwork"`
	ConfigInfo        []RequestResponse `yaml:"configinfo"`
	NetworkInfo       []RequestResponse `yaml:"networkinfo"`
	GetMsgResponse    []RequestResponse `yaml:"getmsgresponse"`
	SendMessageString string            `yaml:"sendmesssagestring"`
}

type RequestResponse struct {
	Request          string `yaml:"request"`
	Response         string `yaml:"response,omitempty""`
	NegativeResponse string `yaml:"negativeresponse,omitempty""`
	WaitForResponse  string `yaml:"waitforresponse,omitempty""`
}

func ScanPorts() {
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		log.Fatal("No serial ports found!")
	}
	for _, port := range ports {
		fmt.Printf("Found port: %v\n", port)
	}

}
