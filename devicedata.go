package senbiotpkg

import (
	"go.bug.st/serial.v1"
)

func ConfigInfo(port serial.Port, c Setup) {
	for _, v := range c.ConfigInfo {
		ReadWritePort(port, v)
	}
}
