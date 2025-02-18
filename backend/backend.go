package backend

import (
	"context"

	"github.com/headblockhead/landmine/models"
)

type Backend interface {
	List(ctx context.Context, request models.ListRecordsRequest) (response models.ListRecordsResponse, err error)
	Create(ctx context.Context, request models.CreateRecordsRequest) (response models.CreateRecordsResponse, err error)
	DeleteMultiple(ctx context.Context, request models.DeleteRecordsRequest) (response models.DeleteRecordsResponse, err error)
}
