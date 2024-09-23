//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./mock/$GOFILE
package services

import (
	"context"
	"gin_server/dto"
	"gin_server/models"
	"gin_server/repositories"
	"gin_server/utils"
)

type IFormatService interface {
	Create(ctx context.Context, req dto.FormatDataRequest) (models.FormatData, error)
}

type FormatService struct {
	FormatRepo repositories.IFormateRepo
}

var _ IFormatService = (*FormatService)(nil)

func NewFormatService(c context.Context) *FormatService {
	return &FormatService{
		FormatRepo: repositories.NewFormatRepo(c),
	}
}

func (f *FormatService) Create(ctx context.Context, req dto.FormatDataRequest) (models.FormatData, error) {
	id := utils.UniqueID(24)

	data := models.FormatData{
		ID:        id,
		MachineID: req.MachineID,
		Judgement: req.Judgement,
	}
	for _, v := range req.Traceability {
		data.AddTraceabilityDetail(v.Name, v.Data)
	}
	return data, f.FormatRepo.Create(ctx, data)
}
