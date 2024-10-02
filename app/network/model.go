package network

import (
	"time"

	"github.com/lib/pq"
)

type (
	CustomerLogin struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	RegisterUser struct {
		CustomerName   string `json:"customername" validate:"required"`
		CustomerPasswd string `json:"password" validate:"required"`
		CustomerAddr   string `json:"address" validate:"required" `
		CustomerNo     string `json:"phoneno" validate:"required" `
		CustomerEmail  string `json:"email" validate:"required"`
		CustomerPic    string `json:"pic" validate:"required"`
	}

	PikapPointsLoc struct {
		PickupLat  float64 `json:"lat" `
		PickupLong float64 `json:"lon" `
	}

	UpdateBin struct {
		BinName     string `json:"binname"`
		BinCapacity string `json:"bincapacity"`
	}

	UpdateStat struct {
		BinStatusName string `json:"statusname"`
	}

	UpdateCustomer struct {
		CustName   string `json:"customername"`
		CustPasswd string `json:"customerpassword"`
		CustAddr   string `json:"customeraddress"`
		CustNo     string `json:"customerphonenumber"`
		CustEmail  string `json:"customeremail"`
		CustPic    string `json:"customerpic"`
	}

	SvReport struct {
		// CustomerID  string    `json:"customerid" validate:"required"`
		PickupID    string `json:"pickupid" validate:"required"`
		EmpID       string `json:"empid" validate:"required"`
		ImageData   string `json:"imagedata" validate:"required"`
		PickupState int    `json:"status" validate:"required"`
		// BinStatus   []BinInfo `json:"binstat"  `
	}

	LorryReport struct {
		PickupID  string `json:"pickupid" validate:"required"`
		ImageData string `json:"imagedata" validate:"required"`
		EmpID     string `json:"empid"`
	}

	BinInfo struct {
		BinID  string `json:"binid" validate:"required"`
		Status string `json:"status" validate:"required"`
	}

	UpdateRegion struct {
		RegionName string `json:"regionname"`
	}

	UpdateRoute struct {
		RouteName string `json:"routename"`
		RegionID  string `json:"region"`
	}
	AssingBin struct {
		BinID   string `json:"binid"`
		BinName string `json:"binname"`
	}

	CustomerOverview struct {
		CustName    string `json:"custname"`
		CustID      string `json:"custid"`
		TotalPickup string `json:"ttlcollect"`
		Collect     string `json:"collect"`
		NotCollect  string `json:"notcollect"`
	}

	RegionOverview struct {
		CustID      string `json:"custid"`
		RegionName  string `json:"regionname"`
		RegionID    string `json:"regionid"`
		TotalPickup string `json:"ttlcollect"`
		Collect     string `json:"collect"`
		Nollect     string `json:"notcollect"`
	}

	RouteOverview struct {
		RouteName   string `json:"routename"`
		RouteID     string `json:"routeid"`
		TotalPickup string `json:"ttlcollect"`
		Collect     string `json:"collect"`
		Nollect     string `json:"notcollect"`
	}

	OverallCollection struct {
		Timestamp time.Time `json:"timestamp"`
		Collect   int       `json:"collect"`
	}

	TotalPickupGraph struct {
		TotalPickup int `json:"totalpickup"`
	}

	OverallCollectionList struct {
		CustName       string `json:"custname"`
		TotalCollected int    `json:"totalcollected"`
		Collected      int    `json:"collected"`
	}

	CutomerLog struct {
		PickupName string `json:"pickupname"`
		Date       string `json:"date"`
		Time       string `json:"time"`
		CreatedAt  string `json:"createdat"`
	}

	MapStatus struct {
		PickupID     string         `json:"pickupid"`
		CustName     string         `json:"custname"`
		CustID       string         `json:"custid"`
		RegionName   string         `json:"regionname"`
		RegionID     string         `json:"regionid"`
		RouteName    string         `json:"routename"`
		RouteID      string         `json:"routeid"`
		PickupName   string         `json:"pickupname"`
		PickupStatus string         `json:"status"`
		PickupLat    string         `json:"lat"`
		PickupLong   string         `json:"long"`
		CreatedAt    pq.StringArray `gorm:"type:string[]"`
		ColImgPath   pq.StringArray `gorm:"type:string[]"`
		DriverName   pq.StringArray `gorm:"type:string[]"`
		Complaint    pq.StringArray `gorm:"type:string[]"`
	}

	CustomerRating struct {
		CustName string `json:"pickupname"`
		Star1    int32  `json:"star1"`
		Star2    int32  `json:"star2"`
		Star3    int32  `json:"star3"`
		Star4    int32  `json:"star4"`
		Star5    int32  `json:"star5"`
	}

	RegPPLatLon struct {
		Id         string  `json:"pickupid" `
		PickupLat  float64 `json:"lat" `
		PickupLong float64 `json:"lon" `
	}

	BinReg struct {
		Bintypeid string `json:"bintypeid"`
		Quantity  int    `json:"quantity"`
	}

	PickupPointsDetails struct {
		PickupID   string        `json:"pickupid"`
		RouteId    string        `json:"routeid"`
		CustId     string        `json:"custid"`
		RegionId   string        `json:"regionid"`
		PickupName string        `json:"pickupname"`
		PickupLat  float64       `json:"pickuplat"`
		PickupLong float64       `json:"pickuplong"`
		BinList    []BinReg      `json:"binlist"`
		PickupFreq pq.Int64Array `gorm:"type:int[]"`
		Radius     string        `json:"radius"`
		PPStatId   string        `json:"ppstatid"`
	}

	TestYe struct {
		Data []PickupPointsDetails `json:"data"`
	}
)
