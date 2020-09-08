package mqtt

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/SasukeBo/log"
	"github.com/SasukeBo/pmes-device-monitor/orm"
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

func handleMessage(payload string) {
	fmt.Printf("handle message: %s\n", payload)
	mac, status, total, ng, errorIndex, err := analyzeMessage(payload)
	if err != nil {
		log.Errorln(err)
		return
	}
	var device orm.Device
	if err := device.GetByMAC(mac); err != nil {
		log.Errorln(err)
		return
	}

	var produceLog orm.DeviceProduceLog
	produceLog.Record(mac, total, ng)
	var statusLog orm.DeviceStatusLog
	statusLog.Record(mac, status, errorIndex)
}

func analyzeMessage(msg string) (mac string, status int, total int, ng int, errorIndex []int, err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("analyzeMessage %s\n failed: %v\n", msg, err)
		}
	}()
	words, err := hexToWords(strings.TrimSpace(msg))
	if err != nil {
		return
	}

	if len(words) < 8 {
		err = errors.New("illegal payload length")
		return
	}

	// mac
	mac = wordsToMAC(words[0:3])
	// 状态
	status = wordToStatus(words[3])
	// 产量
	total = wordsToAmount(words[4:6])
	// 不良
	ng = wordsToAmount(words[6:8])

	if len(words) > 8 {
		errorIndex = wordsToErrorIdxs(words[8:])
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

func wordsToMAC(words [][]byte) string {
	if len(words) < 3 {
		return ""
	}
	return fmt.Sprintf("%s%s%s", hex.EncodeToString(words[0]), hex.EncodeToString(words[1]), hex.EncodeToString(words[2]))
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
