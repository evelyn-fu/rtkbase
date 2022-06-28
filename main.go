package main

import (
	"fmt"
	"log"
	//"io"
	//"bufio"

	"github.com/go-gnss/rtcm/rtcm3"
    "github.com/jacobsa/go-serial/serial"

	"rtcmReading/configure/base"
)

func main() {
	options := serial.OpenOptions{
		PortName: "/dev/serial/by-path/platform-fd500000.pcie-pci-0000:01:00.0-usb-0:1.3:1.0", // change to base port
		BaudRate: 115200,
		DataBits: 8,
		StopBits: 1,
		MinimumReadSize: 1,
	}

	// Open the port.
	readPort, err := serial.Open(options)
	if err != nil {
		log.Fatalf("serial.Open: %v", err)
	}
	defer readPort.Close()


	// options = serial.OpenOptions{
	// 	PortName: "/dev/serial/by-path/platform-fd500000.pcie-pci-0000:01:00.0-usb-0:1.4:1.0", // change to rover port
	// 	BaudRate: 115200,
	// 	DataBits: 8,
	// 	StopBits: 1,
	// 	MinimumReadSize: 1,
	// }

	// // Open the port.
	// writePort, err := serial.Open(options)
	// if err != nil {
	// 	log.Fatalf("serial.Open: %v", err)
	// }
	// defer writePort.Close()

	// w := bufio.NewWriter(writePort)
	// r := io.TeeReader(readPort, w)

	base.disableAll()
	scanner := rtcm3.NewScanner(readPort)

	for err == nil {
		msg, err := scanner.NextMessage()
		if err != nil {
			log.Fatalf("%s\n", err)
		}

		fmt.Printf("Msg %d\n", msg.Number())
	}
}