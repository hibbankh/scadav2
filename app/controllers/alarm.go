package controllers

import (
	"app/model"
	"app/network"
	"fmt"
	db "framework/database"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm/dialects/postgres"
)

func AlarmRoute() *chi.Mux {
	router := chi.NewMux()
	router.Get("/daily", GetAlarmLogDaily)
	router.Get("/weekly", GetAlarmLogWeekly)
	router.Get("/monthly", GetAlarmLogMonthly)
	router.Get("/yearly", GetAlarmLogYearly)
	return router
}

/*
	start - datetime
	end - datetime

	filter by
		view:
			-	daily
			-	weekly
			-	monthly
			- 	yearly
		calendar:
			- date range
		Site

*/

type AlarmLogRes struct {
	AlarmId         uint           `json:"alarm_id"`
	DestinationId   uint           `json:"destination_id"`
	DestinationCode string         `json:"destination_code"`
	InstrumentId    uint           `json:"instrument_id"`
	TimeLct         string         `json:"time_lct"`
	Priority        string         `json:"priority"`
	State           string         `json:"state"`
	Node            string         `json:"node"`
	Group           string         `json:"group"`
	TagName         string         `json:"tag_name"`
	Description     string         `json:"description"`
	Type            string         `json:"type"`
	Limit           string         `json:"limit"`
	CurrentValue    string         `json:"current_value"`
	Operator        string         `json:"operator"`
	AlarmDuration   string         `json:"alarm_duration"`
	UnAckDuration   string         `json:"un_act_duration"`
	CreatedAt       string         `json:"created_at"`
	Instrument      postgres.Jsonb `gorm:"type:json;default:'{}'"`
}

var stmt = `
	with cte_destination as (
		select destination_id, destination_code, destination_name
		from destinations
	)
	, cte_incinerator as (
		select incinerator_id, incinerator_code, json_agg(json_build_object(
		'destination_id', cte_d.destination_id,
		'destination_code', cte_d.destination_code,
		'destination_name', cte_d.destination_name
		)) as destination
		from incinerators i 
		left join cte_destination cte_d using (destination_id)
		group by 1,2
	)
	, cte_instrument as (
		select instrument_id, instrument_code, instrument_name
		, json_agg(json_build_object(
		'incinerator_id', cte_i.incinerator_id,
		'incinerator_code', cte_i.incinerator_code,
		'destination', cte_i.destination
		)) as incinerator
		from instruments i
		left join cte_incinerator cte_i using (incinerator_id)
		group by 1,2,3
	)
	select 
		al.alarm_id, al.destination_id, al.time_lct, al.priority, al.state, al.node, al.group, al.tag_name, al.description,
		al.type, al.limit, al.current_value, al.operator, al.alarm_duration, al.un_ack_duration, al.created_at
		,json_agg(json_build_object(
			'instrument_id', cte_i.instrument_id,
			'instrument_code', cte_i.instrument_code,
			'instrument_name', cte_i.instrument_name,
			'incinerator', cte_i.incinerator
		)) as instrument
	from alarm_logs al
	left join cte_instrument cte_i using (instrument_id)
	where destination_id = ? and time_lct between ? and ?
	group by 1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16
	order by created_at desc
`

func GetAlarmLogDaily(w http.ResponseWriter, r *http.Request) {
	// respModel := []AlarmLogRes{}
	respModel := []AlarmLogRes{}
	var destinationCode, from, to string
	if r.URL.Query().Has("destination_code"); r.URL.Query().Get("destination_code") != "" {
		destinationCode = r.URL.Query().Get("destination_code")
	} else {
		network.ResponseJSON(w, true, http.StatusNotFound, "404 page not found")
		return
	}

	destId := model.GetDestinationId(db.GetInstance(), destinationCode)
	if r.URL.Query().Has("date") && r.URL.Query().Get("date") != "" {
		from = r.URL.Query().Get("date")
		parsedDate, err := time.Parse("2006-01-02", from)
		if err != nil {
			http.Error(w, "Invalid date format", http.StatusBadRequest)
			return
		}
		// Set to start of today
		to = parsedDate.Add(24*time.Hour-time.Second).AddDate(0, 0, 1).Format("2006-01-02")
	} else {
		today := time.Now().UTC().Truncate(24 * time.Hour)
		from = today.Format("2006-01-02")
		// Set to the end of today
		to = today.Add(24*time.Hour - time.Second).Format("2006-01-02")
	}
	tx := db.GetInstance().Raw(stmt, destId, from, to).Scan(&respModel)

	if tx.Error != nil {
		network.ResponseJSON(w, true, http.StatusBadRequest, ErrProcessError)
		return
	}

	network.ResponseJSON(w, false, http.StatusOK, respModel)
}

func GetAlarmLogWeekly(w http.ResponseWriter, r *http.Request) {
	respModel := []AlarmLogRes{}

	now := time.Now()
	from, to, err := processDateParams(r, now)
	var destinationCode string

	if r.URL.Query().Has("destination_code"); r.URL.Query().Get("destination_code") != "" {
		destinationCode = r.URL.Query().Get("destination_code")
	} else {
		network.ResponseJSON(w, true, http.StatusNotFound, "404 page not found")
		return
	}

	if err != nil {
		fmt.Println(err)
		network.ResponseJSON(w, true, http.StatusNotFound, "404 page not found")
		return
	}
	destId := model.GetDestinationId(db.GetInstance(), destinationCode)

	tx := db.GetInstance().Raw(stmt, destId, from, to).Scan(&respModel)

	if tx.Error != nil {
		network.ResponseJSON(w, true, http.StatusBadRequest, ErrProcessError)
		return
	}
	network.ResponseJSON(w, false, http.StatusOK, respModel)
}
func GetAlarmLogMonthly(w http.ResponseWriter, r *http.Request) {
	respModel := []AlarmLogRes{}

	now := time.Now()
	from, to, err := processDateParams(r, now)
	var destinationCode string

	if r.URL.Query().Has("destination_code"); r.URL.Query().Get("destination_code") != "" {
		destinationCode = r.URL.Query().Get("destination_code")
	} else {
		network.ResponseJSON(w, true, http.StatusNotFound, "404 page not found")
		return
	}

	if err != nil {
		fmt.Println(err)
		network.ResponseJSON(w, true, http.StatusNotFound, "404 page not found")
		return
	}
	destId := model.GetDestinationId(db.GetInstance(), destinationCode)

	tx := db.GetInstance().Raw(stmt, destId, from, to).Scan(&respModel)

	if tx.Error != nil {
		network.ResponseJSON(w, true, http.StatusBadRequest, ErrProcessError)
		return
	}
	network.ResponseJSON(w, false, http.StatusOK, respModel)
}

func GetAlarmLogYearly(w http.ResponseWriter, r *http.Request) {}
