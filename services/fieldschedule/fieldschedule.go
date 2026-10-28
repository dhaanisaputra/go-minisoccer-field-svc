package service

import (
	"context"
	"field-service/common/util"
	"field-service/constants"
	errFieldSchedule "field-service/constants/error/fieldschedule"
	"field-service/domain/dto"
	"field-service/domain/models"
	"field-service/repositories"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type FieldScheduleService struct {
	repository repositories.IRepositoryRegistry
}

type IFieldScheduleService interface {
	GetAllWithPagination(context.Context, *dto.FieldScheduleRequestParam) (*util.PaginationResult, error)
	GetAllByFieldIDAndDate(context.Context, string, string) ([]dto.FieldScheduleForBookingResponse, error)
	GetByUUID(context.Context, string) (*dto.FieldScheduleResponse, error)
	GenerateScheduleForOneMonth(context.Context, *dto.GenerateFieldScheduleForOneMonthRequest) error
	Create(context.Context, *dto.FieldScheduleRequest) error
	Update(context.Context, string, *dto.UpdateFieldScheduleRequest) (*dto.FieldScheduleResponse, error)
	UpdateStatus(context.Context, *dto.UpdateStatusFieldScheduleRequest) error
	Delete(context.Context, string) error
}

func NewFieldScheduleService(repository repositories.IRepositoryRegistry) IFieldScheduleService {
	return &FieldScheduleService{repository: repository}
}

func (s *FieldScheduleService) GetAllWithPagination(ctx context.Context,
	param *dto.FieldScheduleRequestParam) (*util.PaginationResult, error) {
	fieldSchedules, total, err := s.repository.GetFieldSchedule().FindAllWithPagination(ctx, param)
	if err != nil {
		return nil, err
	}

	fieldScheduleResult := make([]dto.FieldScheduleResponse, 0, len(fieldSchedules))
	for _, fieldSchedule := range fieldSchedules {
		fieldScheduleResult = append(fieldScheduleResult, dto.FieldScheduleResponse{
			UUID:         fieldSchedule.UUID,
			FieldName:    fieldSchedule.Field.Name,
			PricePerHour: fieldSchedule.Field.PricePerHour,
			Date:         fieldSchedule.Date.Format("2026-01-02"),
			Status:       fieldSchedule.Status.GetStatusString(),
			Time:         fmt.Sprintf("%s %s", fieldSchedule.Time.StartTime, fieldSchedule.Time.EndTime),
			CreatedAt:    fieldSchedule.CreatedAt,
			UpdatedAt:    fieldSchedule.UpdatedAt,
		})

	}

	pagination := &util.PaginationParam{
		Count: total,
		Page:  param.Page,
		Limit: param.Limit,
		Data:  fieldScheduleResult,
	}

	response := util.GeneratePagination(*pagination)
	return &response, nil
}

func (s *FieldScheduleService) convertMonthName(inputDate string) string {
	date, err := time.Parse(time.DateOnly, inputDate)
	if err != nil {
		return ""
	}

	idnMonth := map[string]string{
		"Jan": "Jan",
		"Feb": "Feb",
		"Mar": "Mar",
		"Apr": "Apr",
		"May": "Mei",
		"Jun": "Jun",
		"Jul": "Jul",
		"Aug": "Ags",
		"Sep": "Sep",
		"Oct": "Okt",
		"Nov": "Nov",
		"Dec": "Des",
	}

	formattedDate := date.Format("02 Jan")
	day := formattedDate[:3]
	month := formattedDate[3:]
	formattedDate = fmt.Sprintf("%s %s", day, idnMonth[month])
	return formattedDate

}

func (s *FieldScheduleService) GetAllByFieldIDAndDate(ctx context.Context, uuid, date string) ([]dto.FieldScheduleForBookingResponse, error) {
	field, err := s.repository.GetField().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	fieldSchedules, err := s.repository.GetFieldSchedule().FindAllByFieldIDAndDate(ctx, int(field.ID), date)
	if err != nil {
		return nil, err
	}

	fieldScheduleResult := make([]dto.FieldScheduleForBookingResponse, 0, len(fieldSchedules))
	for _, schedule := range fieldSchedules {
		pricePerHour := float64(schedule.Field.PricePerHour)
		fieldScheduleResult = append(fieldScheduleResult, dto.FieldScheduleForBookingResponse{
			UUID:         schedule.UUID,
			Date:         s.convertMonthName(schedule.Date.Format(time.DateOnly)),
			Time:         schedule.Time.StartTime,
			Status:       schedule.Status.GetStatusString(),
			PricePerHour: util.RupiahFormat(&pricePerHour),
		})
	}

	return fieldScheduleResult, nil
}

func (s *FieldScheduleService) GetByUUID(ctx context.Context, uuid string) (*dto.FieldScheduleResponse, error) {
	fieldSchedule, err := s.repository.GetFieldSchedule().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	return &dto.FieldScheduleResponse{
		UUID:         fieldSchedule.UUID,
		FieldName:    fieldSchedule.Field.Name,
		PricePerHour: fieldSchedule.Field.PricePerHour,
		Date:         fieldSchedule.Date.Format(time.DateOnly),
		Status:       fieldSchedule.Status.GetStatusString(),
		Time:         fmt.Sprintf("%s - %s", fieldSchedule.Time.StartTime, fieldSchedule.Time.EndTime),
		CreatedAt:    fieldSchedule.CreatedAt,
		UpdatedAt:    fieldSchedule.UpdatedAt,
	}, nil
}

func (s *FieldScheduleService) Create(ctx context.Context, req *dto.FieldScheduleRequest) error {
	field, err := s.repository.GetField().FindByUUID(ctx, req.FieldID)
	if err != nil {
		return err
	}

	fieldSchedule := make([]models.FieldSchedule, 0, len(req.TimeIDs))
	dateParse, _ := time.Parse(time.DateOnly, req.Date)
	for _, timeID := range req.TimeIDs {
		scheduleTime, err := s.repository.GetTime().FindByUUID(ctx, timeID)
		if err != nil {
			return err
		}

		schedule, err := s.repository.GetFieldSchedule().FindByDateAndTimeID(ctx, req.Date, int(scheduleTime.ID), int(field.ID))
		if err != nil {
			return err
		}

		if schedule != nil {
			return errFieldSchedule.ErrFieldScheduleIsExist
		}

		fieldSchedule = append(fieldSchedule, models.FieldSchedule{
			UUID:    uuid.New(),
			FieldID: field.ID,
			TimeID:  scheduleTime.ID,
			Date:    dateParse,
			Status:  constants.Available,
		})

		err = s.repository.GetFieldSchedule().Create(ctx, fieldSchedule)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *FieldScheduleService) GenerateScheduleForOneMonth(ctx context.Context, req *dto.GenerateFieldScheduleForOneMonthRequest) error {
	field, err := s.repository.GetField().FindByUUID(ctx, req.FieldID)
	if err != nil {
		return err
	}

	times, err := s.repository.GetTime().FindAll(ctx)
	if err != nil {
		return err
	}

	numberOfDays := 30
	fieldSchedule := make([]models.FieldSchedule, 0, numberOfDays)
	now := time.Now().Add(time.Duration(1) * 24 * time.Hour)
	for i := 0; i < numberOfDays; i++ {
		currDate := now.AddDate(0, 0, i)
		for _, item := range times {
			schedule, err := s.repository.GetFieldSchedule().FindByDateAndTimeID(ctx, currDate.Format(time.DateOnly), int(item.ID), int(field.ID))
			if err != nil {
				return err
			}

			if schedule != nil {
				return errFieldSchedule.ErrFieldScheduleIsExist
			}

			fieldSchedule = append(fieldSchedule, models.FieldSchedule{
				UUID:    uuid.New(),
				FieldID: field.ID,
				TimeID:  item.ID,
				Date:    currDate,
				Status:  constants.Available,
			})
		}
	}

	err = s.repository.GetFieldSchedule().Create(ctx, fieldSchedule)
	if err != nil {
		return err
	}

	return nil
}

func (s *FieldScheduleService) Update(ctx context.Context, uuid string, req *dto.UpdateFieldScheduleRequest) (*dto.FieldScheduleResponse, error) {
	fieldSchedule, err := s.repository.GetFieldSchedule().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	scheduleTime, err := s.repository.GetTime().FindByUUID(ctx, req.TimeID)
	if err != nil {
		return nil, err
	}

	// cek jadwal available berdasarkan tanggal dan waktu
	isTimeExist, err := s.repository.GetFieldSchedule().FindByDateAndTimeID(ctx, req.Date, int(scheduleTime.ID), int(fieldSchedule.FieldID))
	if err != nil {
		return nil, err
	}

	if isTimeExist != nil && req.Date != fieldSchedule.Date.Format(time.DateOnly) {
		checkDate, err := s.repository.GetFieldSchedule().FindByDateAndTimeID(ctx, req.Date, int(scheduleTime.ID), int(fieldSchedule.FieldID))
		if err != nil {
			return nil, err
		}

		if checkDate != nil {
			return nil, errFieldSchedule.ErrFieldScheduleIsExist
		}
	}

	dateParsed, _ := time.Parse(time.DateOnly, req.Date)
	fieldResult, err := s.repository.GetFieldSchedule().Update(ctx, uuid, &models.FieldSchedule{
		Date:   dateParsed,
		TimeID: scheduleTime.ID,
	})
	if err != nil {
		return nil, err
	}

	return &dto.FieldScheduleResponse{
		UUID:         fieldResult.UUID,
		FieldName:    fieldResult.Field.Name,
		Date:         fieldResult.Date.Format(time.DateOnly),
		PricePerHour: fieldResult.Field.PricePerHour,
		Status:       fieldResult.Status.GetStatusString(),
		Time:         fmt.Sprintf("%s %s", fieldResult.Time.StartTime, fieldResult.Time.EndTime),
		CreatedAt:    fieldResult.CreatedAt,
		UpdatedAt:    fieldResult.UpdatedAt,
	}, nil
}

func (s *FieldScheduleService) UpdateStatus(ctx context.Context, req *dto.UpdateStatusFieldScheduleRequest) error {
	for _, item := range req.FieldScheduleIDs {
		_, err := s.repository.GetFieldSchedule().FindByUUID(ctx, item)
		if err != nil {
			return err
		}

		err = s.repository.GetFieldSchedule().UpdateStatus(ctx, constants.Booked, item)
		if err != nil {
			return err
		}

	}
	return nil
}

func (s *FieldScheduleService) Delete(ctx context.Context, uuid string) error {
	_, err := s.repository.GetFieldSchedule().FindByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	err = s.repository.GetFieldSchedule().Delete(ctx, uuid)
	if err != nil {
		return err
	}
	return nil
}
