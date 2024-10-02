package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	db "framework/database"
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
)

type JSON json.RawMessage

func (j *JSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	result := json.RawMessage{}
	err := json.Unmarshal(bytes, &result)
	*j = JSON(result)
	return err
}

func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.RawMessage(j).MarshalJSON()
}

/*
	Database Schema
	Further reference for gorm at https://gorm.io/docs/models.html
*/

type GormDefault struct {
	CreatedAt time.Time  `json:"createdAt" `
	UpdatedAt time.Time  `json:"updatedAt" gorm:"not null;default:CURRENT_TIMESTAMP;"`
	DeletedAt *time.Time `json:"deleteAt" `
}

// Incinerator
type (
	Destination struct {
		DestinationId   uint   `gorm:"primary_key;auto_increment"`
		DestinationCode string `json:"destination_code"`
		DestinationName string `json:"destination_name"`
		GormDefault
	}

	Incinerator struct {
		IncineratorId   uint `gorm:"primary_key:auto_increment"`
		IncineratorCode string
		DestinationId   uint
		Destination     Destination `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
		Line            string
		IsActive        bool `gorm:"default:TRUE"`
		GormDefault
	}

	Instrument struct {
		InstrumentId   uint `gorm:"primary_key:auto_increment"`
		IncineratorId  uint
		Incinerator    Incinerator `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
		InstrumentCode string
		InstrumentName string `gorm:"unique_index"`
		IsActive       bool   `gorm:"default:TRUE"`
		GormDefault
	}

	Sensor struct {
		SensorId          uint `gorm:"primary_key:auto_increment"`
		InstrumentId      uint
		Instrument        Instrument `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
		Label             string
		IsActive          bool `gorm:"default:TRUE"`
		UnitOfMeasurement string
		Measure           string
		GormDefault
	}

	Reading struct {
		ReadingId uint `gorm:"primary_key:auto_increment"`
		SensorId  uint
		Sensor    Sensor `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
		Value     string
		ReadAt    time.Time
		IsActive  bool `gorm:"default:true"`
		GormDefault
	}

	InstrumentReadingLog struct {
		Id              uint `gorm:"primary_key:auto_increment"`
		IncineratorCode string
		DestinationCode string
		InstrumentName  string
		Data            postgres.Jsonb `gorm:"type:jsonb;default:'{}'"`
		ReadAt          time.Time
		GormDefault
	}
)

// Alarm
type (
	AlarmLog struct {
		AlarmId       uint `gorm:"primary_key:auto_increment"`
		DestinationId uint
		Destination   Destination `gorm:"constraint:OnUpdate:CASCADE,OnDelete:Set NULL"`
		InstrumentId  uint
		Instrument    Instrument `gorm:"constraint:OnUpdate:CASCADE,OnDelete:Set NULL"`
		TimeLct       time.Time
		Priority      int
		State         string
		Node          string
		Group         string
		TagName       string
		Description   string
		Type          string
		Limit         string
		CurrentValue  string
		Operator      string
		AlarmDuration string
		UnAckDuration string
		GormDefault
	}
)

/*

 */
func Execute() {
	db.GetInstance().AutoMigrate(
		&Destination{},
		&Incinerator{},
		&Instrument{},
		&Sensor{},
		&Reading{},
		&InstrumentReadingLog{},
		&AlarmLog{},
	)

	db.GetInstance().AutoMigrate()

}
