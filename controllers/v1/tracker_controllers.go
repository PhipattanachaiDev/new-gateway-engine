package controllers

import (
	"context"
	"fmt"
	"naturelink/configs/database"
	models "naturelink/models/tracker_models"
	"time"
)

func ProcessTrackerController(data *models.FlattenedCombined) error {
	tableName, err := GenerateTableName(data.TimeStamp)
	if err != nil {
		return fmt.Errorf("failed to generate table name: %v", err)
	}

	trackerID, err := fetchTrackerID(data.Imei)
	if err != nil {
		return err
	}

	err = insertTrackerEvent(trackerID, tableName, data)
	if err != nil {
		return fmt.Errorf("insert tracker event error: %w", err)
	}

	return nil

}

// เหลือ เพิ่ม IsEngineOn
func insertTrackerEvent(trackerID int, tableName string, data *models.FlattenedCombined) error {
	dbRt, err := database.ConnDBRealtimeConfig()
	if err != nil {
		return fmt.Errorf("failed to connect to Realtime database: %v", err)
	}

	queryInsert := `SELECT public.insert_event(
		$1::SMALLINT, $2::TIMESTAMP, $3::SMALLINT, $4::SMALLINT,
		$5::NUMERIC, $6::NUMERIC, $7::INTEGER,
		$8::SMALLINT, $9::SMALLINT, $10::NUMERIC,
		$11::SMALLINT, $12::SMALLINT,
		$13::BIT(1), $14::BIT(1), $15::TEXT`
	_, err = dbRt.Exec(context.Background(), queryInsert, trackerID, data.TimeStamp, data.EventCode, data.Speed, data.Latitude, data.Longitude, data.Mileage, data.Direction, data.Altitude, data.HDOP, data.Satellites, data.GsmSignal, tableName)
	if err != nil {
		return fmt.Errorf("failed to execute insert query: %v", err)
	}
	return nil
}

func fetchTrackerID(imei string) (int, error) {
	db, err := database.ConnDBEzViewLiteConfig()
	if err != nil {
		return 0, fmt.Errorf("connect to EZViewLite DB failed: %w", err)
	}
	defer db.Close()

	var trackerID int
	query := `SELECT ezv_get_tracker_id_by_imei($1)`
	err = db.QueryRow(context.Background(), query, imei).Scan(&trackerID)
	if err != nil {
		return 0, fmt.Errorf("query tracker ID failed: %w", err)
	}
	if trackerID == 0 {
		return 0, fmt.Errorf("tracker ID not found for IMEI: %s", imei)
	}

	return trackerID, nil
}

func GenerateTableName(timeStamp string) (string, error) {
	parsedDatetime, err := time.Parse("2006-01-02 15:04:05", timeStamp)
	if err != nil {
		return "", fmt.Errorf("error parsing generate table datetime: %v", err)
	}

	formattedDatetime := parsedDatetime.Format("20060102")
	tableName := fmt.Sprintf("tracker_location_%s", formattedDatetime)
	return tableName, nil
}
