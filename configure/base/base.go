package base

import (
	"log"
    "github.com/jacobsa/go-serial/serial"
)

const (
	UBX_SYNCH_1 = 0xB5
	UBX_SYNCH_2 = 0x62
	UBX_RTCM_1005 = 0x05;   // Stationary RTK reference ARP
	UBX_RTCM_1074 = 0x4A;   // GPS MSM4
	UBX_RTCM_1077 = 0x4D;   // GPS MSM7
	UBX_RTCM_1084 = 0x54;   // GLONASS MSM4
	UBX_RTCM_1087 = 0x57;   // GLONASS MSM7
	UBX_RTCM_1094 = 0x5E;   // Galileo MSM4
	UBX_RTCM_1097 = 0x61;   // Galileo MSM7
	UBX_RTCM_1124 = 0x7C;   // BeiDou MSM4
	UBX_RTCM_1127 = 0x7F;   // BeiDou MSM7
	UBX_RTCM_1230 = 0xE6;   // GLONASS code-phase biases, set to once every 10 seconds
	COM_PORT_UART2 = 2
	COM_PORT_USB= 3
	UBX_RTCM_MSB = 0xF5
	UBX_CLASS_CFG = 0x06
	UBX_CFG_MSG = 0x01
	UBX_CFG_TMODE3 = 0x71
	MAX_PAYLOAD_SIZE = 256
	UBX_CFG_PRT = 0x00
	UBX_CFG_CFG = 0x09
	VAL_CFG_SUBSEC_IOPORT = 0x00000001  // ioPort - communications port settings (causes IO system reset!)
	VAL_CFG_SUBSEC_MSGCONF = 0x00000002  // msgConf - message configuration

	UBX_NMEA_MSB = 0xF0 // All NMEA enable commands have 0xF0 as MSB. Equal to UBX_CLASS_NMEA
	UBX_NMEA_DTM = 0x0A // GxDTM (datum reference)
	UBX_NMEA_GAQ = 0x45 // GxGAQ (poll a standard message (if the current talker ID is GA))
  	UBX_NMEA_GBQ = 0x44 // GxGBQ (poll a standard message (if the current Talker ID is GB))
  	UBX_NMEA_GBS = 0x09 // GxGBS (GNSS satellite fault detection)
  	UBX_NMEA_GGA = 0x00 // GxGGA (Global positioning system fix data)
  	UBX_NMEA_GLL = 0x01 // GxGLL (latitude and long, whith time of position fix and status)
  	UBX_NMEA_GLQ = 0x43 // GxGLQ (poll a standard message (if the current Talker ID is GL))
  	UBX_NMEA_GNQ = 0x42 // GxGNQ (poll a standard message (if the current Talker ID is GN))
  	UBX_NMEA_GNS = 0x0D // GxGNS (GNSS fix data)
  	UBX_NMEA_GPQ = 0x40 // GxGPQ (poll a standard message (if the current Talker ID is GP))
  	UBX_NMEA_GQQ = 0x47 // GxGQQ (poll a standard message (if the current Talker ID is GQ))
  	UBX_NMEA_GRS = 0x06 // GxGRS (GNSS range residuals)
  	UBX_NMEA_GSA = 0x02 // GxGSA (GNSS DOP and Active satellites)
  	UBX_NMEA_GST = 0x07 // GxGST (GNSS Pseudo Range Error Statistics)
  	UBX_NMEA_GSV = 0x03 // GxGSV (GNSS satellites in view)
  	UBX_NMEA_RLM = 0x0B // GxRMC (Return link message (RLM))
  	UBX_NMEA_RMC = 0x04 // GxRMC (Recommended minimum data)
  	UBX_NMEA_TXT = 0x41 // GxTXT (text transmission)
  	UBX_NMEA_VLW = 0x0F // GxVLW (dual ground/water distance)
  	UBX_NMEA_VTG = 0x05 // GxVTG (course over ground and Ground speed)
  	UBX_NMEA_ZDA = 0x08 // GxZDA (Time and Date)
)

func disableAll() {
	disableRTCMCommand(UBX_RTCM_1005, COM_PORT_UART2)
	disableRTCMCommand(UBX_RTCM_1074, COM_PORT_UART2)
	disableRTCMCommand(UBX_RTCM_1084, COM_PORT_UART2)
	disableRTCMCommand(UBX_RTCM_1094, COM_PORT_UART2)
	disableRTCMCommand(UBX_RTCM_1124, COM_PORT_UART2)
	disableRTCMCommand(UBX_RTCM_1230, COM_PORT_UART2)
	saveAllConfigs()
}

