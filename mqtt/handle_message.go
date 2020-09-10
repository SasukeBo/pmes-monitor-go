package mqtt

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SasukeBo/log"
	"github.com/SasukeBo/pmes-device-monitor/orm"
	"github.com/SasukeBo/pmes-device-monitor/websocket"
	"math"
	"strings"
)

const (
	deviceStatusRunning          = 16
	deviceStatusStopped          = 17
	deviceStatusRunningWithError = 18
	deviceStatusOffline          = 32
	deviceStatusStoppedWithError = 33
)

type analyzeResult struct {
	Status     int
	Total      int
	Ng         int
	ErrorIndex []int
}

var ErrIllegalPayload = errors.New("illegal payload length")

func handleMessage(payload string) {
	payload = strings.TrimSpace(payload)
	if len(payload) < 52 {
		log.Errorln(ErrIllegalPayload)
		return
	}

	mac := payload[0:12]
	var device orm.Device
	if err := device.GetByMAC(mac); err != nil {
		log.Errorln(err)
		return
	}

	result, err := analyzeMessage(payload[12:])
	if err != nil {
		log.Errorln(err)
		return
	}
	message, err := json.Marshal(&result)
	if err == nil {
		websocket.Publish(fmt.Sprintf("device_%v", device.ID), message)
	}

	if result.Status == orm.DeviceStatusError {
		fmt.Printf("[ErrorCode]: %v\n", payload)
	}

	var produceLog orm.DeviceProduceLog
	produceLog.Record(&device, result.Total, result.Ng)
	var statusLog orm.DeviceStatusLog
	statusLog.Record(&device, result.Status, result.ErrorIndex)
}

func analyzeMessage(msg string) (result analyzeResult, err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("analyzeMessage %s\n failed: %v\n", msg, err)
		}
	}()
	words, err := hexToWords(msg)
	if err != nil {
		return
	}

	if len(words) < 5 {
		err = ErrIllegalPayload
		return
	}

	// 状态
	result.Status = wordToStatus(words[0])
	// 产量
	result.Total = wordsToAmount(words[1:3])
	// 不良
	result.Ng = wordsToAmount(words[3:5])

	if len(words) > 9 {
		result.ErrorIndex = wordsToErrorIdxs(words[9:])
	}
	return
}

func wordsToErrorIdxs(words [][]byte) []int {
	var idxs []int
	for i, word := range words {
		j := bytesToInt(word)
		if j == 0 {
			continue
		}
		for k := 0; k < 16; k++ {
			compare := int(math.Pow(2, float64(k)))
			if j&compare == compare {
				idxs = append(idxs, i*16+k)
			}
		}
	}

	return idxs
}

func wordsToAmount(words [][]byte) int {
	if len(words) != 2 {
		return 0
	}
	var amountBytes []byte
	amountBytes = append(amountBytes, words[1]...)
	amountBytes = append(amountBytes, words[0]...)
	return bytesToInt(amountBytes)
}

func wordToStatus(word []byte) int {
	statusCode := bytesToInt(word)
	var status int
	switch statusCode {
	case deviceStatusRunning:
		status = orm.DeviceStatusRunning
		fmt.Println("status: Running")
	case deviceStatusStopped:
		status = orm.DeviceStatusStopped
		fmt.Println("status: Stopped")
	case deviceStatusStoppedWithError, deviceStatusRunningWithError:
		status = orm.DeviceStatusError
		fmt.Println("status: Error")
	case deviceStatusOffline:
		status = orm.DeviceStatusShutdown
		fmt.Println("status: Offline")
	}
	return status
}

func hexToWords(hexStr string) ([][]byte, error) {
	var length = len(hexStr)
	var words [][]byte
	for i := 0; i < length; i = i + 4 {
		if length < i+4 {
			break
		}
		word, err := hex.DecodeString(hexStr[i : i+4])
		if err != nil {
			return words, err
		}

		words = append(words, word)
	}

	return words, nil
}

func bytesToInt(bytes []byte) int {
	var result int
	var length = len(bytes)
	for i := 0; i < length; i++ {
		result = result + int(bytes[i])*int(math.Pow(16*16, float64(length-i-1)))
	}

	return result
}
