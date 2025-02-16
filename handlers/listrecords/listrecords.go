package listrecords

import (
	"log/slog"
	"net/http"

	"github.com/headblockhead/landmine/backend"
	"github.com/headblockhead/landmine/models"
	"github.com/headblockhead/landmine/respond"
)

func New(log *slog.Logger, backend backend.Backend) Handler {
	return Handler{
		log:     log,
		backend: backend,
	}
}

type Handler struct {
	log     *slog.Logger
	backend backend.Backend
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req, validationFailures := models.NewListRecordsRequest(r)
	if len(validationFailures) > 0 {
		br := models.NewBadRequest("bad request", validationFailures)
		respond.WithJSON(w, br, http.StatusBadRequest)
		return
	}

	resp, err := h.backend.List(r.Context(), req)
	if err != nil {
		h.log.Error("failed to list records", slog.Any("error", err))
		respond.WithJSON(w, models.NewError("failed to list records", http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	respond.WithJSON(w, resp, http.StatusOK)
}
