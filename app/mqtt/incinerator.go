package mqtt

import (
	"app/model"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
	"gorm.io/gorm"
)

type InstrumentReadingLog struct {
	DestinationCode string   `json:"destination_code"`
	IncineratorCode string   `json:"incinerator_code"`
	InstrumentName  string   `json:"instrument_name"`
	InstrumentCode  string   `json:"instrument_code"`
	Sensor          []Sensor `json:"sensor"`
}

type Sensor struct {
	InstrumentId      uint   `json:"instrument_id"`
	Label             string `json:"label"`
	Value             string `json:"value"`
	UnitOfMeasurement string `json:"unit_of_measurement"`
	Measure           string `json:"measure"`
	ReadAt            string `json:"read_at"`
}

func registerIncinerator(db *gorm.DB, destinationId uint, incineratorCode string) uint {
	incinerator := model.Incinerator{
		DestinationId:   destinationId,
		IncineratorCode: incineratorCode,
	}

	db.Create(&incinerator)
	return incinerator.IncineratorId
}

type InstrumentR struct {
	Name string
	Code string
}

func createInstrument(db *gorm.DB, incId uint, inst InstrumentR) (uint, error) {
	newInstrument := model.Instrument{
		IncineratorId:  incId,
		InstrumentCode: inst.Code,
		InstrumentName: inst.Name,
	}
	err := db.Create(&newInstrument).Error
	if err != nil {
		return 0, err
	}
	return newInstrument.InstrumentId, nil
}

func getInstrument(db *gorm.DB, incineratorId uint, name string, code string) (uint, error) {
	var instrument model.Instrument
	db.Where(&model.Instrument{
		IncineratorId:  incineratorId,
		InstrumentName: name,
		InstrumentCode: code,
	}).Find(&instrument)

	if instrument.InstrumentId == 0 {
		instId, err := createInstrument(db, incineratorId, InstrumentR{Name: name, Code: code})
		if err != nil {
			fmt.Println("ERR_CREATE_INST", err)
			return 0, err
		}
		return instId, nil
	}

	return instrument.InstrumentId, nil
}

type SensorReq struct{}

func storeReading(db *gorm.DB, sId uint, value string, read_at time.Time) error {
	reading := model.Reading{
		SensorId: sId,
		Value:    value,
		ReadAt:   read_at,
	}
	if err := db.Create(&reading).Error; err != nil {
		return err
	}
	return nil
}

func createSensor(db *gorm.DB, instId uint, label string, uom string, measure string) (uint, error) {
	newSensor := model.Sensor{
		InstrumentId:      instId,
		Label:             label,
		UnitOfMeasurement: uom,
		Measure:           measure,
	}
	err := db.Create(&newSensor).Error
	if err != nil {
		fmt.Println("ERR_CREATE_SENSOR", err)
		return 0, err
	}
	return newSensor.SensorId, nil

}

func getSensor(db *gorm.DB, instId uint, label string, uom string, measure string) (uint, error) {
	var sensor model.Sensor
	db.Where(&Sensor{
		Label:        label,
		InstrumentId: instId,
	}).First(&sensor)

	if sensor.SensorId == 0 {
		sid, err := createSensor(db, instId, label, uom, measure)
		if err != nil {
			return 0, err
		}
		return sid, nil
	}
	return sensor.SensorId, nil
}

func getIncinerator(db *gorm.DB, destinationId uint, incineratorCode string) uint {
	var incinerator model.Incinerator
	db.Where(&model.Incinerator{DestinationId: destinationId, IncineratorCode: incineratorCode}).Find(&incinerator)
	if incinerator.IncineratorId == 0 {
		fmt.Println("incinerator not found")
		return registerIncinerator(db, destinationId, incineratorCode)
	}
	return incinerator.IncineratorId
}

func StoreInstrumentReadingLog(db *gorm.DB, irl InstrumentReadingLog) error {
	fmt.Println(irl.Sensor)
	b, err := json.Marshal(irl.Sensor)

	if err != nil {
		return errors.New("ERR_INST_READ_LOG")
	}

	fmt.Println("storing instrument reading logs...")
	go db.Create(&model.InstrumentReadingLog{
		IncineratorCode: irl.IncineratorCode,
		DestinationCode: irl.DestinationCode,
		InstrumentName:  irl.InstrumentName,
		Data: postgres.Jsonb{
			RawMessage: b,
		},
	})
	return nil
}

type MassSensorReading struct {
	IncineratorInstrumentSensorId uint
	Value                         []string
	ReadAt                        time.Time
}

/*
List all sensors id
params:
	instId uint instrument id

return:
	[]sensor_id uint
*/
type FilteredIncineratorSensor struct {
	InstrumentSensorId uint   `json:"instrument_sensor_id"`
	SensorId           uint   `json:"sensor_id"`
	Label              string `json:"label"`
}

func GetAllSensor(db *gorm.DB, instId uint) []FilteredIncineratorSensor {
	result := []FilteredIncineratorSensor{}
	tx := db.
		Select("readings.reading_id", "readings.sensor_id", "s.label").
		Table("readings").
		Joins("left join sensors as s on s.sensor_id = readings.sensor_id").
		Where("readings.instrument_id = ?", instId).Find(&result)

	if tx.Error != nil {
		fmt.Printf("ERR_GET_ALL_SENS: \n%v\n", tx.Error)
	}
	return result
}

func groupSensors(sensors []Sensor) map[string][]string {
	grouped := make(map[string][]string)

	for _, sensor := range sensors {
		label := sensor.Label
		value := sensor.Value
		measurement := sensor.UnitOfMeasurement
		grouped[label] = append(grouped[label], value, measurement)
	}
	return grouped
}
