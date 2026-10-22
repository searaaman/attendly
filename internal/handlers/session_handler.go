package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/searaaman/attendly/internal/config"
	"github.com/searaaman/attendly/internal/domain"
	"github.com/searaaman/attendly/internal/services"
	"gorm.io/gorm"
)

// withSlots preloads a session's slots in time order, and the members in each slot.
func withSlots(db *gorm.DB) *gorm.DB {
	return db.
		Preload("TimeSlots", func(db *gorm.DB) *gorm.DB { return db.Order("start_time, id") }).
		Preload("TimeSlots.Members", func(db *gorm.DB) *gorm.DB { return db.Order("name") })
}

func loadSession(ctx *gin.Context, sessionID uint) (domain.Session, bool) {
	var session domain.Session
	err := withSlots(config.DB).First(&session, sessionID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "session not found",
		})
		return session, false
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to load session",
		})
		return session, false
	}
	return session, true
}

func CreateSession(ctx *gin.Context) {
	var request domain.CreateSessionRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	session, err := services.NewSession(request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := config.DB.Create(&session).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create session",
		})
		return
	}
	ctx.JSON(http.StatusCreated, session)
}

func ListSessions(ctx *gin.Context) {
	sessions := []domain.Session{}
	err := withSlots(config.DB).Order("start_time DESC").Find(&sessions).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to load sessions",
		})
		return
	}
	ctx.JSON(http.StatusOK, sessions)
}

func GetSession(ctx *gin.Context) {
	sessionID, ok := parseID(ctx, "id")
	if !ok {
		return
	}
	session, ok := loadSession(ctx, sessionID)
	if !ok {
		return
	}
	ctx.JSON(http.StatusOK, session)
}

func DeleteSession(ctx *gin.Context) {
	sessionID, ok := parseID(ctx, "id")
	if !ok {
		return
	}
	if _, ok := loadSession(ctx, sessionID); !ok {
		return
	}

	err := services.DeleteSession(config.DB, sessionID)
	if errors.Is(err, services.ErrSessionHasPayments) {
		ctx.JSON(http.StatusConflict, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete session",
		})
		return
	}
	ctx.Status(http.StatusNoContent)
}

func CreateTimeSlot(ctx *gin.Context) {
	sessionID, ok := parseID(ctx, "id")
	if !ok {
		return
	}

	var request domain.CreateTimeSlotRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	session, ok := loadSession(ctx, sessionID)
	if !ok {
		return
	}
	if err := services.ValidateTimeSlot(session, request.StartTime, request.EndTime); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	timeSlot := domain.TimeSlot{
		SessionID:   sessionID,
		StartTime:   request.StartTime,
		EndTime:     request.EndTime,
		UnitsBooked: request.UnitsBooked,
		Members:     []domain.Member{},
	}
	if err := config.DB.Create(&timeSlot).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create timeslot",
		})
		return
	}
	ctx.JSON(http.StatusCreated, timeSlot)
}

func UpdateTimeSlot(ctx *gin.Context) {
	timeSlotID, ok := parseID(ctx, "id")
	if !ok {
		return
	}

	var request domain.UpdateTimeSlotRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var timeSlot domain.TimeSlot
	if err := config.DB.First(&timeSlot, timeSlotID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "timeslot not found",
		})
		return
	}

	err := config.DB.Model(&timeSlot).Update("units_booked", request.UnitsBooked).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update timeslot",
		})
		return
	}
	ctx.JSON(http.StatusOK, timeSlot)
}

func DeleteTimeSlot(ctx *gin.Context) {
	timeSlotID, ok := parseID(ctx, "id")
	if !ok {
		return
	}

	var timeSlot domain.TimeSlot
	if err := config.DB.First(&timeSlot, timeSlotID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "timeslot not found",
		})
		return
	}

	if err := services.DeleteTimeSlot(config.DB, timeSlotID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete timeslot",
		})
		return
	}
	ctx.Status(http.StatusNoContent)
}

func GetSessionBilling(ctx *gin.Context) {
	sessionID, ok := parseID(ctx, "id")
	if !ok {
		return
	}
	session, ok := loadSession(ctx, sessionID)
	if !ok {
		return
	}
	ctx.JSON(http.StatusOK, services.CalculateSessionBills(session))
}
