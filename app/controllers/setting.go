package controllers

import (
	"app/model"
	"app/network"
	"fmt"
	db "framework/database"

	utils "framework/utils/common"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
)

type DestinationModel struct {
	DestinationId   uint   `json:"destination_id"`
	DestinationCode string `json:"destination_code"`
	DestinationName string `json:"destination_name"`
}

func SettingRoute() *chi.Mux {
	router := chi.NewMux()
	router.Get("/destination", GetDestination)
	router.Post("/destination", StoreDestination)
	router.Get("/incinerator", GetIncinerator)
	router.Post("/incinerator", StoreIncinerator)
	// router.Patch("/incinerator", PatchIncinerator)
	router.Get("/instrument", GetInstrument)
	router.Post("/instrument", StoreInstrument)

	router.Get("/instrument/{instrument_id}", GetInstrumentData)
	return router
}

func GetDestination(w http.ResponseWriter, r *http.Request) {
	var destModel []DestinationModel
	tx := db.GetInstance().Model(model.Destination{}).Scan(&destModel)

	if tx.Error != nil {
		network.ResponseJSON(w, true, http.StatusBadRequest, destModel)
		return
	}
	network.ResponseJSON(w, false, http.StatusOK, destModel)

}

func StoreDestination(w http.ResponseWriter, r *http.Request) {
	var destModel model.Destination
	var destReq model.Destination
	err := network.ReadJSONData(r, &destReq)

	if err != nil {
		network.ResponseJSON(w, true, http.StatusBadRequest, "Invalid parameter")
		return
	}

	validateInput := []utils.InputPattern{
		{Input: destReq.DestinationCode, RegPattern: "^[a-zA-Z0-9]+$", PatternDescription: "Destination code is required."},
		{Input: destReq.DestinationName, RegPattern: "^[a-zA-Z0-9]+$", PatternDescription: "Destination name is required."},
	}

	pattern, patternErr := utils.CheckPattern(validateInput)
	if patternErr {
		network.ResponseJSON(w, true, http.StatusBadRequest, pattern)
		return
	}

	tx := db.GetInstance().Where("destination_code = ?", destReq.DestinationCode).Table("destinations").Scan(&destModel)

	if tx.RowsAffected == 0 {
		// Create
		tx = db.GetInstance().Table("destinations").Create(&destReq)
		if tx.Error != nil {
			fmt.Println(tx.Error.Error())
			network.ResponseJSON(w, true, http.StatusInternalServerError, ErrDbTransFail)
			return
		}
		network.ResponseJSON(w, false, http.StatusCreated, "Entry Added")
	} else {
		tx = db.GetInstance().Model(&destModel).Updates(&destReq)
		if tx.Error != nil {
			fmt.Println(tx.Error.Error())
			network.ResponseJSON(w, true, http.StatusInternalServerError, ErrDbTransFail)
			return
		}
		network.ResponseJSON(w, false, http.StatusCreated, "Entry Updated")
	}
}

type IncineratorModel struct {
	DestinationId   uint   `json:"destination_id"`
	DestinationCode string `json:"destination_code"`
	IncineratorId   uint   `json:"incinerator_id"`
	IncineratorCode string `json:"incinerator_code"`
}

func GetIncinerator(w http.ResponseWriter, r *http.Request) {
	var incModel []IncineratorModel
	destinationCode := r.URL.Query().Get("destination_code")

	tx := db.GetInstance().Model(model.Incinerator{}).
		Raw(`
			select i.incinerator_id, i.incinerator_code, d.destination_id, d.destination_code
			from incinerators i
			left join destinations d on d.destination_id = i.destination_id
			where d.destination_code = ?
		`, destinationCode).
		Scan(&incModel)

	if tx.Error != nil {
		network.ResponseJSON(w, true, http.StatusBadRequest, incModel)
		return
	}
	network.ResponseJSON(w, false, http.StatusOK, incModel)
}

type IncineratorReq struct {
	DestinationId   uint   `json:"destination_id"`
	IncineratorCode string `json:"incinerator_code"`
}

func StoreIncinerator(w http.ResponseWriter, r *http.Request) {
	// var incModel model.Incinerator
	var incReq IncineratorReq
	err := network.ReadJSONData(r, &incReq)

	if err != nil {
		network.ResponseJSON(w, true, http.StatusBadRequest, "Invalid parameter")
		return
	}

	validateInput := []utils.InputPattern{
		{Input: incReq.IncineratorCode, RegPattern: "^[a-zA-Z0-9]+$", PatternDescription: "Incinerator code is required."},
		// {Input: incReq.DestinationId, RegPattern: "^[a-zA-Z0-9]+$", PatternDescription: "Destination name is required."},
	}

	pattern, patternErr := utils.CheckPattern(validateInput)
	if patternErr {
		network.ResponseJSON(w, true, http.StatusBadRequest, pattern)
		return
	}

	// tx := db.GetInstance().Where("destination_id = ?", incReq.DestinationId).Table("incinerators").Scan(&incModel)

	// if tx.RowsAffected == 0 {
	// Create
	tx := db.GetInstance().Table("incinerators").Create(&incReq)
	if tx.Error != nil {
		fmt.Println(tx.Error.Error())
		network.ResponseJSON(w, true, http.StatusInternalServerError, ErrDbTransFail)
		return
	}
	network.ResponseJSON(w, false, http.StatusCreated, "Entry Added")
	// } else {
	// 	tx = db.GetInstance().Model(&incModel).Updates(&incReq)
	// 	if tx.Error != nil {
	// 		fmt.Println(tx.Error.Error())
	// 		network.ResponseJSON(w, true, http.StatusInternalServerError, ErrDbTransFail)
	// 		return
	// 	}
	// 	network.ResponseJSON(w, false, http.StatusCreated, "Entry Updated")
	// }

}

