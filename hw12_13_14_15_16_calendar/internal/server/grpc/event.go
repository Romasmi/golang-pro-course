package grpcserver

import (
	"context"
	"fmt"

	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/convert"
	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/domain"
	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/usecases"
	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/pkg/api"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) CreateEvent(ctx context.Context, req *api.CreateEventRequest) (*api.CreateEventResponse, error) {
	uc := s.usecases[usecases.CreateEvent]
	res, err := uc.Do(ctx, protoToDomain(req.Event))
	if err != nil {
		return nil, err
	}
	return &api.CreateEventResponse{Event: domainToProto(res.(domain.Event))}, nil
}

func (s *Server) UpdateEvent(ctx context.Context, req *api.UpdateEventRequest) (*api.UpdateEventResponse, error) {
	uc := s.usecases[usecases.UpdateEvent]
	res, err := uc.Do(ctx, protoToDomain(req.Event))
	if err != nil {
		return nil, err
	}
	return &api.UpdateEventResponse{Event: domainToProto(res.(domain.Event))}, nil
}

func (s *Server) GetEvent(ctx context.Context, req *api.GetEventRequest) (*api.GetEventResponse, error) {
	uc := s.usecases[usecases.GetEvent]
	res, err := uc.Do(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &api.GetEventResponse{Event: domainToProto(res.(domain.Event))}, nil
}

func (s *Server) ListEvents(ctx context.Context, _ *api.ListEventsRequest) (*api.ListEventsResponse, error) {
	uc := s.usecases[usecases.ListEvents]
	res, err := uc.Do(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &api.ListEventsResponse{Events: domainToProtoSlice(res.([]domain.Event))}, nil
}

func (s *Server) ListEventsByInterval(
	ctx context.Context,
	req *api.ListEventsByIntervalRequest,
) (*api.ListEventsResponse, error) {
	uc := s.usecases[usecases.ListEventsByInterval]

	interval := convert.IntervalFromProto(req.Interval)
	if !interval.IsValid() {
		return nil, fmt.Errorf("invalid interval: %s", interval)
	}

	ucReq := usecases.ListEventsByIntervalRequest{
		Date:     req.Date.AsTime(),
		Interval: convert.IntervalFromProto(req.Interval),
	}

	res, err := uc.Do(ctx, ucReq)
	if err != nil {
		return nil, err
	}
	return &api.ListEventsResponse{Events: domainToProtoSlice(res.([]domain.Event))}, nil
}

func (s *Server) Hello(ctx context.Context, _ *api.HelloRequest) (*api.HelloResponse, error) {
	uc := s.usecases[usecases.Hello]
	res, err := uc.Do(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &api.HelloResponse{Message: res.(string)}, nil
}

func protoToDomain(e *api.Event) domain.Event {
	if e == nil {
		return domain.Event{}
	}
	return domain.Event{
		ID:           e.Id,
		Title:        e.Title,
		Description:  e.Description,
		StartAt:      e.StartAt.AsTime(),
		EndAt:        e.EndAt.AsTime(),
		UserID:       e.UserId,
		RemindBefore: e.RemindBefore,
	}
}

func domainToProto(e domain.Event) *api.Event {
	return &api.Event{
		Id:           e.ID,
		Title:        e.Title,
		Description:  e.Description,
		StartAt:      timestamppb.New(e.StartAt),
		EndAt:        timestamppb.New(e.EndAt),
		UserId:       e.UserID,
		RemindBefore: e.RemindBefore,
	}
}

func domainToProtoSlice(events []domain.Event) []*api.Event {
	res := make([]*api.Event, 0, len(events))
	for _, e := range events {
		res = append(res, domainToProto(e))
	}
	return res
}
