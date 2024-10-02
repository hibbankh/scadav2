package model

import (
	"fmt"

	"gorm.io/gorm"
)

func CreateDestination(db *gorm.DB, dCode string) uint {
	dest := Destination{
		DestinationCode: dCode,
	}
	db.Create(&dest)
	return dest.DestinationId
}

func GetDestinationId(db *gorm.DB, dCode string) uint {
	var destination Destination

	db.Where("destination_code = ?", dCode).Find(&destination)
	if destination.DestinationId == 0 {
		fmt.Println("Create destination")
		return CreateDestination(db, dCode)
	}
	return destination.DestinationId
}
