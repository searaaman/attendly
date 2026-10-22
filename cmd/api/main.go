package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/searaaman/attendly/internal/config"
	"github.com/searaaman/attendly/internal/docs"
	"github.com/searaaman/attendly/internal/handlers"
	"github.com/searaaman/attendly/internal/middleware"
)

func main() {
	if err := config.Load(); err != nil {
		log.Fatal("Invalid configuration: ", err)
	}
	config.ConnectDatabase()

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     config.Cfg.CORSOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.GET("/health", handlers.HealthHandler)
	r.GET("/docs", docs.UIHandler)
	r.GET("/docs/openapi.yaml", docs.SpecHandler)
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())

	auth.GET("/sessions", handlers.ListSessions)
	auth.POST("/sessions", handlers.CreateSession)
	auth.GET("/sessions/:id", handlers.GetSession)
	auth.DELETE("/sessions/:id", handlers.DeleteSession)
	auth.POST("/sessions/:id/timeslots", handlers.CreateTimeSlot)
	auth.GET("/sessions/:id/billing", handlers.GetSessionBilling)

	auth.PUT("/timeslots/:id", handlers.UpdateTimeSlot)
	auth.DELETE("/timeslots/:id", handlers.DeleteTimeSlot)
	auth.POST("/timeslots/:id/members", handlers.AssignMemberToTimeSlot)
	auth.DELETE("/timeslots/:id/members/:memberId", handlers.RemoveMemberFromTimeSlot)

	auth.GET("/members", handlers.GetMembers)
	auth.POST("/members", handlers.CreateMember)
	auth.GET("/members/:id/ledger", handlers.GetMemberLedger)
	auth.GET("/ledger", handlers.GetLedger)

	auth.GET("/payments", handlers.ListPayments)
	auth.POST("/payments", handlers.CreatePayment)
	auth.DELETE("/payments/:id", handlers.DeletePayment)

	if err := r.Run(":" + config.Cfg.ServerPort); err != nil {
		log.Fatal(err)
	}
}
