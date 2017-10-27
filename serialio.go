package senbiotpkg

import (
	"fmt"
	"go.bug.st/serial.v1"
	"log"
	"time"
)

func ReadResponse(port serial.Port) (response string) {
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

func ReadWritePort(port serial.Port, v RequestResponse) (response string) {
	var n int
	var err error
	fmt.Printf("%s\n", v.Request)
	n, err = port.Write([]byte(fmt.Sprintf("%s\r\n", v.Request)))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sent %v bytes\n", n)
	resultString := ReadResponse(port)
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
