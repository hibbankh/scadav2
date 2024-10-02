package controllers

import (
	"app/network"
	"errors"
	"fmt"
	db "framework/database"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm/dialects/postgres"
)

type OverviewIncineratorRes struct {
	DestinationId   uint           `json:"destination_id"`
	DestinationCode string         `json:"destination_code"`
	DestinationName string         `json:"destination_name"`
	IncineratorId   uint           `json:"incinerator_Id"`
	IncineratorCode string         `json:"incinerator_code"`
	Instruments     postgres.Jsonb `gorm:"type:json;default:'{}'"`
}

func IncineratorRoute() *chi.Mux {
	router := chi.NewMux()
	router.Get("/data", getOverviewIncinerator)          //latest data by (destination_code)site and line(incinerator_code)
	router.Get("/data/hourly", getIncineratorByTime)     //filter by site, get hourly data of date
	router.Get("/data/daily", getIncineratorByRequest)   //filter by site, get data for start_date,end_date | group by day
	router.Get("/data/weekly", getIncineratorByRequest)  //filter by site, get data for start_date,end_date | group by week number
	router.Get("/data/monthly", getIncineratorByRequest) //filter by site, get data for start_date,end_date | group by month
	router.Get("/data/yearly", getIncineratorByRequest)  //filter by site, get data for start_date,end_date | group by year
	return router
}

type IncineratorOverviewRes struct {
	IncineratorId     uint   `json:"incinerator_id"`
	InstrumentCode    string `json:"InstrumentCode"`
	InstrumentId      string `json:"instrument_id"`
	InstrumentName    string `json:"instrument_name"`
	SensorId          string `json:"sensor_id"`
	Label             string `json:"label"`
	UnitOfMeasurement string `json:"unit_of_measurement"`
	Measure           string `json:"measure"`
	Value             string `json:"value"`
	ReadAt            string `json:"read_at"`
	CreatedAt         string `json:"created_at"`
}

// Get latest incinerator sensor readings
func getOverviewIncinerator(w http.ResponseWriter, r *http.Request) {
	incineratorId := r.URL.Query().Get("incinerator_id")
	destinationCode := r.URL.Query().Get("destination_code")

	if incineratorId == "" {
		network.ResponseJSON(w, true, http.StatusBadRequest, PrmNotComplt)
		return
	}

	now := time.Now()
	fmt.Println("okkkkkkkkkkkkkkkkkkkkkkkkkkk")
	respModel := []IncineratorOverviewRes{}
	query := `with info as (
		select incinerator_id, instrument_code , instrument_name, s.instrument_id, s.sensor_id, s."label", s.unit_of_measurement, s.measure
		from sensors s 
		left join instruments i using (instrument_id)
		left join readings r using (sensor_id)
		where incinerator_id = ?
		and exists ( select 1 from incinerators where incinerator_id = ? and destination_id = (select destination_id from destinations where destination_code = ?))
		group by 1,2,3,4,5,6,7,8
)
,latest_readings as (
		select
				r.sensor_id,
				r.value,
				r.read_at,
				r.created_at,
				row_number() over (partition by r.sensor_id order by r.read_at desc) as rn
		from readings r
		order by r.read_at desc
)
select * 
from info
left join latest_readings lr on info.sensor_id = lr.sensor_id 
and lr.rn = 1
	`
	tx := db.GetInstance().Raw(query, incineratorId, incineratorId, destinationCode).Scan(&respModel)
	if tx.Error != nil {
		network.ResponseJSON(w, true, http.StatusBadRequest, "ERR_DSP_OVI")
		return
	}
	network.ResponseJSON(w, false, http.StatusOK, respModel)
	fmt.Printf("execution_time: %s\n", time.Since(now))
}

/*
	Daily Incinerator Reading
*/
func getIncineratorByTime(w http.ResponseWriter, r *http.Request) {
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
	cteDailyReadingStmt := getCteReadingByHours()
	cteInstrumentStmt := getCteInstrumentV4()
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

	end := start + " 23:59:59"
	tx := db.GetInstance().Raw(formattedQuery, start, end, start, incineratorId, instrumentId, incineratorId, destinationCode).Scan(&respModel)
	if tx.Error != nil {
		network.ResponseJSON(w, true, http.StatusBadRequest, "ERR_DSP_OVI")
		return
	}
	network.ResponseJSON(w, false, http.StatusOK, respModel)
	fmt.Printf("execution_time: %s\n", time.Since(now))
}

