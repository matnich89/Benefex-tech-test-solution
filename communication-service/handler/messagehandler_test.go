package handler

import (
	"encoding/json"
	"fmt"
	common "github.com/matnich89/benefex/common/model"
	"github.com/matnich89/benefex/communcation/model"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"sync"
	"testing"
	"time"
)

type MockFanBaseStore struct {
	mock.Mock
}

func (m *MockFanBaseStore) GetFansForArtist(artist string) ([]model.Fan, error) {
	args := m.Called(artist)
	return args.Get(0).([]model.Fan), args.Error(1)
}

type MockEmailSender struct {
	mu              sync.Mutex
	numberOfInvokes int
	wg              sync.WaitGroup
}

func (m *MockEmailSender) SendEmail(email model.FanEmail) {
	m.mu.Lock()
	defer m.mu.Unlock()
	defer m.wg.Done()
	m.numberOfInvokes++
}

func TestHandleMessage_Success(t *testing.T) {

	numberOfFans := 5000

	mockFanBaseDb := &MockFanBaseStore{}
	mockEmailClient := &MockEmailSender{}

	mockEmailClient.wg.Add(numberOfFans)

	fans := generateFans(numberOfFans)
	mockFanBaseDb.On("GetFansForArtist", "TestArtist").Return(fans, nil)

	handler := NewMessageHandler(make(chan error), mockFanBaseDb, mockEmailClient)

	release := common.Release{
		Artist:      "TestArtist",
		Title:       "New Album",
		ReleaseDate: time.Now().AddDate(0, 1, 0),
	}

	body, _ := json.Marshal(release)

	delivery := amqp.Delivery{Body: body}

	handler.HandleMessage(delivery)

	mockEmailClient.wg.Wait()

	mockFanBaseDb.AssertExpectations(t)
	assert.Equal(t, numberOfFans, mockEmailClient.numberOfInvokes)
}

func generateFans(numFans int) []model.Fan {
	fans := make([]model.Fan, numFans)
	for i := 0; i < numFans; i++ {
		fans[i] = model.Fan{EmailAddress: fmt.Sprintf("fan%d@example.com", i+1)}
	}
	return fans
}
