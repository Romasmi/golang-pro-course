package grpcserver

import (
	"testing"
	"time"

	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/domain"
	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/pkg/api"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestProtoToDomain(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	protoEvent := &api.Event{
		Id:          "1",
		Title:       "Title",
		Description: "Desc",
		StartAt:     timestamppb.New(now),
		EndAt:       timestamppb.New(now.Add(time.Hour)),
		UserId:      "user1",
	}

	domainEvent := protoToDomain(protoEvent)

	assert.Equal(t, protoEvent.Id, domainEvent.ID)
	assert.Equal(t, protoEvent.Title, domainEvent.Title)
	assert.Equal(t, protoEvent.Description, domainEvent.Description)
	assert.True(t, protoEvent.StartAt.AsTime().Equal(domainEvent.StartAt))
	assert.True(t, protoEvent.EndAt.AsTime().Equal(domainEvent.EndAt))
	assert.Equal(t, protoEvent.UserId, domainEvent.UserID)
}

func TestDomainToProto(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	domainEvent := domain.Event{
		ID:          "1",
		Title:       "Title",
		Description: "Desc",
		StartAt:     now,
		EndAt:       now.Add(time.Hour),
		UserID:      "user1",
	}

	protoEvent := domainToProto(domainEvent)

	assert.Equal(t, domainEvent.ID, protoEvent.Id)
	assert.Equal(t, domainEvent.Title, protoEvent.Title)
	assert.Equal(t, domainEvent.Description, protoEvent.Description)
	assert.Equal(t, domainEvent.StartAt.Unix(), protoEvent.StartAt.AsTime().Unix())
	assert.Equal(t, domainEvent.EndAt.Unix(), protoEvent.EndAt.AsTime().Unix())
	assert.Equal(t, domainEvent.UserID, protoEvent.UserId)
}

func TestProtoToDomain_Nil(t *testing.T) {
	domainEvent := protoToDomain(nil)
	assert.Equal(t, domain.Event{}, domainEvent)
}

func TestDomainToProtoSlice(t *testing.T) {
	events := []domain.Event{
		{ID: "1", Title: "E1"},
		{ID: "2", Title: "E2"},
	}

	protos := domainToProtoSlice(events)

	assert.Len(t, protos, 2)
	assert.Equal(t, "1", protos[0].Id)
	assert.Equal(t, "2", protos[1].Id)
}
