package service

import (
	"context"
	"field-service/domain/dto"
	"field-service/domain/models"
	"field-service/repositories"
)

type TimeService struct {
	repository repositories.IRepositoryRegistry
}

type ITimeService interface {
	GetAll(context.Context) ([]dto.TimeResponse, error)
	GetByUUID(context.Context, string) (*dto.TimeResponse, error)
	Create(context.Context, *dto.TimeRequest) (*dto.TimeResponse, error)
}

func NewTimeService(repository repositories.IRepositoryRegistry) ITimeService {
	return &TimeService{repository: repository}
}

func (s *TimeService) GetAll(ctx context.Context) ([]dto.TimeResponse, error) {
	times, err := s.repository.GetTime().FindAll(ctx)
	if err != nil {
		return nil, err
	}
	timeResult := make([]dto.TimeResponse, 0, len(times))
	for _, time := range times {
		timeResult = append(timeResult, dto.TimeResponse{
			UUID:      time.UUID,
			StartTime: time.StartTime,
			EndTime:   time.EndTime,
			CreatedAt: time.CreatedAt,
			UpdatedAt: time.UpdatedAt,
		})
	}

	return timeResult, nil
}

func (s *TimeService) GetByUUID(ctx context.Context, uuid string) (*dto.TimeResponse, error) {
	time, err := s.repository.GetTime().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	return &dto.TimeResponse{
		UUID:      time.UUID,
		StartTime: time.StartTime,
		EndTime:   time.EndTime,
		CreatedAt: time.CreatedAt,
		UpdatedAt: time.UpdatedAt,
	}, nil
}

func (s *TimeService) Create(ctx context.Context, req *dto.TimeRequest) (*dto.TimeResponse, error) {
	time := &dto.TimeRequest{
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}

	timeResult, err := s.repository.GetTime().Create(ctx, &models.Time{
		StartTime: time.StartTime,
		EndTime:   time.EndTime,
	})
	if err != nil {
		return nil, err
	}

	return &dto.TimeResponse{
		UUID:      timeResult.UUID,
		StartTime: timeResult.StartTime,
		EndTime:   timeResult.EndTime,
		CreatedAt: timeResult.CreatedAt,
		UpdatedAt: timeResult.UpdatedAt,
	}, nil
}
