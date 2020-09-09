// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"time"
)

type Dashboard struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	DeviceTotal  int    `json:"deviceTotal"`
	RunningTotal int    `json:"runningTotal"`
	ErrorTotal   int    `json:"errorTotal"`
}

type DashboardDevice struct {
	ID               int       `json:"id"`
	Number           string    `json:"number"`
	Status           string    `json:"status"`
	Total            int       `json:"total"`
	Ng               int       `json:"ng"`
	Durations        []int     `json:"durations"`
	Errors           []string  `json:"errors"`
	LastProduceLogID int       `json:"lastProduceLogID"`
	LastStatusLogID  int       `json:"lastStatusLogID"`
	LastStatusTime   time.Time `json:"lastStatusTime"`
}

type DashboardDeviceErrorsResponse struct {
	Category []string `json:"category"`
	Data     []int    `json:"data"`
}

type DashboardDeviceFreshResponse struct {
	ProduceLogs []*DeviceProduceLog `json:"produceLogs"`
	StatusLogs  []*DeviceStatusLog  `json:"statusLogs"`
}

type DashboardDeviceStatusResponse struct {
	Stopped int `json:"stopped"`
	Running int `json:"running"`
	Offline int `json:"offline"`
	Error   int `json:"error"`
}

type DashboardOverviewAnalyzeResponse struct {
	Total      int     `json:"total"`
	Ng         int     `json:"ng"`
	Activation float64 `json:"activation"`
}

type DashboardWrap struct {
	Total      int          `json:"total"`
	Dashboards []*Dashboard `json:"dashboards"`
}

type DeviceProduceLog struct {
	ID       int `json:"id"`
	DeviceID int `json:"deviceID"`
	Total    int `json:"total"`
	Ng       int `json:"ng"`
}

type DeviceStatusLog struct {
	ID       int      `json:"id"`
	DeviceID int      `json:"deviceID"`
	Messages []string `json:"messages"`
	Status   string   `json:"status"`
	Duration int      `json:"duration"`
}
