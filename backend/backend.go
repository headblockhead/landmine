package backend

import (
	"context"

	"github.com/headblockhead/landmine/models"
)

type Backend interface {
	List(ctx context.Context, request models.ListGetRequest) (response models.ListGetResponse, err error)
}
