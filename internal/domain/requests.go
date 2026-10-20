package domain

import "time"

type CreateMemberRequest struct {
	Name  string `json:"name" binding:"required"`
	Phone string `json:"phone"`
}

type CreateSessionRequest struct {
	StartTime time.Time `json:"start_time" binding:"required"`
	EndTime   time.Time `json:"end_time" binding:"required"`
	UnitRate  float64   `json:"unit_rate" binding:"gte=0"`

	// When HourlySlots is true, the session is created with one slot per hour
	// between StartTime and EndTime, each booked with DefaultUnits units.
	HourlySlots  bool `json:"hourly_slots"`
	DefaultUnits int  `json:"default_units" binding:"gte=0"`
}

type CreateTimeSlotRequest struct {
	StartTime   time.Time `json:"start_time" binding:"required"`
	EndTime     time.Time `json:"end_time" binding:"required"`
	UnitsBooked int       `json:"units_booked" binding:"gte=0"`
}

type UpdateTimeSlotRequest struct {
	UnitsBooked int `json:"units_booked" binding:"gte=0"`
}

type AssignMemberRequest struct {
	MemberID uint `json:"member_id" binding:"required"`
}

type CreatePaymentRequest struct {
	MemberID  uint    `json:"member_id" binding:"required"`
	SessionID uint    `json:"session_id" binding:"required"`
	Amount    float64 `json:"amount" binding:"gt=0"`
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