func EnableAll() {
	enableRTCMCommand(UBX_RTCM_1005, COM_PORT_UART2, 1)
	enableRTCMCommand(UBX_RTCM_1074, COM_PORT_UART2, 1)
	enableRTCMCommand(UBX_RTCM_1084, COM_PORT_UART2, 1)
	enableRTCMCommand(UBX_RTCM_1094, COM_PORT_UART2, 1)
	enableRTCMCommand(UBX_RTCM_1124, COM_PORT_UART2, 1)
	enableRTCMCommand(UBX_RTCM_1230, COM_PORT_UART2, 5)
	saveAllConfigs()
}

func enableNMEA() {
	enableRTCMCommand(UBX_NMEA_GLL, COM_PORT_UART2, 1)
	enableRTCMCommand(UBX_NMEA_GSA, COM_PORT_UART2, 1)
	enableRTCMCommand(UBX_NMEA_GSV, COM_PORT_UART2, 1)
	enableRTCMCommand(UBX_NMEA_RMC, COM_PORT_UART2, 1)
	enableRTCMCommand(UBX_NMEA_VTG, COM_PORT_UART2, 1)
	enableRTCMCommand(UBX_NMEA_GGA, COM_PORT_UART2, 1)
	saveAllConfigs()
}

func setStaticPosition(ecefXOrLat int, ecefXOrLatHP int, ecefYOrLon int, ecefYOrLonHP int, ecefZOrAlt int, ecefZOrAltHP int, latLong bool) {
	cls := UBX_CLASS_CFG
	id := UBX_CFG_TMODE3
	msg_len := 40

	payloadCfg := make([]byte, 256)
	payloadCfg[2] = byte(2)

	if (latLong == true) {
    	payloadCfg[3] = (1 << 0); // Set mode to fixed. Use LAT/LON/ALT.
	}

	// Set ECEF X or Lat
	payloadCfg[4] = byte((ecefXOrLat >> 8 * 0) & 0xFF) // LSB
	payloadCfg[5] = byte((ecefXOrLat >> 8 * 1) & 0xFF)
	payloadCfg[6] = byte((ecefXOrLat >> 8 * 2) & 0xFF)
	payloadCfg[7] = byte((ecefXOrLat >> 8 * 3) & 0xFF) // MSB
  
	// Set ECEF Y or Long
	payloadCfg[8] = byte((ecefYOrLon >> 8 * 0) & 0xFF) // LSB
	payloadCfg[9] = byte((ecefYOrLon >> 8 * 1) & 0xFF)
	payloadCfg[10] = byte((ecefYOrLon >> 8 * 2) & 0xFF)
	payloadCfg[11] = byte((ecefYOrLon >> 8 * 3) & 0xFF) // MSB
  
	// Set ECEF Z or Altitude
	payloadCfg[12] = byte((ecefZOrAlt >> 8 * 0) & 0xFF) // LSB
	payloadCfg[13] = byte((ecefZOrAlt >> 8 * 1) & 0xFF)
	payloadCfg[14] = byte((ecefZOrAlt >> 8 * 2) & 0xFF)
	payloadCfg[15] = byte((ecefZOrAlt >> 8 * 3) & 0xFF) // MSB
  
	// Set high precision parts
	payloadCfg[16] = byte(ecefXOrLatHP)
	payloadCfg[17] = byte(ecefYOrLonHP)
	payloadCfg[18] = byte(ecefZOrAltHP)
	sendCommand(cls, id, msg_len, payloadCfg)
}

func disableRTCMCommand(messageNumber int, portId int) {
	enableRTCMCommand(messageNumber, portId, 0)
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
	//default to enable usb on with same sendRate
	payloadCfg[2 + COM_PORT_USB] = byte(sendRate)

	sendCommand(cls, id, msg_len, payloadCfg)
}

func sendCommand(cls int, id int, msg_len int, payloadCfg []byte) ([]byte){
	checksumA, checksumB := calcChecksum(cls, id, msg_len, payloadCfg)

	log.Print(byte(checksumA))

	options := serial.OpenOptions {
		PortName: "/dev/serial/by-path/platform-fd500000.pcie-pci-0000:01:00.0-usb-0:1.3:1.0", // change to base port
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