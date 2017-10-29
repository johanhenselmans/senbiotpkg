package senbiotpkg

import (
	"encoding/hex"
	"fmt"
	"log"
)

func EncodeMessageByte(messagebyte []byte) string {
	//fmt.Println(messagebyte)
	dst := make([]byte, hex.EncodedLen(len(messagebyte)))
	hex.Encode(dst, messagebyte)
	sendString := fmt.Sprintf("%s", dst)
	//fmt.Println(sendString)
	return sendString
}

func EncodeMessageString(messageString string) string {
	//fmt.Println(messageString)
	messagebyte := []byte(messageString)
	dst := make([]byte, hex.EncodedLen(len(messagebyte)))
	hex.Encode(dst, messagebyte)
	sendString := fmt.Sprintf("%s", dst)
	//fmt.Println(sendString)
	return sendString
}

func DecodeMessageByte(messagebyte []byte) string {
	decoded := make([]byte, hex.DecodedLen(len(messagebyte)))
	n, err := hex.Decode(decoded, messagebyte)
	if err != nil {
		log.Fatal(err)
	}
	receiveString := fmt.Sprintf("%s", decoded[:n])
	return receiveString
}

func DecodeMessageString(message string) string {
	decoded, _ := hex.DecodeString(message)
	receiveString := fmt.Sprintf("%s", decoded)
	return receiveString
}
