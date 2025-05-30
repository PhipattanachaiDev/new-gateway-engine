package controllers

import "time"

type Card struct {
	DriverCode    string
	DriverName    string
	IsLifelong    string
	TypeID        string
	IssueDate     time.Time
	ExpireDate    time.Time
	SmartCardCode string
	IdCard        string
	CreateAt      time.Time
}

