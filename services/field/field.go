package services

import (
	"bytes"
	"context"
	"field-service/common/storage"
	"field-service/common/util"
	errConstant "field-service/constants/error"
	"field-service/domain/dto"
	"field-service/domain/models"
	"field-service/repositories"
	"fmt"
	"io"
	"mime/multipart"
	"path"
	"strings"
	"time"

	"github.com/google/uuid"
)

type FieldService struct {
	repository repositories.IRepositoryRegistry
	storage    storage.IStorageClient
}

type IFieldService interface {
	GetAllWithPagination(context.Context, *dto.FieldRequestParam) (*util.PaginationResult, error)
	GetAllWithoutPagination(ctx context.Context) ([]dto.FieldResponse, error)
	GetByUUID(ctx context.Context, uuid string) (*dto.FieldResponse, error)
	Create(context.Context, *dto.FieldRequest) (*dto.FieldResponse, error)
	Update(context.Context, string, *dto.UpdateFieldRequest) (*dto.FieldResponse, error)
	Delete(context.Context, string) error
}

func NewFieldService(repository repositories.IRepositoryRegistry, store storage.IStorageClient) IFieldService {
	return &FieldService{repository: repository, storage: store}
}

func (s *FieldService) GetAllWithPagination(ctx context.Context,
	param *dto.FieldRequestParam) (*util.PaginationResult, error) {
	fields, total, err := s.repository.GetField().FindAllWithPagination(ctx, param)
	if err != nil {
		return nil, err
	}

	fieldResult := make([]dto.FieldResponse, 0, len(fields))
	for _, field := range fields {
		fieldResult = append(fieldResult, dto.FieldResponse{
			UUID:         field.UUID,
			Code:         field.Code,
			Name:         field.Name,
			PricePerHour: field.PricePerHour,
			Images:       field.Images,
			CreatedAt:    field.CreatedAt,
			UpdatedAt:    field.UpdatedAt,
		})
	}

	pagination := &util.PaginationParam{
		Count: total,
		Page:  param.Page,
		Limit: param.Limit,
		Data:  fieldResult,
	}

	response := util.GeneratePagination(*pagination)
	return &response, nil
}

func (s *FieldService) GetAllWithoutPagination(ctx context.Context) ([]dto.FieldResponse, error) {
	fields, err := s.repository.GetField().FindAllWithoutPagination(ctx)
	if err != nil {
		return nil, err
	}

	fieldResult := make([]dto.FieldResponse, 0, len(fields))
	for _, field := range fields {
		fieldResult = append(fieldResult, dto.FieldResponse{
			UUID:         field.UUID,
			Name:         field.Name,
			PricePerHour: field.PricePerHour,
			Images:       field.Images,
		})
	}

	return fieldResult, nil
}

func (s *FieldService) GetByUUID(ctx context.Context, uuid string) (*dto.FieldResponse, error) {
	field, err := s.repository.GetField().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	return &dto.FieldResponse{
		UUID:         field.UUID,
		Code:         field.Code,
		Name:         field.Name,
		PricePerHour: field.PricePerHour,
		Images:       field.Images,
		CreatedAt:    field.CreatedAt,
		UpdatedAt:    field.UpdatedAt,
	}, nil
}

func (s *FieldService) validateUpload(images []multipart.FileHeader) error {
	if images == nil || len(images) == 0 {
		return errConstant.ErrInvalidUploadFile
	}

	for _, image := range images {
		if image.Size > 5*1024*1024 {
			return errConstant.ErrSizeTooBig
		}
	}
	return nil
}

// proses upload file
func (s *FieldService) processAndUploadImages(ctx context.Context, image multipart.FileHeader) (string, error) {
	file, err := image.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	buffer := new(bytes.Buffer)
	_, err = io.Copy(buffer, file)
	if err != nil {
		return "", err
	}

	webpData, err := convertToWebP(buffer.Bytes())
	if err != nil {
		return "", err
	}

	base := strings.TrimSuffix(path.Base(image.Filename), path.Ext(image.Filename))
	filename := fmt.Sprintf("images/%d-%s.webp", time.Now().Unix(), base)
	url, err := s.storage.UploadFile(ctx, filename, webpData)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (s *FieldService) uploadImage(ctx context.Context, images []multipart.FileHeader) ([]string, error) {
	err := s.validateUpload(images)
	if err != nil {
		return nil, err
	}

	urls := make([]string, 0, len(images))
	for _, image := range images {
		url, err := s.processAndUploadImages(ctx, image)
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}

	return urls, nil
}

func (s *FieldService) Create(ctx context.Context, req *dto.FieldRequest) (*dto.FieldResponse, error) {
	imageUrl, err := s.uploadImage(ctx, req.Images)
	if err != nil {
		return nil, err
	}

	field, err := s.repository.GetField().Create(ctx, &models.Field{
		Code:         req.Code,
		Name:         req.Name,
		PricePerHour: req.PricePerHour,
		Images:       imageUrl,
	})
	if err != nil {
		return nil, err
	}

	response := dto.FieldResponse{
		UUID:         field.UUID,
		Code:         field.Code,
		Name:         field.Name,
		PricePerHour: field.PricePerHour,
		Images:       field.Images,
		CreatedAt:    field.CreatedAt,
		UpdatedAt:    field.UpdatedAt,
	}

	return &response, nil
}

func (s *FieldService) Update(ctx context.Context, uuidParam string, req *dto.UpdateFieldRequest) (*dto.FieldResponse, error) {
	field, err := s.repository.GetField().FindByUUID(ctx, uuidParam)
	if err != nil {
		return nil, err
	}

	var imageUrl []string
	if req.Images == nil {
		imageUrl = field.Images
	} else {
		imageUrl, err = s.uploadImage(ctx, req.Images)
		if err != nil {
			return nil, err
		}
	}

	fieldResult, err := s.repository.GetField().Update(ctx, uuidParam, &models.Field{
		Code:         req.Code,
		Name:         req.Name,
		PricePerHour: req.PricePerHour,
		Images:       imageUrl,
	})
	if err != nil {
		return nil, err
	}

	uuidParse, _ := uuid.Parse(uuidParam)
	return &dto.FieldResponse{
		UUID:         uuidParse,
		Code:         fieldResult.Code,
		Name:         fieldResult.Name,
		PricePerHour: fieldResult.PricePerHour,
		Images:       fieldResult.Images,
		CreatedAt:    fieldResult.CreatedAt,
		UpdatedAt:    fieldResult.UpdatedAt,
	}, nil
}

func (s *FieldService) Delete(ctx context.Context, uuid string) error {
	_, err := s.repository.GetField().FindByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	err = s.repository.GetField().Delete(ctx, uuid)
	if err != nil {
		return err
	}

	return nil
}
