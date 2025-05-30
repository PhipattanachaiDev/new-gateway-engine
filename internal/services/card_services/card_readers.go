package cardservices

import (
	"fmt"
	"math"
	log "naturelink/internal/services/log_services"
	models "naturelink/models/card_models"
	"strings"
)

func ProcessCardReaders(tmpMVT []string) string {
    // fmt.Println("🚀🚀 ~ funcProcessCardReaders ~ tmpMVT : ", tmpMVT)
	var card = models.CardControllers{}

	var BolRFID, BolQRCode bool
	var CompareCard, tmpTrack2IdCard, tmpTrack3Prefix, tmpTrack3Name, tmpTrack3SurName, tmpStrDriverName string
	var bolTrack3 bool
	var strExternalData string

	str4Replace := strings.ReplaceAll(tmpMVT[19], " ", "_")
    // fmt.Println("🚀🚀 ~ funcProcessCardReaders ~ str4Replace : ", str4Replace)
	str4Replace = strings.ReplaceAll(str4Replace, "\x00", "_")
    // fmt.Println("🚀🚀 ~ funcProcessCardReaders ~ str4Replace : ", str4Replace)
	var _tmpMVT []string

	if strings.HasPrefix(str4Replace, "%") {
		_tmpMVT = strings.Split(str4Replace, "?")
        fmt.Println("🚀🚀 ~ if strings.HasPrefix ~ _tmpMVT : ", _tmpMVT)
	} else {
		_tmpMVT = strings.Split(str4Replace, " ")
        fmt.Println("🚀🚀 ~ else strings.HasPrefix ~ _tmpMVT : ", _tmpMVT)
	}

	if len(_tmpMVT) > 1 {
		if strings.HasPrefix(_tmpMVT[0], "%") {
			bolTrack3 = true
			tmpTrack3 := strings.Split(_tmpMVT[0], "$")
			tmpTrack2 := strings.Split(_tmpMVT[1], "=")
			tmpTrack2IdCard := strings.ReplaceAll(tmpTrack2[0], ";", "")

			if len(tmpTrack2IdCard) > 6 {
				tmpTrack2IdCard = tmpTrack2IdCard[6:]
			} else {
				tmpTrack2IdCard = tmpTrack2IdCard[:0]
			}

			card.IdCard = tmpTrack2IdCard

			if len(tmpTrack3) > 1 {
				tmpTrack3SurName = strings.ReplaceAll(tmpTrack3[0], "%", "")
				tmpTrack3SurName = strings.ReplaceAll(tmpTrack3SurName, "^", "")
				tmpTrack3SurName = strings.ReplaceAll(tmpTrack3SurName, " ", "")
				tmpTrack3SurName = strings.ReplaceAll(tmpTrack3SurName, "_", "")

				tmpTrack3Name = tmpTrack3[1]
				tmpTrack3Prefix = strings.ReplaceAll(tmpTrack3[2], "^", "")
				tmpTrack3Prefix = strings.ReplaceAll(tmpTrack3Prefix, "\n", "")
				tmpTrack3Prefix = strings.ReplaceAll(tmpTrack3Prefix, "\r", "")
				tmpTrack3SurName = strings.ReplaceAll(tmpTrack3SurName, "\n", "")
				tmpTrack3SurName = strings.ReplaceAll(tmpTrack3SurName, "\r", "")

				tmpStrDriverName = tmpTrack3Prefix + tmpTrack3Name + " " + tmpTrack3SurName
				tmpStrDriverName = strings.ReplaceAll(tmpStrDriverName, "^", "")
				tmpStrDriverName = strings.ReplaceAll(tmpStrDriverName, "?", "")
				tmpStrDriverName = strings.ReplaceAll(tmpStrDriverName, "\n", "")
			} else if len(tmpTrack3) == 1 {
				tmpDName := strings.Split(tmpTrack3[0], "_")
				tmpStrDriverName = tmpDName[3] + "." + tmpDName[2] + " " + tmpDName[1]
			}

			card.DriverName = tmpStrDriverName
            // fmt.Println("🚀🚀 ~ Card.DriverName : ", Card.DriverName)

			tmpCardID := strings.Split(_tmpMVT[2], "_")
            // fmt.Println("🚀🚀 ~ if strings.HasPrefix ~ tmpCardID : ", tmpCardID)

			if len(tmpCardID) > 28 {
				CompareCard = strings.TrimSpace(tmpCardID[2]) + strings.TrimSpace(tmpCardID[14]) + strings.TrimSpace(tmpCardID[26]) + strings.TrimSpace(tmpCardID[28])
			} else {
				CompareCard = strings.Join(tmpCardID, ",")
			}

			CompareCard = strings.ReplaceAll(CompareCard, "?", "")
			CompareCard = strings.ReplaceAll(CompareCard, "\n", "")
			CompareCard = strings.ReplaceAll(CompareCard, "\r", "")

			card.TypeID = tmpCardID[2]
			card.DriverCode = tmpCardID[26]
			card.IsLifelong = tmpCardID[14]
			card.SmartCardCode = CompareCard

		} else {
			BolRFID = true
			if strings.HasPrefix(_tmpMVT[0], "QR[") {
				BolQRCode = true
				CompareCard = strings.TrimSpace(_tmpMVT[0][strings.Index(_tmpMVT[0], ":")+1:])
				_tmpMVT[0] = CompareCard

			} else {
				BolQRCode = false
				CompareCard = fmt.Sprintf("%d", convertHexToDec(strings.TrimSpace(_tmpMVT[0])))
			}
		}
	}

	if BolRFID {
		if BolQRCode {
			strExternalData = fmt.Sprintf("%s,%s,-,-,-,-,-", strExternalData, strings.TrimSpace(_tmpMVT[0]))
		} else {
			strExternalData = fmt.Sprintf("%s,%d,-,-,-,-,-", strExternalData, convertHexToDec(strings.TrimSpace(_tmpMVT[0])))
		}
	} else {
		if bolTrack3 {
			strExternalData = fmt.Sprintf("%s,%s,%s,%s", strExternalData, strings.TrimSpace(CompareCard), strings.TrimSpace(tmpStrDriverName), tmpTrack2IdCard)
		} else {
			tmpCard := fmt.Sprintf("%s,%s,%s,%s,-,-", strings.TrimSpace(_tmpMVT[0]), strings.TrimSpace(_tmpMVT[12]), strings.TrimSpace(_tmpMVT[24]), strings.TrimSpace(_tmpMVT[26]))
			tmpCard = strings.ReplaceAll(tmpCard, "\n", "")
			tmpCard = strings.ReplaceAll(tmpCard, "\r", "")
			tmpCard = strings.ReplaceAll(tmpCard, "$", "")
			tmpCard = strings.ReplaceAll(tmpCard, "\"", "")
			tmpCard = strings.ReplaceAll(tmpCard, " ", "")
			tmpCard = strings.ReplaceAll(tmpCard, "_", "")
			tmpCard = strings.ReplaceAll(tmpCard, ";", "")
			tmpCard = strings.ReplaceAll(tmpCard, "=", "")
			tmpCard = strings.ReplaceAll(tmpCard, ")", "")
			tmpCard = strings.ReplaceAll(tmpCard, "(", "")
			tmpCard = strings.ReplaceAll(tmpCard, "*", "")
			tmpCard = strings.ReplaceAll(tmpCard, ".", "")
			tmpCard = strings.ReplaceAll(tmpCard, "?", "")
			tmpCard = strings.ReplaceAll(tmpCard, "^", "")
			tmpCard = strings.ReplaceAll(tmpCard, "%", "")

			strExternalData = fmt.Sprintf("%s,%s", strExternalData, tmpCard)
		}
	}
	return strExternalData
}

func convertHexToDec(hexString string) int64 {
	var result int64
	length := len(hexString)

	for i := 0; i < length; i++ {
		char := hexString[i]
		var value int64

		if char >= '0' && char <= '9' {
			value = int64(char - '0')
		} else if char >= 'A' && char <= 'F' {
			value = int64(char - 'A' + 10)
		} else {
			log.ServicesWarning(fmt.Errorf("invalid hex digit: %c", char))
			continue
		}

		result += value * int64(math.Pow(16, float64(length-i-1)))
	}

	return result
}
