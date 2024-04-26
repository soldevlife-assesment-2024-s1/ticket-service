package usecases_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"ticket-service/internal/module/ticket/mocks"
	"ticket-service/internal/module/ticket/models/entity"
	"ticket-service/internal/module/ticket/models/response"
	"ticket-service/internal/module/ticket/usecases"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	uc   usecases.Usecases
	repo *mocks.Repositories
	ctx  context.Context
)

func setup() {
	repo = new(mocks.Repositories)
	uc = usecases.New(repo, nil, nil)
	ctx = context.Background()
}
func TestGetTicketByRegionName(t *testing.T) {
	// Setup
	setup()
	defer repo.AssertExpectations(t)

	t.Run("success", func(t *testing.T) {
		regionName := "exampleRegion"

		// Mock repository calls
		ticket := entity.Ticket{
			ID:        1,
			Capacity:  100,
			Region:    regionName,
			EventDate: time.Now(),
			CreatedAt: time.Time{},
			UpdatedAt: sql.NullTime{},
			DeletedAt: sql.NullTime{},
		}
		ticketDetails := []entity.TicketDetail{
			{
				ID:        1,
				TicketID:  ticket.ID,
				Level:     "exampleLevel",
				Stock:     10,
				BasePrice: 100.00,
				CreatedAt: time.Time{},
				UpdatedAt: sql.NullTime{},
				DeletedAt: sql.NullTime{},
			},
		}

		repo.On("FindTicketByRegionName", ctx, regionName).Return(ticket, nil)
		repo.On("FindTicketDetailByTicketID", ctx, ticket.ID).Return(ticketDetails, nil)

		// Call the function
		resp, err := uc.GetTicketByRegionName(ctx, regionName)

		// Assertions
		require.NoError(t, err)
		require.Len(t, resp, 1)

		expectedTicket := response.Ticket{
			ID:        ticketDetails[0].ID,
			Region:    ticket.Region,
			EventDate: ticket.EventDate,
			Level:     ticketDetails[0].Level,
			Price:     ticketDetails[0].BasePrice,
			Stock:     ticketDetails[0].Stock,
		}
		require.Equal(t, expectedTicket, resp[0])
	})

	t.Run("error", func(t *testing.T) {
		regionName := "exampleRegionError"

		// Mock repository calls
		repo.On("FindTicketByRegionName", ctx, regionName).Return(entity.Ticket{}, errors.New("error"))

		// Call the function
		resp, err := uc.GetTicketByRegionName(ctx, regionName)

		// Assertions
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestDecrementTicketStock(t *testing.T) {
	// Setup
	setup()
	defer repo.AssertExpectations(t)

	t.Run("success", func(t *testing.T) {
		ticketDetailID := int64(1)
		totalTicket := int64(5)

		// Mock repository calls
		ticketDetail := entity.TicketDetail{
			ID:        ticketDetailID,
			TicketID:  1,
			Level:     "exampleLevel",
			Stock:     10,
			BasePrice: 100.00,
			CreatedAt: time.Time{},
			UpdatedAt: sql.NullTime{},
			DeletedAt: sql.NullTime{},
		}

		ticketDetailUpdate := entity.TicketDetail{
			ID:        ticketDetailID,
			TicketID:  1,
			Level:     "exampleLevel",
			Stock:     5,
			BasePrice: 100.00,
			CreatedAt: time.Time{},
			UpdatedAt: sql.NullTime{},
			DeletedAt: sql.NullTime{},
		}

		ticketDetails := []entity.TicketDetail{
			{
				ID:        ticketDetailID,
				TicketID:  1,
				Level:     "exampleLevel",
				Stock:     10,
				BasePrice: 100.00,
				CreatedAt: time.Time{},
				UpdatedAt: sql.NullTime{},
				DeletedAt: sql.NullTime{},
			},
		}

		ticket := entity.Ticket{
			ID:        1,
			Capacity:  100,
			Region:    "exampleRegion",
			EventDate: time.Now(),
			CreatedAt: time.Time{},
			UpdatedAt: sql.NullTime{},
			DeletedAt: sql.NullTime{},
		}

		repo.On("FindTicketDetail", ctx, ticketDetailID).Return(ticketDetail, nil)
		repo.On("FindTicketDetailByTicketID", ctx, ticketDetailID).Return(ticketDetails, nil)
		repo.On("FindTicketByID", ctx, ticketDetail.TicketID).Return(ticket, nil)
		repo.On("UpdateTicketDetail", ctx, ticketDetailUpdate).Return(nil)

		// Call the function
		err := uc.DecrementTicketStock(ctx, ticketDetailID, totalTicket)

		// Assertions
		require.NoError(t, err)
	})

	t.Run("error find ticket detail", func(t *testing.T) {
		ticketDetailID := int64(1)
		totalTicket := int64(15)

		// Mock repository calls
		ticketDetail := entity.TicketDetail{
			ID:        ticketDetailID,
			TicketID:  1,
			Level:     "exampleLevel",
			Stock:     10,
			BasePrice: 100.00,
			CreatedAt: time.Time{},
			UpdatedAt: sql.NullTime{},
			DeletedAt: sql.NullTime{},
		}

		repo.On("FindTicketDetail", ctx, ticketDetailID).Return(ticketDetail, nil)

		// Call the function
		err := uc.DecrementTicketStock(ctx, ticketDetailID, totalTicket)

		// Assertions
		require.Error(t, err)
		require.EqualError(t, err, "stock not enough")
	})
}

func TestIncrementTicketStock(t *testing.T) {
	// Setup
	setup()
	defer repo.AssertExpectations(t)

	t.Run("success", func(t *testing.T) {
		ticketDetailID := int64(1)
		totalTicket := int64(5)

		// Mock repository calls
		ticketDetail := entity.TicketDetail{
			ID:        ticketDetailID,
			TicketID:  1,
			Level:     "exampleLevel",
			Stock:     10,
			BasePrice: 100.00,
			CreatedAt: time.Time{},
			UpdatedAt: sql.NullTime{},
			DeletedAt: sql.NullTime{},
		}

		ticketDetailUpdate := entity.TicketDetail{
			ID:        ticketDetailID,
			TicketID:  1,
			Level:     "exampleLevel",
			Stock:     15,
			BasePrice: 100.00,
			CreatedAt: time.Time{},
			UpdatedAt: sql.NullTime{},
			DeletedAt: sql.NullTime{},
		}

		repo.On("FindTicketDetail", ctx, ticketDetailID).Return(ticketDetail, nil)
		repo.On("UpdateTicketDetail", ctx, ticketDetailUpdate).Return(nil)

		// Call the function
		err := uc.IncrementTicketStock(ctx, ticketDetailID, totalTicket)

		// Assertions
		require.NoError(t, err)
	})
}

func TestCheckStockTicket(t *testing.T) {
	// Setup
	setup()
	defer repo.AssertExpectations(t)

	t.Run("success", func(t *testing.T) {
		ticketDetailID := 1

		// Mock repository calls
		ticketDetail := entity.TicketDetail{
			ID:        int64(ticketDetailID),
			TicketID:  1,
			Level:     "exampleLevel",
			Stock:     10,
			BasePrice: 100.00,
			CreatedAt: time.Time{},
			UpdatedAt: sql.NullTime{},
			DeletedAt: sql.NullTime{},
		}

		repo.On("FindTicketDetail", ctx, int64(ticketDetailID)).Return(ticketDetail, nil)

		// Call the function
		result, err := uc.CheckStockTicket(ctx, ticketDetailID)

		// Assertions
		require.NoError(t, err)
		require.Equal(t, response.StockTicket{Stock: 10}, result)
	})
}

func TestInquiryTicketAmount(t *testing.T) {
	// Setup
	setup()
	defer repo.AssertExpectations(t)

	t.Run("success", func(t *testing.T) {
		ticketID := int64(1)
		totalTicket := 5

		// Mock repository calls
		ticket := entity.Ticket{
			ID:        ticketID,
			Capacity:  100,
			Region:    "exampleRegion",
			EventDate: time.Now(),
			CreatedAt: time.Time{},
			UpdatedAt: sql.NullTime{},
			DeletedAt: sql.NullTime{},
		}

		ticketDetails := []entity.TicketDetail{
			{
				ID:        1,
				TicketID:  ticket.ID,
				Level:     "exampleLevel",
				Stock:     10,
				BasePrice: 100.00,
				CreatedAt: time.Time{},
				UpdatedAt: sql.NullTime{},
				DeletedAt: sql.NullTime{},
			},
			{
				ID:        2,
				TicketID:  ticket.ID,
				Level:     "exampleLevel",
				Stock:     10,
				BasePrice: 100.00,
				CreatedAt: time.Time{},
				UpdatedAt: sql.NullTime{},
				DeletedAt: sql.NullTime{},
			},
			{
				ID:        3,
				TicketID:  ticket.ID,
				Level:     "exampleLevel",
				Stock:     10,
				BasePrice: 100.00,
				CreatedAt: time.Time{},
				UpdatedAt: sql.NullTime{},
				DeletedAt: sql.NullTime{},
			},
		}

		repo.On("FindTicketDetail", ctx, ticketDetails[0].ID).Return(ticketDetails[0], nil)

		// Call the function
		result, err := uc.InquiryTicketAmount(ctx, ticketID, totalTicket)

		// Assertions
		require.NoError(t, err)
		require.Equal(t, response.InquiryTicketAmount{
			TotalTicket: 5,
			TotalAmount: 500.00,
		}, result)
	})

	t.Run("error find ticket", func(t *testing.T) {
		ticketID := int64(2)
		totalTicket := 5

		// Mock repository calls
		repo.On("FindTicketDetail", ctx, ticketID).Return(entity.TicketDetail{}, errors.New("error"))

		// Call the function
		result, err := uc.InquiryTicketAmount(ctx, ticketID, totalTicket)

		// Assertions
		require.Error(t, err)
		assert.Equal(t, response.InquiryTicketAmount{}, result)
	})
}

// func TestShowTickets(t *testing.T) {
// 	// Setup
// 	setup()
// 	defer repo.AssertExpectations(t)

// 	t.Run("success", func(t *testing.T) {
// 		// Mock repository calls
// 		tickets := []entity.Ticket{
// 			{
// 				ID:        1,
// 				Capacity:  100,
// 				Region:    "exampleRegion",
// 				EventDate: time.Now(),
// 				CreatedAt: time.Time{},
// 				UpdatedAt: sql.NullTime{},
// 				DeletedAt: sql.NullTime{},
// 			},
// 			{
// 				ID:        2,
// 				Capacity:  100,
// 				Region:    "exampleRegion",
// 				EventDate: time.Now(),
// 				CreatedAt: time.Time{},
// 				UpdatedAt: sql.NullTime{},
// 				DeletedAt: sql.NullTime{},
// 			},
// 		}
// 		page := 1
// 		pageSize := 10

// 		repo.On("FindAllTickets", ctx).Return(tickets, nil)

// 		// Call the function
// 		result, _, _, err := uc.ShowTickets(ctx, page, pageSize, 1)

// 		// Assertions
// 		require.NoError(t, err)
// 		require.Len(t, result, 2)
// 	})

// 	t.Run("error", func(t *testing.T) {
// 		page := 1
// 		pageSize := 10
// 		// Mock repository calls
// 		repo.On("FindAllTickets", ctx).Return([]entity.Ticket{}, errors.New("error"))

// 		// Call the function
// 		result, _, _, err := uc.ShowTickets(ctx, page, pageSize, 2)

// 		// Assertions
// 		require.Error(t, err)
// 		assert.Nil(t, result)
// 	})
// }