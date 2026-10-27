package repositories

import (
	"context"
	"errors"
	errWrap "field-service/common/error"
	"field-service/constants"
	errConstant "field-service/constants/error"
	errFieldSchedule "field-service/constants/error/fieldschedule"
	"field-service/domain/dto"
	"field-service/domain/models"
	"fmt"

	"gorm.io/gorm"
)

type FieldScheduleRepository struct {
	db *gorm.DB
}

type IFieldScheduleRepository interface {
	FindAllWithPagination(context.Context, *dto.FieldScheduleRequestParam) ([]models.FieldSchedule, int64, error)
	FindAllByFieldIDAndDate(context.Context, int, string) ([]models.FieldSchedule, error)
	FindByUUID(context.Context, string) (*models.FieldSchedule, error)
	FindByDateAndTimeID(context.Context, string, int, int) (*models.FieldSchedule, error)
	Create(context.Context, []models.FieldSchedule) error
	Update(context.Context, string, *models.FieldSchedule) (*models.FieldSchedule, error)
	UpdateStatus(context.Context, constants.FieldScheduleStatus, string) error
	Delete(context.Context, string) error
}

func NewFieldScheduleRepository(db *gorm.DB) IFieldScheduleRepository {
	return &FieldScheduleRepository{db: db}
}

func (r *FieldScheduleRepository) FindAllWithPagination(
	ctx context.Context,
	param *dto.FieldScheduleRequestParam) ([]models.FieldSchedule, int64, error) {
	var (
		fieldSchedule []models.FieldSchedule
		sort          string
		total         int64
	)

	if param.SortColumn != "" {
		sort = fmt.Sprintf("%s %s", param.SortColumn, param.SortOrder)
	} else {
		sort = "created_at desc"
	}

	limit := param.Limit
	offset := (param.Page - 1) * limit

	err := r.db.WithContext(ctx).
		Preload("Field").
		Preload("Time").
		Limit(limit).
		Offset(offset).
		Order(sort).
		Find(&fieldSchedule).
		Error
	if err != nil {
		return nil, 0, errWrap.WrapError(errConstant.ErrSQLError)
	}

	err = r.db.WithContext(ctx).Model(&fieldSchedule).Count(&total).Error
	if err != nil {
		return nil, 0, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return fieldSchedule, total, nil
}

func (r *FieldScheduleRepository) FindAllByFieldIDAndDate(ctx context.Context, fieldID int, date string) ([]models.FieldSchedule, error) {
	var fieldSchedule []models.FieldSchedule
	err := r.db.WithContext(ctx).
		Preload("Field").
		Preload("Time").
		Where("field_id = ?", fieldID).
		Where("date = ?", date).
		Joins("left join times on field_schedules.time_id = times.id").
		Order("times.start_time asc").
		Find(&fieldSchedule).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return fieldSchedule, nil
}

func (r *FieldScheduleRepository) FindByUUID(ctx context.Context, uuid string) (*models.FieldSchedule, error) {
	var fieldSchedule models.FieldSchedule
	err := r.db.WithContext(ctx).
		Preload("Field").
		Preload("Time").
		Where("uuid = ?", uuid).
		First(&fieldSchedule).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errWrap.WrapError(errFieldSchedule.ErrFieldScheduleNotFound)
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return &fieldSchedule, nil
}

func (r *FieldScheduleRepository) FindByDateAndTimeID(ctx context.Context, date string, timeID int, fieldID int) (*models.FieldSchedule, error) {
	var fieldSchedule models.FieldSchedule
	err := r.db.WithContext(ctx).
		Where("date = ?", date).
		Where("tiem_id = ?", timeID).
		Where("field_id = ?", fieldID).
		First(&fieldSchedule).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errWrap.WrapError(errFieldSchedule.ErrFieldScheduleNotFound)
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return &fieldSchedule, nil
}

func (r *FieldScheduleRepository) Create(ctx context.Context, req []models.FieldSchedule) error {
	err := r.db.WithContext(ctx).Create(&req).Error
	if err != nil {
		return errWrap.WrapError(errConstant.ErrSQLError)
	}

	return nil
}

func (r *FieldScheduleRepository) Update(ctx context.Context, uuid string, req *models.FieldSchedule) (*models.FieldSchedule, error) {
	fieldSchedule, err := r.FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}
	fieldSchedule.Date = req.Date

	err = r.db.WithContext(ctx).Save(&fieldSchedule).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return fieldSchedule, nil
}

func (r *FieldScheduleRepository) UpdateStatus(ctx context.Context, status constants.FieldScheduleStatus, uuid string) error {
	fieldSchedule, err := r.FindByUUID(ctx, uuid)
	if err != nil {
		return err
	}
	fieldSchedule.Status = status

	err = r.db.WithContext(ctx).Save(&fieldSchedule).Error
	if err != nil {
		return errWrap.WrapError(errConstant.ErrSQLError)
	}

	return nil
}

func (r *FieldScheduleRepository) Delete(ctx context.Context, uuid string) error {
	err := r.db.WithContext(ctx).Where("uuid = ?", uuid).Delete(&models.FieldSchedule{}, uuid).Error
	if err != nil {
		return errWrap.WrapError(errConstant.ErrSQLError)
	}
	return nil
}
