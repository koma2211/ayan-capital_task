package entities

import "time"

type Event struct {
	OrderType  string    `json:"orderType" binding:"required"`
	SessionID  string    `json:"sessionId" binding:"required"`
	Card       string    `json:"card" binding:"required"`
	WebSiteURL string    `json:"websiteUrl" binding:"required"`
	Date       time.Time `json:"eventDate" binding:"required"`
}