type InstrumentModel struct {
	DestinationId   uint   `json:"destination_id"`
	IncineratorId   uint   `json:"incinerator_id"`
	IncineratorCode string `json:"incinerator_code"`
	InstrumentId    uint   `json:"instrument_id"`
	InstrumentName  string `json:"instrument_name"`
	InstrumentCode  string `json:"instrument_code"`
}

func GetInstrument(w http.ResponseWriter, r *http.Request) {
	var instModel []InstrumentModel

	destinationCode := r.URL.Query().Get("destination_code")
	incineratorId := r.URL.Query().Get("incinerator_id")

	tx := db.GetInstance().
		Raw(`
			select *
			from instruments i
			left join incinerators i2 on i2.incinerator_id = i.incinerator_id
			where
				i.incinerator_id = ?
			and exists (
				select 1 from destinations d 
				where d.destination_code = ?
			)
		`, incineratorId, destinationCode).
		Scan(&instModel)

	if tx.Error != nil {
		network.ResponseJSON(w, true, http.StatusBadRequest, instModel)
		return
	}
	network.ResponseJSON(w, false, http.StatusOK, instModel)
}

type InstrumentReq struct {
	IncineratorId  uint   `json:"incinerator_id"`
	InstrumentName string `json:"instrument_name"`
	InstrumentCode string `json:"instrument_code"`
}

func StoreInstrument(w http.ResponseWriter, r *http.Request) {
	var instModel model.Instrument
	var instReq InstrumentReq
	err := network.ReadJSONData(r, &instReq)

	if err != nil {
		network.ResponseJSON(w, true, http.StatusBadRequest, "Invalid parameter")
		return
	}

	validateInput := []utils.InputPattern{
		{Input: instReq.InstrumentCode, RegPattern: "^[a-zA-Z0-9]+$", PatternDescription: "Instrument code is required."},
		{Input: instReq.InstrumentName, RegPattern: "^[a-zA-Z0-9]+$", PatternDescription: "Instrument name is required."},
	}

	pattern, patternErr := utils.CheckPattern(validateInput)
	if patternErr {
		network.ResponseJSON(w, true, http.StatusBadRequest, pattern)
		return
	}

	tx := db.GetInstance().Where(&model.Instrument{
		IncineratorId:  instReq.IncineratorId,
		InstrumentCode: instReq.InstrumentCode,
	}).Table("instruments").Scan(&instModel)
	fmt.Println(tx.RowsAffected)
	if tx.RowsAffected == 0 {
		// Create
		tx := db.GetInstance().Table("instruments").Create(&instReq)
		if tx.Error != nil {
			fmt.Println(tx.Error.Error())
			network.ResponseJSON(w, true, http.StatusInternalServerError, ErrDbTransFail)
			return
		}
		network.ResponseJSON(w, false, http.StatusCreated, "Entry Added")
	} else {
		tx = db.GetInstance().Model(&instModel).Updates(&instReq)
		if tx.Error != nil {
			fmt.Println(tx.Error.Error())
			network.ResponseJSON(w, true, http.StatusInternalServerError, ErrDbTransFail)
			return
		}
		network.ResponseJSON(w, false, http.StatusCreated, "Entry Updated")
	}

}

type LatestInstrumentReading struct {
	InstrumentId uint   `json:"instrument_id"`
	SensorId     uint   `json:"sensor_id"`
	Label        string `json:"label"`
	Value        string `json:"value"`
	ReadAt       string `json:"read_at"`
}

func GetInstrumentData(w http.ResponseWriter, r *http.Request) {
	// now := time.Now()
	var latestReadingModel []LatestInstrumentReading
	var path = strings.Split(r.URL.Path, "/")
	instrumentId := path[len(path)-1]

	stmt := `
	select
		s.instrument_id,
		r.sensor_id,
		s.label,
		concat(max(r.value), s.unit_of_measurement) as value,
		max(r.read_at) as read_at
	from
		readings r
	left join sensors s on r.sensor_id = s.sensor_id
	where s.instrument_id = ?
	group by
		s.instrument_id,
		r.sensor_id,
		s.label,
		s.unit_of_measurement
	order by s.label
	`

	tx := db.GetInstance().Raw(stmt, instrumentId).Scan(&latestReadingModel)

	if tx.Error != nil {
		network.ResponseJSON(w, true, http.StatusNotFound, "404 Not Found")
		return
	}

	network.ResponseJSON(w, false, http.StatusOK, latestReadingModel)
}
