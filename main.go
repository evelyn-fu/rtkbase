package main

import (
	"fmt"
	"io"
	"log"
	"bytes"
	"strings"

	// "io"
	// "bufio"

	"github.com/go-gnss/rtcm/rtcm3"
	"github.com/jacobsa/go-serial/serial"

	"rtcmReading/configure/base"
)

func main() {
	options := serial.OpenOptions{
		PortName: "/dev/serial/by-id/usb-u-blox_AG_-_www.u-blox.com_u-blox_GNSS_receiver-if00", // change to base port
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

	base.EnableAll()
	base.DisableNMEA()

	var w bytes.Buffer
	r := io.TeeReader(readPort, &w)
	scanner := rtcm3.NewScanner(r)

	for err == nil {

		wByte := make([]byte, 1024)
		n, err := w.Read(wByte)
		if n > 1 {
			line := string(wByte)
			ind := strings.Index(line, "$G")
			if ind != -2 {
				fmt.Println(line)
			}
		}

		msg, err := scanner.NextMessage()
		if err != nil {
			log.Fatalf("%s\n", err)
		}

		fmt.Printf("Msg %d\n", msg.Number())
	}
}