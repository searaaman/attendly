package domain

import "time"

type Session struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	StartTime time.Time  `json:"start_time"`
	EndTime   time.Time  `json:"end_time"`
	UnitRate  float64    `json:"unit_rate"`
	TimeSlots []TimeSlot `json:"time_slots"`
}
