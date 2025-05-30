package models

type Buffer struct {
	FixedStart        string
	ProtocolVersion   string
	FrameID           string
	Imei              string
	DataLength        string
	CodecID           string
	NumberOfRecords   uint16
	EventCode         string
	TimeStamp         string
	Latitude          string
	Longitude         string
	Status            string
	Satellites        string
	CheckSum          string
}