package controllers

import (
	"app/network"
	"fmt"
	db "framework/database"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

func IncineratorRouteV2() *chi.Mux {
	router := chi.NewMux()
	// router.Get("/data", getOverviewIncinerator)			//
	router.Get("/data/daily", getIncineratorByTimeV2)  //filter by site, get hourly data of date
	router.Get("/data/weekly", getIncineratorByDayV2)  //filter by site, get data for start_date,end_date of week
	router.Get("/data/monthly", getIncineratorByDayV2) //filter by site, get data for start_date,end_date of month

	router.Get("/data/yearly", getIncineratorByMonthV2) //filter by site, get data for monthly of year
	return router
}

/*
	Daily Incinerator Reading
*/
func getIncineratorByTimeV2(w http.ResponseWriter, r *http.Request) {
	now := time.Now()

	var start, destinationCode, instrumentId, incineratorId string

	if r.URL.Query().Has("from"); r.URL.Query().Get("from") != "" {
		start = r.URL.Query().Get("from")
	} else {
		start = now.Format("2006-01-02")
	}

	if r.URL.Query().Has("destination_code"); r.URL.Query().Get("destination_code") != "" {
		destinationCode = r.URL.Query().Get("destination_code")
	} else {
		destinationCode = "00"
	}

	if r.URL.Query().Has("incinerator_id"); r.URL.Query().Get("incinerator_id") != "" {
		incineratorId = r.URL.Query().Get("incinerator_id")
	} else {
		incineratorId = "00"
	}

	if r.URL.Query().Has("instrument_id"); r.URL.Query().Get("instrument_id") != "" {
		instrumentId = r.URL.Query().Get("instrument_id")
	} else {
		instrumentId = "00"
	}

	respModel := []OverviewIncineratorRes{}
	// cteDailyReadingStmt := getCteReadingByHoursV2()
	// cteInstrumentStmt := getCteInstrumentV2()
	cteDailyReadingStmt := getCteReadingByHoursV3()
	cteInstrumentStmt := getCteInstrumentV3()
	cteIncineratorStmt := getCteIncineratorV2()

	query := ` with  %s, %s, %s
	select
		d.destination_id,
		d.destination_code,
		d.destination_name,
		cte_i.incinerator_id,
		cte_i.incinerator_code,
		cte_i.line,
		cte_i.instruments
	from destinations d
	left join cte_incinerators cte_i on d.destination_id = cte_i.destination_id
	where d.destination_code = ?
	`
	formattedQuery := fmt.Sprintf(query, cteDailyReadingStmt, cteInstrumentStmt, cteIncineratorStmt)

	tx := db.GetInstance().Raw(formattedQuery, start, incineratorId, instrumentId, incineratorId, destinationCode).Scan(&respModel)
	if tx.Error != nil {
		network.ResponseJSON(w, true, http.StatusBadRequest, "ERR_DSP_OVI")
		return
	}
	network.ResponseJSON(w, false, http.StatusOK, respModel)
	fmt.Printf("execution_time: %s\n", time.Since(now))
}

/*
	Weekly Incinerator Reading
*/
func getIncineratorByDayV2(w http.ResponseWriter, r *http.Request) {
	now := time.Now()

	start, end, err := processDateParams(r, now)
	var destinationCode, instrumentId, incineratorId string

	if r.URL.Query().Has("destination_code"); r.URL.Query().Get("destination_code") != "" {
		destinationCode = r.URL.Query().Get("destination_code")
	} else {
		destinationCode = "00"
	}

	if r.URL.Query().Has("incinerator_id"); r.URL.Query().Get("incinerator_id") != "" {
		incineratorId = r.URL.Query().Get("incinerator_id")
	} else {
		incineratorId = "00"
	}

	if r.URL.Query().Has("instrument_id"); r.URL.Query().Get("instrument_id") != "" {
		instrumentId = r.URL.Query().Get("instrument_id")
	} else {
		instrumentId = "00"
	}

	if err != nil {
		fmt.Println(err)
		network.ResponseJSON(w, true, http.StatusNotFound, "404 page not found")
		return
	}

	respModel := []OverviewIncineratorRes{}
	cteDailyReadingStmt := getCteReadingByDaysV2()
	cteInstrumentStmt := getCteInstrumentV3()
	cteIncineratorStmt := getCteIncineratorV2()

	query := ` with  %s, %s, %s
	select
		d.destination_id,
		d.destination_code,
		d.destination_name,
		cte_i.incinerator_id,
		cte_i.incinerator_code,
		cte_i.line,
		cte_i.instruments
	from destinations d
	left join cte_incinerators cte_i on d.destination_id = cte_i.destination_id
	where d.destination_code = ?
	`
	formattedQuery := fmt.Sprintf(query, cteDailyReadingStmt, cteInstrumentStmt, cteIncineratorStmt)

	tx := db.GetInstance().Raw(formattedQuery,
		start,
		end,
		incineratorId, instrumentId, incineratorId,
		destinationCode).Scan(&respModel)
	if tx.Error != nil {
		network.ResponseJSON(w, true, http.StatusBadRequest, "ERR_DSP_OVI")
		return
	}
	network.ResponseJSON(w, false, http.StatusOK, respModel)
	fmt.Printf("execution_time: %s\n", time.Since(now))
}

/*
	Monthly of Year Incinerator Reading

	Query Params
		- year
		- destination(site)
*/
func getIncineratorByMonthV2(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	var year string
	var destinationCode, instrumentId, incineratorId string

	if r.URL.Query().Has("destination_code"); r.URL.Query().Get("destination_code") != "" {
		destinationCode = r.URL.Query().Get("destination_code")
	} else {
		destinationCode = "00"
	}

	if r.URL.Query().Has("incinerator_id"); r.URL.Query().Get("incinerator_id") != "" {
		incineratorId = r.URL.Query().Get("incinerator_id")
	} else {
		incineratorId = "00"
	}

	if r.URL.Query().Has("instrument_id"); r.URL.Query().Get("instrument_id") != "" {
		instrumentId = r.URL.Query().Get("instrument_id")
	} else {
		instrumentId = "00"
	}

	if r.URL.Query().Has("from"); r.URL.Query().Get("from") != "" {
		year = r.URL.Query().Get("from")
	} else {
		year = fmt.Sprint(now.Year())
	}

	if r.URL.Query().Has("destination_code"); r.URL.Query().Get("destination_code") != "" {
		destinationCode = r.URL.Query().Get("destination_code")
	} else {
		destinationCode = "00"
	}

	respModel := []OverviewIncineratorRes{}
	cteDailyReadingStmt := getCteReadingByMonthsOfYearV2()
	cteInstrumentStmt := getCteInstrumentV3()
	cteIncineratorStmt := getCteIncineratorV2()

	query := ` with  %s, %s, %s
	select
		d.destination_id,
		d.destination_code,
		d.destination_name,
		cte_i.incinerator_id,
		cte_i.incinerator_code,
		cte_i.line,
		cte_i.instruments
	from destinations d
	left join cte_incinerators cte_i on d.destination_id = cte_i.destination_id
	where d.destination_code = ?
	`
	formattedQuery := fmt.Sprintf(query, cteDailyReadingStmt, cteInstrumentStmt, cteIncineratorStmt)

	tx := db.GetInstance().Raw(formattedQuery,
		year,
		incineratorId, instrumentId, incineratorId,
		destinationCode).Scan(&respModel)
	if tx.Error != nil {
		network.ResponseJSON(w, true, http.StatusBadRequest, "ERR_DSP_OVI")
		return
	}
	network.ResponseJSON(w, false, http.StatusOK, respModel)
	fmt.Printf("execution_time: %s\n", time.Since(now))
}
