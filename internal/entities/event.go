package entities

type Event struct {
	OrderType  string `json:"orderType" binding:"required"`
	SessionID  string `json:"sessionId" binding:"required"`
	Card       string `json:"card" binding:"required"`
	Date       string `json:"eventDate" binding:"required"`
	WebSiteURL string `json:"websiteUrl" binding:"required"`
}
