package parsedata

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	// "naturelink/controllers/v1"
	handles "naturelink/internal/services/handlers_services"
	// log "naturelink/internal/services/log_services"
	m "naturelink/models/buffer_models"
	models "naturelink/models/tracker_models"
	"strings"
	"time"
)

type Info struct {
	EventCode  uint16
	TimeStamp  string
	Latitude   string
	Longitude  string
	Status     string
	Satellites int
	Angle      int
	Speed      int
	Signal     int
}

func ParseData(buffer string) string {
	hexData := strings.ReplaceAll(buffer, " ", "")
	bytes, err := hex.DecodeString(hexData)

	// fmt.Printf("🚀🚀 ~ funcParseData ~ bytes : %02X\n", bytes)

	if err != nil {
		return ""
	}

	gps := parseHeader(bytes)

	pos := 16

	fmt.Printf("🚀🚀 ~ funcParseData ~ gps : %+v\n", gps.NumberOfRecords)
	// var flatResults []FlattenedCombined

	for i := 0; i < int(gps.NumberOfRecords); i++ {
		info, posAfterInfo := parseBasicInfo(bytes, pos)
		detail, posAfterDetail := parseTrackingDetail(bytes, posAfterInfo)

		// fmt.Println("Input Port:", detail.InputPort)
		// fmt.Println("Output Port:", detail.OutputPort)

		// if detail.InputPort == 1 {
		// 	var IsEngineOn = 1
		// 	fmt.Println("🚀🚀 ~ fori:=0;i<int ~ IsEngineOn : ", IsEngineOn)
		// }

		isEngineOn := CheckStstusEngine(detail.InputPort, detail.OutputPort)
		// fmt.Println("🚀🚀 ~ isEngineOn : ", isEngineOn)

		flat := models.FlattenedCombined{
			FixedStart:        gps.FixedStart,
			ProtocolVersion:   gps.ProtocolVersion,
			FrameID:           gps.FrameID,
			Imei:              gps.Imei,
			DataLength:        gps.DataLength,
			CodecID:           gps.CodecID,
			NumberOfRecords:   gps.NumberOfRecords,
			EventCode:         info.EventCode,
			TimeStamp:         info.TimeStamp,
			Latitude:          info.Latitude,
			Longitude:         info.Longitude,
			Status:            info.Status,
			Satellites:        info.Satellites,
			Direction:         info.Angle,
			Speed:             info.Speed,
			GsmSignal:         info.Signal,
			Networkfordetails: detail.Networkfordetails,
			Battery:           detail.Battery,
			External:          detail.External,
			AD1:               detail.AD1,
			AD2:               detail.AD2,
			HDOP:              detail.HDOP,
			Altitude:          detail.Altitude,
			IsEngineOn:        isEngineOn,
			// InputPort:         detail.InputPort,
			// OutputPort:        detail.OutputPort,
			DetailStatus: detail.Status,
			Mileage:      detail.Mileage,
			RunningTime:  detail.RunningTime,
			StatusInfo:   detail.StatusInfo,
			BaseStation:  detail.BaseStation,
			Axis:         detail.Axis,
			CheckSum:     gps.CheckSum,
		}

		fmt.Printf("🚀🚀 ~ flat : %v\n", flat)
		// flatResults = append(flatResults, flat)
		pos = posAfterDetail // update pos for next record

		if detail.RFID != "" {
			fmt.Println("RFID:", detail.RFID)
		}
		// else {
		// 	err := controllers.ProcessTrackerController(&flat)
		// 	if err != nil {
		// 		log.ServicesError(fmt.Errorf("error processing tracker controller: %w", err))
		// 	}
		// }
	}

	// fmt.Println("🚀🚀 ~ flatResults : ", flatResults)
	// printFlattened(flatResults)

	if len(bytes) >= 2 {
		gps.CheckSum = fmt.Sprintf("%02X", bytes[len(bytes)-2])
	}

	return gps.Imei
}

// ---------- ส่วนฟังก์ชันย่อย ----------

func parseHeader(bytes []byte) m.Buffer {
	return m.Buffer{
		FixedStart:      fmt.Sprintf("%02X%02X", bytes[0], bytes[1]),
		ProtocolVersion: fmt.Sprintf("%02X", bytes[2]),
		FrameID:         fmt.Sprintf("%02X", bytes[3]),
		Imei:            extractIMEI(bytes[4:12]),
		DataLength:      fmt.Sprintf("%d", int(bytes[12])|int(bytes[13])<<8),
		CodecID:         fmt.Sprintf("%02X", bytes[14]),
		NumberOfRecords: uint16(bytes[15]),
	}
}

func extractIMEI(b []byte) string {
	var imei string
	for _, byteVal := range b {
		high := byteVal >> 4
		low := byteVal & 0x0F
		if len(imei) > 0 || high != 0 {
			imei += fmt.Sprintf("%X", high)
		}
		if len(imei) > 0 || low != 0 {
			imei += fmt.Sprintf("%X", low)
		}
	}
	return imei
}

