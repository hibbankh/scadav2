package mqtt

import (
	"app/model"
	"fmt"

	"gorm.io/gorm"
)

func storeAlarmLog(db *gorm.DB, logReq model.AlarmLog) bool {
	tx := db.Create(&logReq)

	if tx.Error != nil {
		fmt.Printf("ERR_STORE_ALARM_LOG\n \r\ndata:%v \n db_err: %v", logReq, tx.Error)
		return false
	}
	return true
}
