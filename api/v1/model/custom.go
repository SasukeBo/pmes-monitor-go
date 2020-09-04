package model

import "time"

type DeviceType struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	ErrorCodeID int       `json:"errorCodeID"`
	CreatedAt   time.Time `json:"createdAt"`
}

type Device struct {
	ID           int       `json:"id"`
	CreatedAt    time.Time `json:"createdAt"`
	Number       string    `json:"number"`
	DeviceTypeID int       `json:"deviceTypeID"`
	Mac          string    `json:"mac"`
	Address      string    `json:"address"`
	Status       string    `json:"status"`
}
