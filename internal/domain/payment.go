package domain

import (
	"time"

	"gorm.io/gorm"
)

// Payment records money a member handed over for a session.
// It is deliberately not tied to individual time slots.
type Payment struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	MemberID  uint    `json:"member_id"`
	SessionID uint    `json:"session_id"`
	Amount    float64 `json:"amount"`

	Member  *Member  `json:"member,omitempty"`
	Session *Session `json:"-"`
}
