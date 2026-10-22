package services

import (
	"errors"
	"time"

	"github.com/searaaman/attendly/internal/domain"
	"gorm.io/gorm"
)

var (
	ErrInvalidTimeRange   = errors.New("end time must be after start time")
	ErrSlotOutsideSession = errors.New("time slot must fall within the session")
	ErrSessionHasPayments = errors.New("session has payments recorded; delete those first")
)

// BuildHourlySlots cuts [start, end) into one-hour slots. If the range is not a
// whole number of hours, the last slot is shorter.
func BuildHourlySlots(start, end time.Time, units int) []domain.TimeSlot {
	var slots []domain.TimeSlot
	for slotStart := start; slotStart.Before(end); slotStart = slotStart.Add(time.Hour) {
		slotEnd := slotStart.Add(time.Hour)
		if slotEnd.After(end) {
			slotEnd = end
		}
		slots = append(slots, domain.TimeSlot{
			StartTime:   slotStart,
			EndTime:     slotEnd,
			UnitsBooked: units,
			Members:     []domain.Member{},
		})
	}
	return slots
}

func NewSession(request domain.CreateSessionRequest) (domain.Session, error) {
	if !request.EndTime.After(request.StartTime) {
		return domain.Session{}, ErrInvalidTimeRange
	}
	session := domain.Session{
		StartTime: request.StartTime,
		EndTime:   request.EndTime,
		UnitRate:  request.UnitRate,
	}
	if request.HourlySlots {
		session.TimeSlots = BuildHourlySlots(request.StartTime, request.EndTime, request.DefaultUnits)
	}
	return session, nil
}

func ValidateTimeSlot(session domain.Session, start, end time.Time) error {
	if !end.After(start) {
		return ErrInvalidTimeRange
	}
	if start.Before(session.StartTime) || end.After(session.EndTime) {
		return ErrSlotOutsideSession
	}
	return nil
}

// DeleteSession removes a session, its slots and their attendance records.
// Sessions with payments are refused so money records are never silently lost.
func DeleteSession(db *gorm.DB, sessionID uint) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var paymentCount int64
		if err := tx.Model(&domain.Payment{}).Where("session_id = ?", sessionID).Count(&paymentCount).Error; err != nil {
			return err
		}
		if paymentCount > 0 {
			return ErrSessionHasPayments
		}
		// Soft-deleted payments still hold a foreign key to the session, so purge them.
		if err := tx.Unscoped().Where("session_id = ?", sessionID).Delete(&domain.Payment{}).Error; err != nil {
			return err
		}

		var slotIDs []uint
		if err := tx.Model(&domain.TimeSlot{}).Where("session_id = ?", sessionID).Pluck("id", &slotIDs).Error; err != nil {
			return err
		}
		if len(slotIDs) > 0 {
			if err := tx.Exec("DELETE FROM member_time_slots WHERE time_slot_id IN ?", slotIDs).Error; err != nil {
				return err
			}
			if err := tx.Where("id IN ?", slotIDs).Delete(&domain.TimeSlot{}).Error; err != nil {
				return err
			}
		}
		return tx.Delete(&domain.Session{}, sessionID).Error
	})
}

// DeleteTimeSlot removes a slot and its attendance records.
func DeleteTimeSlot(db *gorm.DB, slotID uint) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("DELETE FROM member_time_slots WHERE time_slot_id = ?", slotID).Error; err != nil {
			return err
		}
		return tx.Delete(&domain.TimeSlot{}, slotID).Error
	})
}
