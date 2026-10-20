package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/searaaman/attendly/internal/config"
	"github.com/searaaman/attendly/internal/domain"
	"github.com/searaaman/attendly/internal/services"
)

func CreatePayment(ctx *gin.Context) {
	var request domain.CreatePaymentRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var member domain.Member
	if err := config.DB.First(&member, request.MemberID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "member not found",
		})
		return
	}

	session, ok := loadSession(ctx, request.SessionID)
	if !ok {
		return
	}

	attendedSession := false
	for _, bill := range services.CalculateSessionBills(session) {
		if bill.MemberID == request.MemberID {
			attendedSession = true
			break
		}
	}
	if !attendedSession {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "member is not part of this session",
		})
		return
	}

	payment := domain.Payment{
		MemberID:  request.MemberID,
		SessionID: request.SessionID,
		Amount:    request.Amount,
	}
	if err := config.DB.Create(&payment).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to record payment",
		})
		return
	}
	payment.Member = &member
	ctx.JSON(http.StatusCreated, payment)
}

// ListPayments returns payments, newest first, optionally filtered by
// ?member_id= and/or ?session_id=.
func ListPayments(ctx *gin.Context) {
	query := config.DB.Preload("Member").Order("created_at DESC")
	if memberID := ctx.Query("member_id"); memberID != "" {
		query = query.Where("member_id = ?", memberID)
	}
	if sessionID := ctx.Query("session_id"); sessionID != "" {
		query = query.Where("session_id = ?", sessionID)
	}

	payments := []domain.Payment{}
	if err := query.Find(&payments).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to load payments",
		})
		return
	}
	ctx.JSON(http.StatusOK, payments)
}

func DeletePayment(ctx *gin.Context) {
	paymentID, ok := parseID(ctx, "id")
	if !ok {
		return
	}

	result := config.DB.Delete(&domain.Payment{}, paymentID)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete payment",
		})
		return
	}
	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "payment not found",
		})
		return
	}
	ctx.Status(http.StatusNoContent)
}
