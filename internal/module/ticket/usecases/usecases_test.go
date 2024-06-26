package usecases_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"ticket-service/internal/module/ticket/mocks"
	"ticket-service/internal/module/ticket/models/entity"
	"ticket-service/internal/module/ticket/models/response"
	"ticket-service/internal/module/ticket/usecases"
	"ticket-service/internal/pkg/gorules"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gorules/zen-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	uc                   usecases.Usecases
	repo                 *mocks.Repositories
	ctx                  context.Context
	onlineTicketRulesBre zen.Decision
	p                    message.Publisher
)

type mockPublisher struct{}

// Close implements message.Publisher.
func (m *mockPublisher) Close() error {
	return nil
}

// Publish implements message.Publisher.
func (m *mockPublisher) Publish(topic string, messages ...*message.Message) error {
	return nil
}

func NewMockPublisher() message.Publisher {
	return &mockPublisher{}
}

func setup() {
	repo = new(mocks.Repositories)
	// init business rules engine
	pathOnlineTicket := "../../../../assets/online-ticket-weight.json"
	onlineTicketRulesBre, _ = gorules.Init(pathOnlineTicket)
	fmt.Println("onlineTicketRulesBre", onlineTicketRulesBre)
	p = NewMockPublisher()
	uc = usecases.New(repo, p, onlineTicketRulesBre)
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

func TestShowTickets(t *testing.T) {
	// Setup
	setup()
	defer repo.AssertExpectations(t)

	t.Run("success", func(t *testing.T) {
		// Mock repository calls
		tickets := []entity.Ticket{
			{
				ID:        1,
				Capacity:  100,
				Region:    "exampleRegion",
				EventDate: time.Now(),
				CreatedAt: time.Time{},
				UpdatedAt: sql.NullTime{},
				DeletedAt: sql.NullTime{},
			},
			{
				ID:        2,
				Capacity:  100,
				Region:    "exampleRegion",
				EventDate: time.Now(),
				CreatedAt: time.Time{},
				UpdatedAt: sql.NullTime{},
				DeletedAt: sql.NullTime{}},
		}
		ticketDetails := []entity.TicketDetail{
			{
				ID:        1,
				TicketID:  1,
				Level:     "exampleLevel",
				Stock:     10,
				BasePrice: 100.00,
				CreatedAt: time.Time{},
				UpdatedAt: sql.NullTime{},
				DeletedAt: sql.NullTime{},
			},
			{
				ID:        2,
				TicketID:  2,
				Level:     "exampleLevel",
				Stock:     10,
				BasePrice: 100.00,
				CreatedAt: time.Time{},
				UpdatedAt: sql.NullTime{},
				DeletedAt: sql.NullTime{},
			},
		}

		profileResponse := response.Profile{
			ID:             1,
			UserID:         2,
			FirstName:      "exampleFirstName",
			LastName:       "exampleLastName",
			Address:        "exampleAddress",
			District:       "kota",
			City:           "exampleCity",
			State:          "exampleState",
			Country:        "exampleCountry",
			Region:         "Online",
			Phone:          "examplePhone",
			PersonalID:     "examplePersonalID",
			TypePersonalID: "exampleTypePersonalID",
		}

		responseOnlineTicket := response.OnlineTicket{
			IsSoldOut:      false,
			IsFirstSoldOut: false,
		}
		page := 1
		pageSize := 10

		repo.On("FindTickets", ctx, page, pageSize).Return(tickets, 0, 0, nil)
		repo.On("FindTicketDetails", ctx, page, pageSize).Return(ticketDetails, 0, 0, nil)
		repo.On("GetProfile", ctx, profileResponse.UserID).Return(profileResponse, nil)
		repo.On("GetTicketOnline", ctx, "Online").Return(responseOnlineTicket, nil)
		repo.On("FindTicketByRegionName", ctx, "Online").Return(tickets[0], nil)

		// Call the function
		result, _, _, err := uc.ShowTickets(ctx, page, pageSize, int64(profileResponse.UserID))

		// Assertions
		require.NoError(t, err)
		require.Len(t, result, 2)
	})

	t.Run("error", func(t *testing.T) {
		page := 1
		pageSize := 11
		// Mock repository calls
		repo.On("FindTickets", ctx, page, pageSize).Return([]entity.Ticket{}, 0, 0, errors.New("error"))

		// Call the function
		result, size, page, err := uc.ShowTickets(ctx, page, pageSize, 1)

		// Assertions
		assert.NotNil(t, err)
		assert.Nil(t, result)
		assert.Equal(t, 0, size)
		assert.Equal(t, 0, page)
	})
}
