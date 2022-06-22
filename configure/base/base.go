package base

import (
	"log"
    "github.com/jacobsa/go-serial/serial"
)

const (
	UBX_SYNCH_1 = 0xB5
	UBX_SYNCH_2 = 0x62
	UBX_RTCM_1005 = 0x05
	UBX_RTCM_1074 = 0x4A
	UBX_RTCM_1084 = 0x4D
	UBX_RTCM_1094 = 0x5E
	UBX_RTCM_1124 = 0x7C
	UBX_RTCM_1230 = 0xE6
	COM_PORT_UART2 = 2
	UBX_RTCM_MSB = 0xF5
	UBX_CLASS_CFG = 0x06
	UBX_CFG_MSG = 0x01
	MAX_PAYLOAD_SIZE = 256
	UBX_CFG_PRT = 0x00
	UBX_CFG_CFG = 0x09
	VAL_CFG_SUBSEC_IOPORT = 0x00000001;   // ioPort - communications port settings (causes IO system reset!)
	VAL_CFG_SUBSEC_MSGCONF = 0x00000002;  // msgConf - message configuration

)

func Configure() {
	enableRTCMCommand(05, COM_PORT_UART2, 1)
	enableRTCMCommand(74, COM_PORT_UART2, 1)
	enableRTCMCommand(84, COM_PORT_UART2, 1)
	enableRTCMCommand(94, COM_PORT_UART2, 1)
	saveAllConfigs()
}

func getPortSettings() ([]byte){
	cls := UBX_CLASS_CFG
	id := UBX_CFG_PRT
	msg_len := 1

	payloadCfg := make([]byte, MAX_PAYLOAD_SIZE)
	payloadCfg[0] = COM_PORT_UART2 //set portId in payloadCfg

	payloadCfg = sendCommand(cls, id, msg_len, payloadCfg)
	return payloadCfg
}

func enableRTCMCommand(messageNumber int, portId int, sendRate int) {
	//dont use current port settings actually
	payloadCfg := make([]byte, 256)

	cls := UBX_CLASS_CFG
	id := UBX_CFG_MSG
	msg_len := 8

	payloadCfg[0] = byte(UBX_RTCM_MSB)
	payloadCfg[1] = byte(messageNumber)
	payloadCfg[2 + portId] = byte(sendRate)

	sendCommand(cls, id, msg_len, payloadCfg)
}

func sendCommand(cls int, id int, msg_len int, payloadCfg []byte) ([]byte){
	checksumA, checksumB := calcChecksum(cls, id, msg_len, payloadCfg)

	log.Print(byte(checksumA))

	options := serial.OpenOptions {
		PortName: "/dev/serial/by-path/platform-fd500000.pcie-pci-0000:01:00.0-usb-0:1.4:1.0", // change to base port
		BaudRate: 115200,
		DataBits: 8,
		StopBits: 1,
		MinimumReadSize: 1,
	}

	// Open the port.
	writePort, err := serial.Open(options)
	if err != nil {
		log.Fatalf("serial.Open: %v", err)
	}
	defer writePort.Close()

	//build packet to send over serial
	byteSize := msg_len + 8 //header+checksum+payload
	packet := make([]byte, byteSize)

	//header bytes
	packet[0] = byte(UBX_SYNCH_1)
	packet[1] = byte(UBX_SYNCH_2)
	packet[2] = byte(cls)
	packet[3] = byte(id)
	packet[4] = byte(msg_len & 0xFF) //LSB
	packet[5] = byte(msg_len >> 8) //MSB

	ind := 6
	for i:=0; i<msg_len; i++ {
		packet[ind+i] = payloadCfg[i]
	}
	packet[len(packet)-1] = byte(checksumB)
	packet[len(packet)-2] = byte(checksumA)

	log.Print("Writing: %s", packet)
	writePort.Write(packet)

	//then wait to capture a byte 
	buf := make([]byte, 256)
	n, err := writePort.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	return buf[:n]
} 

func saveAllConfigs() {
	cls := UBX_CLASS_CFG
	id := UBX_CFG_CFG
	msg_len := 12

	payloadCfg:= make([]byte, 256)

	payloadCfg[4] = 0xFF
	payloadCfg[5] = 0xFF

	sendCommand(cls, id, msg_len, payloadCfg)
	
}

func calcChecksum(cls int, id int, msg_len int, payload []byte) (checksumA int, checksumB int){
	checksumA = 0
	checksumB = 0

	checksumA += cls
	checksumB += checksumA

	checksumA += id
	checksumB += checksumA

	checksumA += (msg_len & 0xFF)
	checksumB += checksumA

	checksumA += (msg_len >> 8)
	checksumB += checksumA

	for i:=0; i<msg_len; i++ {
		checksumA += int(payload[i])
		checksumB += checksumA
	}
	return checksumA, checksumB
}