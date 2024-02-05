package handler

import (
	"encoding/json"
	"errors"
	common "github.com/matnich89/benefex/common/model"
	"github.com/matnich89/benefex/distribution/model"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type MockVinylSender struct {
	mock.Mock
}

func (m *MockVinylSender) SendVinylOrder(order model.VinylOrder) error {
	args := m.Called(order)
	return args.Error(0)
}

type MockCDSender struct {
	mock.Mock
}

func (m *MockCDSender) SendCDOrder(order model.CDOrder) error {
	args := m.Called(order)
	return args.Error(0)
}

func TestHandleMessageVinylOrder(t *testing.T) {
	mockVinylSender := new(MockVinylSender)
	mockCDSender := new(MockCDSender)
	errC := make(chan error, 1)

	handler := NewMessageHandler(mockVinylSender, mockCDSender, errC)

	vinylOrder := model.VinylOrder{
		Artist:        "Test Artist",
		Title:         "Test Album",
		Quantity:      100,
		DateOfRelease: "2022-01-01",
	}
	mockVinylSender.On("SendVinylOrder", vinylOrder).Return(nil)

	release := common.Release{
		Artist:       "Test Artist",
		Title:        "Test Album",
		Genre:        "Rock",
		ReleaseDate:  time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		Distribution: []common.ReleaseDistribution{{Type: "vinyl", Qty: 100}},
	}
	body, _ := json.Marshal(release)
	d := amqp.Delivery{Body: body}

	handler.HandleMessage(d)

	close(errC)

	/*
	 This is very hacky... but the 'failed to acknowledge' error is being triggered
	 because we are just using a dummy amqp.Delivery, it will take a lot of work
	 to abstract this to prevent this, so for the sake of time I am using this dodgy check
	 to ensure no other unexpected errors are present on the channel ( sorry :) )
	*/
	for err := range errC {
		require.EqualError(t, err, "failed to acknowledge message")
	}

	mockVinylSender.AssertCalled(t, "SendVinylOrder", mock.AnythingOfType("model.VinylOrder"))

}

func TestHandleMessageVinylOrderWithError(t *testing.T) {
	mockVinylSender := new(MockVinylSender)
	mockCDSender := new(MockCDSender)
	errC := make(chan error, 2)

	handler := NewMessageHandler(mockVinylSender, mockCDSender, errC)

	vinylOrder := model.VinylOrder{
		Artist:        "Test Artist",
		Title:         "Test Album",
		Quantity:      100,
		DateOfRelease: "2022-01-01",
	}

	expectedError := errors.New("vinyl order failed")
	mockVinylSender.On("SendVinylOrder", vinylOrder).Return(expectedError)

	release := common.Release{
		Artist:       "Test Artist",
		Title:        "Test Album",
		Genre:        "Rock",
		ReleaseDate:  time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		Distribution: []common.ReleaseDistribution{{Type: "vinyl", Qty: 100}},
	}
	body, _ := json.Marshal(release)
	d := amqp.Delivery{Body: body}

	handler.HandleMessage(d)

	close(errC)

	receivedError := <-errC
	require.EqualError(t, receivedError, expectedError.Error())

	mockVinylSender.AssertCalled(t, "SendVinylOrder", mock.AnythingOfType("model.VinylOrder"))
}
