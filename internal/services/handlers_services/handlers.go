package handles

import (
	"encoding/binary"
	"fmt"
	m "naturelink/models/tracker_models"
)

type Handler struct {
	Tracking m.Tracking
}

func NewHandler() *Handler {
	return &Handler{
		Tracking: m.Tracking{},
	}
}

// 1-Byte IO
func (h *Handler) Parse1ByteIO(count int, id byte, value byte) (m.Tracking, error) {
	switch id {
	case 0x1A:
		switch value {
		case 0x00:
			h.Tracking.Networkfordetails = "2G"
		case 0x01:
			h.Tracking.Networkfordetails = "3G"
		case 0x02:
			h.Tracking.Networkfordetails = "4G"
		case 0x03:
			h.Tracking.Networkfordetails = "LTE CAT-M1"
		case 0x04:
			h.Tracking.Networkfordetails = "LTE CAT-NB"
		case 0xFF:
			h.Tracking.Networkfordetails = "Unknown"
		default:
			h.Tracking.Networkfordetails = "Unknown"
		}
	}

	return h.Tracking, nil
}

// 2-Byte IO
func (h *Handler) Parse2ByteIO(count int, id byte, value []byte) (m.Tracking, error) {
	switch id {
	case 0x01:
		batteryValue := binary.LittleEndian.Uint16(value) 
		h.Tracking.Battery = fmt.Sprintf("%d", batteryValue) 
		// result = h.Tracking.Battery
	case 0x02:
		externalValue := binary.LittleEndian.Uint16(value)
		h.Tracking.External = fmt.Sprintf("%d", externalValue)
		// result = h.Tracking.External
	case 0x03:
		Ad1Value := binary.LittleEndian.Uint16(value)
		h.Tracking.AD1 = fmt.Sprintf("%d", Ad1Value)
		// result = h.Tracking.AD1
	case 0x04:
		Ad2Value := binary.LittleEndian.Uint16(value)
		h.Tracking.AD2 = fmt.Sprintf("%d", Ad2Value)
		// result = h.Tracking.AD2
	case 0x0B:
		hdopValue := binary.LittleEndian.Uint16(value)
		h.Tracking.HDOP = fmt.Sprintf("%d", hdopValue)
		// result = h.Tracking.HDOP
	case 0x0C:
		altitudeValue := binary.LittleEndian.Uint16(value)
		h.Tracking.Altitude = fmt.Sprintf("%d", altitudeValue)
		// result = h.Tracking.Altitude
	case 0x0F:
		inputValue := binary.LittleEndian.Uint16(value)
		h.Tracking.InputPort = int(inputValue)
		// result = h.Tracking.InputPort
	case 0x10:
		outputValue := binary.LittleEndian.Uint16(value)
		h.Tracking.OutputPort = int(outputValue)
		// result = h.Tracking.OutputPort
	case 0x12:
		statusValue := binary.LittleEndian.Uint16(value)
		h.Tracking.Status = fmt.Sprintf("%d", statusValue)
		// result = h.Tracking.Status
	}

	return h.Tracking, nil
}

// 4-Byte IO
func (h *Handler) Parse4ByteIO(count int, id byte, value []byte) (m.Tracking, error) {
	switch id {
	case 0x0D:
		mileageValue := binary.LittleEndian.Uint32(value)
		h.Tracking.Mileage = fmt.Sprintf("%d", mileageValue)
		// result = h.Tracking.Mileage
	case 0x0E:
		runningTimeValue := binary.LittleEndian.Uint32(value)
		h.Tracking.RunningTime = fmt.Sprintf("%d", runningTimeValue)
		// result = h.Tracking.RunningTime
	}
	return h.Tracking, nil
}

// 8-Byte IO (ข้ามเฉยๆ สำหรับตอนนี้)
func (h *Handler) Parse8ByteIO(count int, id byte, value []byte) (m.Tracking, error) {
	switch id {
	default:
	}

	return h.Tracking, nil
}

// Variable-length IO
func (h *Handler) ParseNotFixed(count int, id byte, value []byte, paramlen int) (m.Tracking, error) {
	var result string
	switch id {
	case 0x11:
		if paramlen != len(value) {
			fmt.Printf("Invalid length for ID 0x11: expected %d, got %d\n", paramlen, len(value))
			// return fmt.Println("Invalid length for ID 0x11"), nil
		}
		mcc := binary.LittleEndian.Uint16(value[0:2])
		mnc := binary.LittleEndian.Uint16(value[2:4])
		lac := binary.LittleEndian.Uint16(value[4:6])
		ci := binary.LittleEndian.Uint32(value[6:10])

		result = fmt.Sprintf("MCC: %d| MNC: %d| LAC: 0x%04X| CI: 0x%08X", mcc, mnc, lac, ci)
		h.Tracking.BaseStation = result
	case 0x18:
		if paramlen != len(value) {
			fmt.Printf("Invalid length for ID 0x18: expected %d, got %d\n", paramlen, len(value))
			// return ""
		}
		x := int16(binary.LittleEndian.Uint16(value[0:2]))
		y := int16(binary.LittleEndian.Uint16(value[2:4]))
		z := int16(binary.LittleEndian.Uint16(value[4:6]))

		result = fmt.Sprintf("X: %d mg, Y: %d mg, Z: %d mg", x, y, z)
		h.Tracking.Axis = result

	}
	return h.Tracking, nil
}


func (h *Handler) Result() m.Tracking {
	return h.Tracking
}