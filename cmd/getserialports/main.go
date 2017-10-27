package main

import (
	"log"
	"fmt"
	"go.bug.st/serial.v1"
)

func main() {
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







