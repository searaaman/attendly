package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/searaaman/attendly/internal/config"
	"github.com/searaaman/attendly/internal/domain"
	"github.com/searaaman/attendly/internal/services"
)

func CreateMember(ctx *gin.Context) {
	var request domain.CreateMemberRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	member := domain.Member{
		Name:  request.Name,
		Phone: request.Phone,
	}
	if err := config.DB.Create(&member).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create member",
		})
		return
	}
	ctx.JSON(http.StatusCreated, member)
}

func GetMembers(ctx *gin.Context) {
	members := []domain.Member{}
	if err := config.DB.Order("name").Find(&members).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to load members",
		})
		return
	}
	ctx.JSON(http.StatusOK, members)
}

func AssignMemberToTimeSlot(ctx *gin.Context) {
	timeSlotID, ok := parseID(ctx, "id")
	if !ok {
		return
	}

	var request domain.AssignMemberRequest
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

	var member domain.Member
	if err := config.DB.First(&member, request.MemberID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "member not found",
		})
		return
	}

	// Appending an already-assigned member is a no-op, so this is safe to retry.
	if err := config.DB.Model(&timeSlot).Association("Members").Append(&member); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to assign member",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "member assigned successfully",
	})
}

func RemoveMemberFromTimeSlot(ctx *gin.Context) {
	timeSlotID, ok := parseID(ctx, "id")
	if !ok {
		return
	}
	memberID, ok := parseID(ctx, "memberId")
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

	member := domain.Member{ID: memberID}
	if err := config.DB.Model(&timeSlot).Association("Members").Delete(&member); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to remove member",
		})
		return
	}
	ctx.Status(http.StatusNoContent)
}

// loadLedgerData fetches everything the ledger calculation needs.
func loadLedgerData() ([]domain.Session, []domain.Payment, error) {
	var sessions []domain.Session
	if err := withSlots(config.DB).Find(&sessions).Error; err != nil {
		return nil, nil, err
	}
	var payments []domain.Payment
	if err := config.DB.Order("created_at DESC").Find(&payments).Error; err != nil {
		return nil, nil, err
	}
	return sessions, payments, nil
}

func GetLedger(ctx *gin.Context) {
	var members []domain.Member
	if err := config.DB.Find(&members).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to load members",
		})
		return
	}
	sessions, payments, err := loadLedgerData()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to load ledger",
		})
		return
	}
	ctx.JSON(http.StatusOK, services.BuildLedger(members, sessions, payments))
}

func GetMemberLedger(ctx *gin.Context) {
	memberID, ok := parseID(ctx, "id")
	if !ok {
		return
	}

	var member domain.Member
	if err := config.DB.First(&member, memberID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "member not found",
		})
		return
	}

	sessions, payments, err := loadLedgerData()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to load ledger",
		})
		return
	}
	ctx.JSON(http.StatusOK, services.BuildMemberLedger(member, sessions, payments))
}
