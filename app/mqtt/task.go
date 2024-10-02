package mqtt

import (
	"app/model"
	custom_time "app/utils/time"
	"encoding/binary"
	"encoding/json"
	"fmt"
	db "framework/database"
	"strconv"
	"time"
)

func InstReadingLog(payload []byte, topic string) {
	start := time.Now()
	var reqData InstrumentReadingLog
	json.Unmarshal([]byte(payload), &reqData)

	// groupedSensors := groupSensors(reqData.Sensor)
	// fmt.Println(groupedSensors)

	//? *****************************************************************************
	//? *****************************************************************************

	// Store Initial Logs
	// err := StoreInstrumentReadingLog(db.GetInstance(), reqData)
	// if err != nil {
	// 	fmt.Println("ERR_STORE_LOG", err)
	// 	return
	// }

	// Get Destination Code
	destId := model.GetDestinationId(db.GetInstance(), reqData.DestinationCode)
	fmt.Printf("destination_id: %d\n", destId)
	incId := getIncinerator(db.GetInstance(), destId, reqData.IncineratorCode)
	fmt.Printf("incinerator_id: %d\n", incId)
	instId, err := getInstrument(db.GetInstance(), incId, reqData.InstrumentName, reqData.InstrumentCode)
	fmt.Printf("instrument_id: %d\n", instId)

	if err != nil {
		return
	}

	if err != nil {
		fmt.Println("ERR_CREATE_INST_INC", err)
	}

	for _, data := range reqData.Sensor {
		sensorId, err := getSensor(db.GetInstance(), instId, data.Label, data.UnitOfMeasurement, data.Measure)
		// fmt.Printf("Destination Id: %d, Incinerator Id: %d, Instrument Id: %d, Sensor Id: %d, Label: %s\n", destId, incId, instId, sensorId, data.Label)
		if err != nil {
			fmt.Printf("error getting sensor id for:- label: %s, uom: %s, measure: %s \n%v", data.Label, data.UnitOfMeasurement, data.Measure, err)
		}

		// store reading
		parsedTime := custom_time.ConvertStringToTime(data.ReadAt)
		err = storeReading(db.GetInstance(), sensorId, data.Value, parsedTime)
		if err != nil {
			fmt.Printf("error while storing sensor value for:- instrument: %d, sensor: %d, value: %s, read_at: %s \n %v", instId, sensorId, data.Value, data.ReadAt, err)
		}
	}

	// //? *****************************************************************************
	// //? *****************************************************************************
	fmt.Printf("start_at: %s, elapsed_time: %s, payload_size: %v bytes\n", start.Format("2006-01-02 15:04:05.00"), time.Since(start), binary.Size(payload))
}

type AlarmLogReq struct {
	DestinationCode string `json:"destination_code"`
	IncineratorCode string `json:"incinerator_code"`
	InstrumentName  string `json:"instrument_name"`
	InstrumentCode  string `json:"instrument_code"`

	Time          string `json:"time"`
	Priority      string `json:"priority"`
	State         string `json:"state"`
	Node          string `json:"node"`
	Group         string `json:"group"`
	Tagname       string `json:"tag_name"`
	Description   string `json:"description"`
	Type          string `json:"type"`
	Limit         string `json:"limit"`
	CurrentValue  string `json:"current_value"`
	AlarmDuration string `json:"alarm_duration"`
	Operator      string `json:"operator"`
	UnAckDuration string `json:"un_ack_duration"`
}

func HandleAlarmLog(payload []byte, topic string) {
	start := time.Now()
	reqData := AlarmLogReq{}

	json.Unmarshal([]byte(payload), &reqData)

	destId := model.GetDestinationId(db.GetInstance(), reqData.DestinationCode)
	fmt.Printf("destination_id: %d\n", destId)
	incId := getIncinerator(db.GetInstance(), destId, reqData.IncineratorCode)
	fmt.Printf("incinerator_id: %d\n", incId)
	instId, _ := getInstrument(db.GetInstance(), incId, reqData.InstrumentName, reqData.InstrumentCode)
	fmt.Printf("instrument_id: %d\n", instId)

	time_lct, err := time.Parse("1/2/2006 3:04:05", reqData.Time)

	if err != nil {
		fmt.Println("ERR_TIME_LCT", reqData.Time, err)
		time_lct = time.Now()
	}
	priority, _ := strconv.Atoi(reqData.Priority)

	// alarmDuration, _ := utils_time.ParseDurationToBigInt(reqData.AlarmDuration)
	// unAckDuration, _ := utils_time.ParseDurationToBigInt(reqData.UnAckDuration)
	alarmDuration := reqData.AlarmDuration
	unAckDuration := reqData.UnAckDuration

	fmt.Println(reqData.AlarmDuration)
	fmt.Println(alarmDuration)
	fmt.Println(reqData.UnAckDuration)
	fmt.Println(unAckDuration)

	alarmModel := model.AlarmLog{
		DestinationId: destId,
		InstrumentId:  instId,
		TimeLct:       time_lct,
		Priority:      priority,
		State:         reqData.State,
		Node:          reqData.Node,
		Group:         reqData.Group,
		TagName:       reqData.Tagname,
		Description:   reqData.Description,
		Type:          reqData.Type,
		Limit:         reqData.Limit,
		CurrentValue:  reqData.CurrentValue,
		Operator:      reqData.Operator,
		AlarmDuration: alarmDuration,
		UnAckDuration: unAckDuration,
	}

	storeAlarmLog(db.GetInstance(), alarmModel)
	fmt.Printf("start_at: %s, elapsed_time: %s, payload_size: %v bytes\n", start.Format("2006-01-02 15:04:05.00"), time.Since(start), binary.Size(payload))
}