func processDateParams(r *http.Request, now time.Time) (string, string, error) {
	path := strings.Split(r.URL.Path, "/")
	pathType := strings.ToLower(path[len(path)-1])
	if pathType != "daily" && pathType != "weekly" && pathType != "monthly" && pathType != "yearly" {
		return "0001-01-01", "0001-01-01", errors.New("only support weekly and monthly path")
	}

	queryStart := r.URL.Query().Get("from")
	queryEnd := r.URL.Query().Get("to")

	// start := utils_time.GetStartDate(now, queryStart, pathType)
	// end := utils_time.GetEndDate(now, start, queryEnd, pathType)

	return queryStart, queryEnd, nil
}

/*
	Weekly Incinerator Reading
*/
func getIncineratorByDay(w http.ResponseWriter, r *http.Request) {
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
	cteDailyReadingStmt := getCteReadingByDays()
	cteInstrumentStmt := getCteInstrumentV4()
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

	tx := db.GetInstance().Raw(formattedQuery, start, end, start, end,
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
func getIncineratorByMonth(w http.ResponseWriter, r *http.Request) {
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

	respModel := []OverviewIncineratorRes{}
	cteDailyReadingStmt := getCteReadingByMonthsOfYear()
	cteInstrumentStmt := getCteInstrumentV4()
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

	start := fmt.Sprintf("%s-01-01", year)
	end := fmt.Sprintf("%s-12-31", year)
	tx := db.GetInstance().Raw(formattedQuery, start, end, year,
		incineratorId, instrumentId, incineratorId,
		destinationCode).Scan(&respModel)
	if tx.Error != nil {
		network.ResponseJSON(w, true, http.StatusBadRequest, "ERR_DSP_OVI")
		return
	}
	network.ResponseJSON(w, false, http.StatusOK, respModel)
	fmt.Printf("execution_time: %s\n", time.Since(now))
}

type OverviewRequestParams struct {
	View struct {
		Daily   bool
		Weekly  bool
		Monthly bool
		Yearly  bool
	}
	Calendar struct {
		Start string
		End   string
	}
	Site string
}

func getOverviewQuery() string {
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
	derivedReadingStmt := getDerivedReadings()
	cteReadingStmt := getCteReadings()
	cteInstrumentStmt := getCteInstrument()
	cteIncineratorStmt := getCteIncinerator()

	rawQuery := `
	with %s, %s, %s, %s
	select
		d.destination_id,
		d.destination_code,
		d.destination_name,
		cte_i.incinerator_id,
		cte_i.incinerator_code,
		cte_i.line,
		cte_i.instruments
	from destinations d
	left join cte_incinerators cte_i on d.destination_id = cte_i.destination_id;
	`
	formattedQuery := fmt.Sprintf(rawQuery, derivedReadingStmt, cteReadingStmt, cteInstrumentStmt, cteIncineratorStmt)
	return formattedQuery
}

type (
	GenericIncinerator struct {
		DestinationCode string
		IncineratorCode string
	}
	GenericInstrument struct {
		Name string
		Code string
	}
	GenericSensor struct {
		Label             string
		UnitOfMeasurement string
		Measure           string
	}
	GenericReading struct {
		Value  string
		ReadAt time.Time
	}
	Incinerator struct {
		GenericIncinerator
	}
)

// func GenerateRandomIncineratorData() {

// 	instrument := make(map[string]interface{})
// 	&instrument{

// 	}
// }

func getIncineratorByRequest(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	path := strings.Split(r.URL.Path, "/")
	pathType := strings.ToLower(path[len(path)-1])
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
	var cteDailyReadingStmt string
	if pathType == "daily" {
		cteDailyReadingStmt = getCteReadingByDay()
	} else if pathType == "weekly" {
		cteDailyReadingStmt = getCteReadingByWeek()
	} else if pathType == "monthly" {
		cteDailyReadingStmt = getCteReadingByMonth()
	} else if pathType == "yearly" {
		cteDailyReadingStmt = getCteReadingByYear()
	}
	cteInstrumentStmt := getCteInstrumentV4()
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

	tx := db.GetInstance().Raw(formattedQuery, start, end, start, end,
		incineratorId, instrumentId, incineratorId,
		destinationCode).Scan(&respModel)
	if tx.Error != nil {
		network.ResponseJSON(w, true, http.StatusBadRequest, "ERR_DSP_OVI")
		return
	}
	network.ResponseJSON(w, false, http.StatusOK, respModel)
	fmt.Printf("execution_time: %s\n", time.Since(now))
}
