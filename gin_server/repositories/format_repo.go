//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./mock/$GOFILE
package repositories

import (
	"context"
	"gin_server/infra"
	"gin_server/models"
)

type IFormateRepo interface {
	Create(ctx context.Context, st models.FormatData) error
}

type FormatRepo struct {
	FormatRepo infra.IMyPostgres
}

var _ IFormateRepo = (*FormatRepo)(nil)

func NewFormatRepo(c context.Context) *FormatRepo {
	return &FormatRepo{
		FormatRepo: infra.NewDB(c),
	}
}

func (f *FormatRepo) Create(ctx context.Context, st models.FormatData) error {
	return f.FormatRepo.Create(ctx, st)
}