func parseBasicInfo(bytes []byte, pos int) (Info, int) {
    // fmt.Printf("🚀🚀 ~ funcparseBasicInfo ~ bytes : %02X\n", bytes)
	var info Info

	info.EventCode = uint16(bytes[pos]) | uint16(bytes[pos+1])<<8
    // fmt.Println("🚀🚀 ~ funcparseBasicInfo ~ info.EventCode : ", info.EventCode)
	// fmt.Printf("Byte Pos :%02X\n",bytes[pos])
	pos += 2

	// d := bytes[pos:]
    // fmt.Printf("🚀🚀 ~ funcparseBasicInfo ~ d : %02X", d)

	// ts := int64(bytes[pos]) | int64(bytes[pos+1])<<8 | int64(bytes[pos+2])<<16 | int64(bytes[pos+3])<<24
	bytes = bytes[pos:pos+4]
	fmt.Printf("Byte Pos :%02X\n", bytes[pos:pos+4])
	secs := binary.LittleEndian.Uint32(bytes)

	// คำนวณเวลา
	base := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	dt := base.Add(time.Duration(secs) * time.Second).Format("2006-01-02 15:04:05")

	fmt.Println("🚀🚀 ~ func main ~ dt :", dt)
	// bytes, err := hex.DecodeString(fmt.Sprint(ts))
	// if err != nil {
	// 	fmt.Println("err",err)
	// }
	// fmt.Printf("debug data: %02X",bytes)
    // fmt.Printf("🚀🚀 ~ funcparseBasicInfo ~ ts : %d\n", ts)
	// t := time.Unix(ts+946684800, 0).UTC()
	// info.TimeStamp = t.Format("2006-01-02 15:04:05")
	pos += 4

	rawLat := int32(binary.LittleEndian.Uint32(bytes[pos : pos+4]))
	info.Latitude = fmt.Sprintf("%.6f", float64(rawLat)/1_000_000)
	pos += 4

	rawLon := int32(binary.LittleEndian.Uint32(bytes[pos : pos+4]))
	info.Longitude = fmt.Sprintf("%.6f", float64(rawLon)/1_000_000)
	pos += 4

	status := binary.LittleEndian.Uint16(bytes[pos : pos+2])
	info.Status = fmt.Sprintf("0x%04X (CSQ: %d)", status, (status>>1)&0x1F)
	info.Signal = int(status & 0x1)
	info.Speed = int(status >> 6)
	pos += 2

	sat := binary.LittleEndian.Uint16(bytes[pos : pos+2])
	info.Satellites = int(sat & 0x7F)
	info.Angle = int(status >> 6)
	pos += 2

	return info, pos
}

func parseTrackingDetail(bytes []byte, pos int) (models.Tracking, int) {
	h := handles.NewHandler()

	// 1 Byte Params
	num1 := int(bytes[pos])
	pos++
	for i := 0; i < num1; i++ {
		h.Parse1ByteIO(num1, bytes[pos], bytes[pos+1])
		pos += 2
	}

	// 2 Byte Params
	num2 := int(bytes[pos])
	pos++
	for i := 0; i < num2; i++ {
		h.Parse2ByteIO(num2, bytes[pos], bytes[pos+1:pos+3])
		pos += 3
	}

	// 4 Byte Params
	num4 := int(bytes[pos])
	pos++
	for i := 0; i < num4; i++ {
		h.Parse4ByteIO(num4, bytes[pos], bytes[pos+1:pos+5])
		pos += 5
	}

	// 8 Byte Params
	num8 := int(bytes[pos])
	pos++
	for i := 0; i < num8; i++ {
		h.Parse8ByteIO(num8, bytes[pos], bytes[pos+1:pos+9])
		pos += 9
	}

	// Variable-length Params
	numVar := int(bytes[pos])
	pos++
	for i := 0; i < numVar; i++ {
		paramID := bytes[pos]
		paramLen := int(bytes[pos+1])
		paramVal := bytes[pos+2 : pos+2+paramLen]
		h.ParseNotFixed(numVar, paramID, paramVal, paramLen)
		pos += 2 + paramLen
	}

	return h.Result(), pos
}

func CheckStstusEngine(inputPort, outputPort int) int {
	if inputPort == 1 && outputPort == 0 {
		return 1
	} else if inputPort == 0 && outputPort == 1 {
		return 0
	}
	return 0
}

// func printFlattened(results []FlattenedCombined) {
// 	for _, f := range results {
// 		fmt.Printf("{%d %s %s %s %s %s %d %s %s %s %s %s %s %s %s %s %s %s %s %s}\n",
// 			f.EventCode, f.TimeStamp, f.Latitude, f.Longitude, f.Status, f.Satellites, f.Speed,
// 			f.Networkfordetails, f.Battery, f.External, f.AD1, f.AD2, f.HDOP,
// 			f.Altitude, f.InputPort, f.OutputPort, f.DetailStatus, f.Mileage,
// 			f.RunningTime, f.BaseStation,
// 		)
// 	}
// }
