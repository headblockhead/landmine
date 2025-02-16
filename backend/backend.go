package backend

import (
	"context"

	"github.com/headblockhead/landmine/models"
)

type Backend interface {
	List(ctx context.Context, request models.ListRecordsRequest) (response models.ListRecordsResponse, err error)
	Create(ctx context.Context, request models.CreateRecordsRequest) (response models.CreateRecordsResponse, err error)
}
