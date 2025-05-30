package cardservices_test

import (
	// "fmt"
	// "strings"
	cardservices "naturelink/internal/services/card_services"
	"testing"
	// cardservices "naturelink/internal/services/card_services"
)

func TestProcessTrackCard(t *testing.T) {
	rawData := []string{
		"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "% ^MUNGMIN$TESTDEE$MR.^^?;6007641100833866815=180919770411=?+ 23 1 9999958 00100 ?", "",
	}
	cardservices.ProcessCardReaders(rawData)
}
