package controllers

import "time"

var loc, _ = time.LoadLocation("Asia/Kuala_Lumpur")

//should be banyak or spaerate , idk
//think on error logging streaming
const (
	ErrInvalidDataStruct = "Unsupported data format"
	ErrProcessError      = "Something not right"
	ErrDbTransFail       = "Error during transaction"
	ErrValidateFail      = "Data validation failed"
	ProcessComplete      = "Data transction complete"
	JSONParseFail        = "Data parse failed"

	//database
	ProcDelCmplt = "Data successfully delete"
	ProcinCmplt  = "Data successfully insert"
	ProcinUpdt   = "Data successfully update"
	RecNotFound  = "Record not available"
	RecIsFound   = "Record available"
	PrmNotComplt = "Invalid parameter"

	//acl
	DntHvePmr = "You dont have permission to access this resources"
	AclAssco  = "There is accout asscited with this ACL"
	//http
	MthdNotAllw = "Method Not Allowed"

	//logout
	SuccessLogOut = "Logout success"

	//Override
	OverOK   = "Scheduler Override success"
	OverAvai = "Scheduler override is ongoing"
	OverFail = "Scheduler override failed"

	//attr
	EtyNotFnd   = "Entity not found"
	TypeDevice  = "DEVICE"
	TypeCluster = "CLUSTER"
	TypeTenant  = "TENANT"
	TypeUser    = "USER"
)
