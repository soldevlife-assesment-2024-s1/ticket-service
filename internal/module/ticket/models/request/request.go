package request

type Pagination struct {
	Page int `json:"page" form:"page" required:"true" validate:"required,numeric"`
	Size int `json:"size" form:"size" required:"true" validate:"required,numeric"`
}

type InquiryTicketAmount struct {
	TicketID    int64 `json:"ticket_id" form:"ticket_id" required:"true" validate:"required"`
	TotalTicket int   `json:"total_ticket" form:"total_ticket" required:"true" validate:"required,numeric"`
}

type CheckStockTicket struct {
	TicketDetailID string `form:"ticket_detail_id"`
}

type StockTicket struct {
	TicketDetailID int64 `json:"ticket_detail_id" form:"ticket_detail_id" validate:"required"`
	TotalTickets   int64 `json:"total_tickets" form:"total_tickets" validate:"required"`
}

type PoisonedQueue struct {
	TopicTarget string      `json:"topic_target" validate:"required"`
	ErrorMsg    string      `json:"error_msg" validate:"required"`
	Payload     interface{} `json:"payload" validate:"required"`
}

type TicketSoldOut struct {
	VenueName string `json:"venue_name" validate:"required"`
	IsSoldOut bool   `json:"is_sold_out" validate:"required"`
}

type OnlineTicketRules struct {
	IsTicketSoldOut      bool  `json:"is_ticket_sold_out" validate:"required"`
	IsTicketFirstSoldOut bool  `json:"is_ticket_first_sold_out" validate:"required"`
	TotalSeat            int64 `json:"total_seat" validate:"required"`
}
