package backend

import (
	"context"

	"github.com/headblockhead/landmine/models"
)

type Backend interface {
	List(ctx context.Context, request models.ListRecordsRequest) (response models.ListRecordsResponse, err error)
}
