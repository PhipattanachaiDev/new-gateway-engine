package models

type FlattenedCombined struct {
	FixedStart        string
	ProtocolVersion   string
	FrameID           string
	Imei              string
	DataLength        string
	CodecID           string
	NumberOfRecords   uint16
	EventCode         uint16
	TimeStamp         string
	Latitude          string
	Longitude         string
	Status            string
	Satellites        int
	Direction         int
	Speed             int
	GsmSignal         int
	Networkfordetails string
	Battery           string
	External          string
	IsReliability     int
	AD1               string
	AD2               string
	HDOP              string
	Altitude          string
	InputPort         int
	OutputPort        int
	DetailStatus      string
	IsEngineOn        int
	Mileage           string
	RunningTime       string
	StatusInfo        string
	BaseStation       string
	Axis              string
	CheckSum          string
}

type Tracking struct {
	Networkfordetails string
	Battery           string
	External          string
	AD1               string
	AD2               string
	HDOP              string
	Altitude          string
	InputPort         int
	OutputPort        int
	Status            string
	Mileage           string
	RunningTime       string
	StatusInfo        string
	RFID              string
	BaseStation       string
	Axis              string
}
