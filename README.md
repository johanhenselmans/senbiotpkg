# Senbiot

[![Build Status](https://travis-ci.org/johanhenselmans/senbiot.svg?branch=master)](https://travis-ci.org/johanhenselmans/senbiot)
[![Go Report Card](https://goreportcard.com/badge/github.com/johanhenselmans/senbiot)](https://goreportcard.com/report/github.com/johanhenselmans/senbiot)
[![GoDoc](https://godoc.org/github.com/johanhenselmans/senbiot?status.svg)](https://godoc.org/github.com/johanhenselmans/senbiot)
[![codecov](https://codecov.io/gh/johanhenselmans/senbiot/branch/master/graph/badge.svg)](https://codecov.io/gh/johanhenselmans/senbiot)


NB-IOT is a LPWAN technology for bi-directional data traffic between devices and centralized cloud platforms. To configure such a device, a whole slew of AT commands is needed. Senbiot tries to simplify that to a standard set of settings and a simple senbiot "message" command


## Installation

To get started please get the Golang toolchains from the [Golang website](https://golang.org/). When you have a working go toolchain you can do:

```
go get github.com/johanhenselmans/senbiotpkg
```

And you are ready to go!

## Included tools

The tools made with senbiotpkg are located in the `cmd` folder of the root of the project. You just do a go build in the specific folder and run the resulting commands. 

### Send a message via an NB-IOT device (senbiot)

Commandline tool to send a message to an NB-``IOT network. Currently the device supported is the SODAQ NBIOT device, at https://shop.sodaq.com/en/nb-iot-shield-deluxe-dual-band-8-20.html. The network that are supported are the Vodafone and T-Mobile networks in the Netherlands. The device requires to have a 'through' connection to the serial port of the ublox device, which can be accomplished by using the Arduino sketch from http://support.sodaq.com/sodaq-one/at/. I have included the Arduino sketch in the folder SerialThrough, You should upload this to your Arduino board. That will make the
 connection transparant from you linux/windows/macos machine, and you can shoot messages to the board.

Install via ```go install github.com/johanhenselmans/cmd/senbiot```

### Check the configuration of your NB-IOT shield (checkconfig)

CommandLine tool to check the configuration of your NB-IOT device. Configuration can be given via the commandline or via a config.yml file. See the example config.yml file included

Install via ```go install github.com/johanhenselmans/cmd/checkconfig```a

### Send a message via your NB-IOT shield (sendmsg)

CommandLine tool to send a message via your preconfigured NB-IOT device. Configuration can be given via the commandline or via a config.yml file. See the example config.yml file included

Install via ```go install github.com/johanhenselmans/cmd/sendmsg```


### Get serialports (getserialports)

Commandline tool to scan serialports on the machine so as to determine which device to use.

Install via ```go install github.com/johanhenselmans/cmd/getserialports```


### Encode a message to be used in a NB-IOT message (encodemessage)

CommandLine tool to encode a message in the way it will be sent via your NB-IOT device. This encoding does not count the lenght of the message, as specified in the NB-IOT message format. 

Install via ```go install github.com/johanhenselmans/cmd/encodemessage``` 


### Decode a message used in a NB-IOT message (decodemessage)

CommandLine tool to encode a decode the hex content of the message that could have been sent from your NB-IOT device to the gateway

Install via ```go install github.com/johanhenselmans/cmd/decodemessage```

### Decode a message to be used in a NB-IOT message (decodebase64message)

CommandLine tool to decode the base64 message that is sent via the OceanConnect gateway of T-Mobile.

Install via ```go install github.com/johanhenselmans/cmd/decodebase64message```


## Plans

I have plans to get the board to work via Firmata, that should make it possible to retrieve the GPS coordinates from the board. 

## Contributing

Please read the [Contribution Guidelines](CONTRIBUTING.md). Furthermore: Fork -> Patch -> Push -> Pull Request

## License

This project is licensed under the MIT License. See the [LICENSE](https://github.com/johanhenselmans/senbiot/blob/master/LICENSE) file for the full license text.
