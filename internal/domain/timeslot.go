package domain

import "time"

type TimeSlot struct {
	ID uint `gorm:"primaryKey" json:"id"`

	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	UnitsBooked int       `json:"units_booked"`

	SessionID uint     `json:"session_id"`
	Members   []Member `gorm:"many2many:member_time_slots;" json:"members"`
}
